package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	offset, _ := ctx.URLParamInt("offset")
	uid := ctx.Values().GetUintDefault("user_id", 0)
	response := service.GetCommentsByOrder(id, util.ToUint(offset), uid)
	ctx.Values().Set("response", response)
}

func GetCommentsByOrderID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	offset, _ := ctx.URLParamInt("offset")
	response := service.GetCommentsByOrderID(id, util.ToUint(offset))
	ctx.Values().Set("response", response)
}

func CreateCommentOverride(ctx iris.Context) {
	oid := ctx.Params().GetUintDefault("user_id", 0)
	aul := &model.CreateCommentJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.CreateCommentOverride(oid, aul)
	ctx.Values().Set("response", response)
}

func CreateComment(ctx iris.Context) {
	oid := ctx.Params().GetUintDefault("user_id", 0)
	aul := &model.CreateCommentJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	response := service.CreateComment(oid, aul)
	ctx.Values().Set("response", response)
}

func DeleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.DeleteCommentByID(id)
	ctx.Values().Set("response", response)
}

func DeleteCommentByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	uid := ctx.Values().GetUintDefault("user_id", 0)
	response := service.DeleteComment(id, uid)
	ctx.Values().Set("response", response)
}
