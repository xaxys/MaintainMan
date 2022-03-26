package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetAnnounceByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAnnounceByID(id, auth)
	ctx.Values().Set("response", response)
}

func GetAllAnnounces(ctx iris.Context) {
	aul := &model.AllAnnounceRequest{}
	if err := ctx.ReadQuery(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllAnnounces(aul, auth)
	ctx.Values().Set("response", response)
}

func GetLatestAnnounces(ctx iris.Context) {
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetLatestAnnounces(param, auth)
	ctx.Values().Set("response", response)
}

func CreateAnnounce(ctx iris.Context) {
	aul := &model.CreateAnnounceRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateAnnounce(aul, auth)
	ctx.Values().Set("response", response)
}

func UpdateAnnounce(ctx iris.Context) {
	aul := &model.UpdateAnnounceRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateAnnounce(id, aul, auth)
	ctx.Values().Set("response", response)
}

func DeleteAnnounce(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteAnnounce(id, auth)
	ctx.Values().Set("response", response)
}

func HitAnnounce(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.HitAnnounce(id, auth)
	ctx.Values().Set("response", response)
}
