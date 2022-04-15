package imagehost

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/google/uuid"
)

type imageResponse struct {
	Data   []byte
	Format string
	ApiRes *model.ApiJson
}

func getImageService(id, param string, auth *model.AuthInfo) *imageResponse {
	// parse param to transformation
	trans, ok := getTransformation(param)
	if !ok {
		if err := dao.CheckPermission(auth.Role, "image.custom"); err != nil {
			return &imageResponse{ApiRes: model.ErrorNoPermissions(err)}
		}
		transParam, err := parseParameters(param)
		if err != nil {
			return &imageResponse{ApiRes: model.ErrorInvalidData(err)}
		}
		trans = newTransformation(&transParam)
	}

	key := id
	if trans != nil {
		key += trans.Hash
	}

	cid, cached := mctx.Cache.Get(key)
	if cached {
		id = cid.(string)
	}

	if !existImage(id, cached) {
		return &imageResponse{ApiRes: model.ErrorNotFound(fmt.Errorf("未找到图片: cached: %v, id: %s", cached, id))}
	}

	image, data, format, err := loadImage(id, cached)
	if err != nil {
		return &imageResponse{ApiRes: model.ErrorQueryDatabase(err)}
	}

	// do transformation
	if trans != nil && !cached {
		uid := parseUUID(id)
		user, err := dao.GetUserByID(uid)
		newAuth := model.AuthInfo{User: uid}
		if err != nil {
			mctx.Logger.Warn(err)
		} else {
			newAuth.Name = user.Name
		}
		image = transformCropAndResize(image, trans, newAuth)
		tid := genUUID(uid)
		format = util.Tenary(imageConfig.GetBool("cache_as_jpeg"), "jpeg", format)
		bytes, err := saveImage(tid, format, image, true)
		if err != nil {
			return &imageResponse{ApiRes: model.ErrorInsertDatabase(err)}
		}
		mctx.Cache.SetWithCost(key, tid, int64(len(bytes)), 0)
		data = bytes
	}

	return &imageResponse{
		Data:   data,
		Format: "image/" + format,
	}
}

func uploadImageService(file multipart.File, auth *model.AuthInfo) *model.ApiJson {
	c, format, err := image.DecodeConfig(file)
	if err != nil {
		return model.ErrorValidation(err)
	}
	file.Seek(0, 0)

	max_pixels := imageConfig.GetInt("upload.max_pixels")
	if max_pixels > 0 && c.Width*c.Height > max_pixels {
		return model.ErrorValidation(fmt.Errorf("图片尺寸过大: %d x %d", c.Width, c.Height))
	}

	max_file_size := imageConfig.GetInt64("upload.max_file_size")
	limit := io.LimitReader(file, max_file_size+1)
	data, err := ioutil.ReadAll(limit)
	if err != nil {
		return model.ErrorInvalidData(err)
	}

	id := genUUID(auth.User)
	saveImage := func(errHandler func(error)) {
		var img image.Image
		if imageConfig.GetBool("save_as_jpeg") {
			img, _, err = image.Decode(bytes.NewReader(data))
			if err != nil {
				errHandler(err)
			}
			if _, err := saveImage(id, "jpeg", img, false); err != nil {
				errHandler(err)
			}
		} else {
			if err := saveImageBytes(id, format, data, false); err != nil {
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
			for _, trans := range getEagerTransformation() {
				imgNew := transformCropAndResize(img, trans, *auth)
				tid := genUUID(auth.User)
				key := tid + trans.Hash
				format = util.Tenary(imageConfig.GetBool("cache_as_jpeg"), "jpeg", format)
				bytes, err := saveImage(tid, format, imgNew, true)
				if err == nil {
					mctx.Cache.SetWithCost(key, tid, int64(len(bytes)), 0)
				}
			}
		}()
	}

	response := model.Success(id, "上传成功")
	if imageConfig.GetBool("upload.async") {
		go saveImage(func(err error) {
			mctx.Logger.Warnf("保存图片失败(id:%s): %+v", id, err)
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
		mctx.Logger.Error("生成uuid失败: %+v", err)
	}
	bigint24 := [8]byte{}
	binary.BigEndian.PutUint64(bigint24[:], uint64(id))
	copy(uuidv1[12:], bigint24[4:])
	return uuidv1.String()
}

func parseUUID(str string) uint {
	uuidv1, err := uuid.Parse(str)
	if err != nil {
		mctx.Logger.Error("解析uuid失败: %+v", err)
	}
	bigint24 := [8]byte{}
	copy(bigint24[4:], uuidv1[12:])
	return uint(binary.BigEndian.Uint64(bigint24[:]))
}
