package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetAnnounceByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
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
	id, _ := ctx.Params().GetUint("id")
	response := service.UpdateAnnounce(id, aul)
	ctx.Values().Set("response", response)
}

func DeleteAnnounceByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	response := service.DeleteAnnounce(id)
	ctx.Values().Set("response", response)
}

func HitAnnounceByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.HitAnnounce(id, uid)
	ctx.Values().Set("response", response)
}
