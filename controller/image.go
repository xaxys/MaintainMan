package controller

import (
	"maintainman/config"
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetImage(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	param := ctx.URLParam("param")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetImage(id, param, auth)
	ctx.Values().Set("response", response)
}

func UploadImage(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(config.ImageConfig.GetInt64("upload.max_file_size"))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	file, _, err := ctx.FormFile("image")
	if err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
	}
	response := service.UploadImage(file, auth)
	ctx.Values().Set("response", util.Tenary(response.Status, response.Data, any(response)))
}
