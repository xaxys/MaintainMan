package service

import (
	"bytes"
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

func GetImage(id, param string, auth *model.AuthInfo) *model.ApiJson {
	// parse param to transformation
	trans := dao.GetTransformation(param)
	if trans == nil {
		if dao.HasPermission(auth.Role, "image.custom") {
			transParam, err := util.ParseParameters(param)
			if err != nil {
				return model.ErrorInvalidData(err)
			}
			trans = util.NewTransformation(&transParam)
		}
	}

	key := id + trans.Hash
	cid, cached := cache.Cache.Get(key)
	if cached {
		id = cid.(string)
	}

	if !dao.ExistImage(id) {
		return model.ErrorNotFound(fmt.Errorf("未找到图片: cached: %v, id: %s", cached, id))
	}

	image, data, format, err := dao.LoadImage(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}

	// do transformation
	if trans != nil && !cached {
		image = util.TransformCropAndResize(image, trans)
		tid := uuid.New().String()
		bytes, err := dao.SaveImage(tid, format, image)
		if err != nil {
			return model.ErrorInsertDatabase(err)
		}
		cache.Cache.SetWithCost(key, tid, int64(len(bytes)), 0)
		data = bytes
	}

	return model.Success(string(data), "获取成功")
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

	eagerTransform := func() {
		img, format, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return
		}
		for _, trans := range dao.GetEagerTransformation() {
			imgNew := util.TransformCropAndResize(img, trans, auth)
			id := uuid.New().String()
			key := id + trans.Hash
			bytes, err := dao.SaveImage(id, format, imgNew)
			if err == nil {
				cache.Cache.SetWithCost(key, id, int64(len(bytes)), 0)
			}
		}
	}

	id := uuid.New().String()
	if config.ImageConfig.GetBool("upload.async") {
		go func() {
			if err := dao.SaveImageBytes(id, format, data); err != nil {
				logger.Logger.Warnf("保存图片失败(id:%s): %+v", id, err)
			}
			go eagerTransform()
		}()
	} else {
		if err := dao.SaveImageBytes(id, format, data); err != nil {
			return model.ErrorInsertDatabase(err)
		}
		go eagerTransform()
	}
	return model.Success(id, "上传成功")
}