package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetAnnounceByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.GetAnnounceByID(id)
	ctx.Values().Set("response", response)
}

func GetAllAnnounces(ctx iris.Context) {
	aul := &model.AllAnnounceJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.GetAllAnnounces(aul)
	ctx.Values().Set("response", response)
}

func GetLatestAnnounces(ctx iris.Context) {
	offset, _ := ctx.URLParamInt("offset")
	response := service.GetLatestAnnounces(util.ToUint(offset))
	ctx.Values().Set("response", response)
}

func CreateAnnounceByID(ctx iris.Context) {
	aul := &model.ModifyAnnounceJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	response := service.CreateAnnounce(aul)
	ctx.Values().Set("response", response)
}

func UpdateAnnounceByID(ctx iris.Context) {
	aul := &model.ModifyAnnounceJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.UpdateAnnounce(id, aul)
	ctx.Values().Set("response", response)
}

func DeleteAnnounceByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.DeleteAnnounce(id)
	ctx.Values().Set("response", response)
}

func HitAnnounceByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	uid := ctx.Values().GetUintDefault("user_id", 0)
	response := service.HitAnnounce(id, uid)
	ctx.Values().Set("response", response)
}
