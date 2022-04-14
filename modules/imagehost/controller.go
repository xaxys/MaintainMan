package imagehost

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getImage godoc
// @Summary 获取图片
// @Description 根据图片UUID 获取图片 可以自定义图片变化
// @Tags image
// @Produce json
// @Produce image/*
// @Param id path string true "Image UUID"
// @Param param query string false "Transformation parameters"
// @Success 200 {object} string "Image data"
// @Failure 400 {object} model.ApiJson{data=[]string} "Error message"
// @Failure 403 {object} model.ApiJson{data=[]string} "Error message"
// @Failure 404 {object} model.ApiJson{data=[]string} "Error message"
// @Failure 500 {object} model.ApiJson{data=[]string} "Error message"
// @Router /v1/image/{id} [get]
func getImage(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	param := ctx.URLParam("param")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getImageService(id, param, auth)
	if response.ApiRes != nil {
		ctx.Values().Set("response", response.ApiRes)
		return
	}
	ctx.ContentType(response.Format)
	ctx.StatusCode(iris.StatusOK)
	ctx.Write(response.Data)
}

// uploadImage godoc
// @Summary 上传图片
// @Description 上传图片
// @Tags image
// @Accept multipart/form-data
// @Produce  json
// @Param image formData file true "Image file"
// @Success 200 {object} model.ApiJson{data=string} "Image UUID"
// @Failure 400 {object} model.ApiJson{data=[]string} "Error message"
// @Failure 500 {object} model.ApiJson{data=[]string} "Error message"
// @Router /v1/image [post]
func uploadImage(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(imageConfig.GetInt64("upload.max_file_size"))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	file, _, err := ctx.FormFile("image")
	if err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
	}
	response := uploadImageService(file, auth)
	ctx.Values().Set("response", response)
}
