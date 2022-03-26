package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetCommentsByOrder(id, param, auth)
	ctx.Values().Set("response", response)
}

func ForceGetCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceGetCommentsByOrder(id, param, auth)
	ctx.Values().Set("response", response)
}

func ForceCreateComment(ctx iris.Context) {
	aul := &model.CreateCommentRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceCreateComment(id, aul, auth)
	ctx.Values().Set("response", response)
}

func CreateComment(ctx iris.Context) {
	aul := &model.CreateCommentRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateComment(id, aul, auth)
	ctx.Values().Set("response", response)
}

func DeleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceDeleteComment(id, auth)
	ctx.Values().Set("response", response)
}

func ForceDeleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteComment(id, auth)
	ctx.Values().Set("response", response)
}
