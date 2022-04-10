package service

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"maintainman/cache"
	"maintainman/config"
	"maintainman/dao"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/util"
	"mime/multipart"

	"github.com/google/uuid"
)

type ImageResponse struct {
	Data   []byte
	Format string
	ApiRes *model.ApiJson
}

func GetImage(id, param string, auth *model.AuthInfo) *ImageResponse {
	// parse param to transformation
	trans, ok := dao.GetTransformation(param)
	if !ok {
		if err := dao.CheckPermission(auth.Role, "image.custom"); err != nil {
			return &ImageResponse{ApiRes: model.ErrorNoPermissions(err)}
		}
		transParam, err := util.ParseParameters(param)
		if err != nil {
			return &ImageResponse{ApiRes: model.ErrorInvalidData(err)}
		}
		trans = util.NewTransformation(&transParam)
	}

	key := id
	if trans != nil {
		key += trans.Hash
	}

	cid, cached := cache.ImageCache.Get(key)
	if cached {
		id = cid.(string)
	}

	if !dao.ExistImage(id, cached) {
		return &ImageResponse{ApiRes: model.ErrorNotFound(fmt.Errorf("未找到图片: cached: %v, id: %s", cached, id))}
	}

	image, data, format, err := dao.LoadImage(id, cached)
	if err != nil {
		return &ImageResponse{ApiRes: model.ErrorQueryDatabase(err)}
	}

	// do transformation
	if trans != nil && !cached {
		uid := parseUUID(id)
		user, err := dao.GetUserByID(uid)
		newAuth := model.AuthInfo{User: uid}
		if err != nil {
			logger.Logger.Warn(err)
		} else {
			newAuth.Name = user.Name
		}
		image = util.TransformCropAndResize(image, trans, newAuth)
		tid := genUUID(uid)
		format = util.Tenary(config.ImageConfig.GetBool("cache_as_jpeg"), "jpeg", format)
		bytes, err := dao.SaveImage(tid, format, image, true)
		if err != nil {
			return &ImageResponse{ApiRes: model.ErrorInsertDatabase(err)}
		}
		cache.ImageCache.SetWithCost(key, tid, int64(len(bytes)), 0)
		data = bytes
	}

	return &ImageResponse{
		Data:   data,
		Format: "image/" + format,
	}
}

func UploadImage(file multipart.File, auth *model.AuthInfo) *model.ApiJson {
	c, format, err := image.DecodeConfig(file)
	if err != nil {
		return model.ErrorValidation(err)
	}
	file.Seek(0, 0)

	max_pixels := config.ImageConfig.GetInt("upload.max_pixels")
	if max_pixels > 0 && c.Width*c.Height > max_pixels {
		return model.ErrorValidation(fmt.Errorf("图片尺寸过大: %d x %d", c.Width, c.Height))
	}

	max_file_size := config.ImageConfig.GetInt64("upload.max_file_size")
	limit := io.LimitReader(file, max_file_size+1)
	data, err := ioutil.ReadAll(limit)
	if err != nil {
		return model.ErrorInvalidData(err)
	}

	id := genUUID(auth.User)
	saveImage := func(errHandler func(error)) {
		var img image.Image
		if config.ImageConfig.GetBool("save_as_jpeg") {
			img, _, err = image.Decode(bytes.NewReader(data))
			if err != nil {
				errHandler(err)
			}
			if _, err := dao.SaveImage(id, "jpeg", img, false); err != nil {
				errHandler(err)
			}
		} else {
			if err := dao.SaveImageBytes(id, format, data, false); err != nil {
				errHandler(err)
			}
		}

		// eager transform
		go func() {
			if img == nil {
				img, format, err = image.Decode(bytes.NewReader(data))
				if err != nil {
					return
				}
			}
			for _, trans := range dao.GetEagerTransformation() {
				imgNew := util.TransformCropAndResize(img, trans, *auth)
				tid := genUUID(auth.User)
				key := tid + trans.Hash
				format = util.Tenary(config.ImageConfig.GetBool("cache_as_jpeg"), "jpeg", format)
				bytes, err := dao.SaveImage(tid, format, imgNew, true)
				if err == nil {
					cache.ImageCache.SetWithCost(key, tid, int64(len(bytes)), 0)
				}
			}
		}()
	}

	response := model.Success(id, "上传成功")
	if config.ImageConfig.GetBool("upload.async") {
		go saveImage(func(err error) {
			logger.Logger.Warnf("保存图片失败(id:%s): %+v", id, err)
		})
	} else {
		saveImage(func(err error) {
			response = model.ErrorInsertDatabase(fmt.Errorf("保存图片失败(id:%s): %+v", id, err))
		})
	}
	return response
}

func genUUID(id uint) string {
	uuidv1, err := uuid.NewUUID()
	if err != nil {
		logger.Logger.Error("生成uuid失败: %+v", err)
	}
	bigint24 := [8]byte{}
	binary.LittleEndian.PutUint64(bigint24[:], uint64(id))
	copy(uuidv1[13:], bigint24[:3])
	return uuidv1.String()
}

func parseUUID(str string) uint {
	uuidv1, err := uuid.Parse(str)
	if err != nil {
		logger.Logger.Error("解析uuid失败: %+v", err)
	}
	bigint24 := [8]byte{}
	copy(bigint24[:3], uuidv1[13:])
	return uint(binary.LittleEndian.Uint64(bigint24[:]))
}
