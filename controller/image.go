package controller

import (
	"maintainman/config"
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetImage godoc
// @Summary Get image
// @Description Get image
// @Tags image
// @Produce json
// @Produce image/*
// @Param id path string true "Image UUID"
// @Param param query string false "Transformation parameters"
// @Success 200 {object} string "Image data"
// @Failure 400 {object} model.ApiJson "Error message"
// @Failure 403 {object} model.ApiJson "Error message"
// @Failure 404 {object} model.ApiJson "Error message"
// @Failure 500 {object} model.ApiJson "Error message"
// @Router /image/{id} [get]
func GetImage(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	param := ctx.URLParam("param")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetImage(id, param, auth)
	if response.ApiRes != nil {
		ctx.Values().Set("response", response.ApiRes)
		return
	}
	ctx.ContentType(response.Format)
	ctx.StatusCode(iris.StatusOK)
	ctx.Write(response.Data)
}

// UploadImage godoc
// @Summary Upload image
// @Description Upload image
// @Tags image
// @Accept multipart/form-data
// @Produce  json
// @Param image formData file true "Image file"
// @Success 200 {object} model.ApiJson{data=string} "Image UUID"
// @Failure 400 {object} model.ApiJson "Error message"
// @Failure 500 {object} model.ApiJson "Error message"
// @Router /image [post]
func UploadImage(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(config.ImageConfig.GetInt64("upload.max_file_size"))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	file, _, err := ctx.FormFile("image")
	if err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
	}
	response := service.UploadImage(file, auth)
	ctx.Values().Set("response", response)
}
