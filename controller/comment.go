package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetCommentsByOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	offset, _ := ctx.URLParamInt("offset")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.GetCommentsByOrder(id, util.ToUint(offset), uid)
	ctx.Values().Set("response", response)
}

func GetCommentsByOrderID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	offset, _ := ctx.URLParamInt("offset")
	response := service.GetCommentsByOrderID(id, util.ToUint(offset))
	ctx.Values().Set("response", response)
}

func CreateCommentOverride(ctx iris.Context) {
	oid, _ := ctx.Params().GetUint("id")
	aul := &model.CreateCommentJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.CreateCommentOverride(oid, aul)
	ctx.Values().Set("response", response)
}

func CreateComment(ctx iris.Context) {
	oid, _ := ctx.Params().GetUint("id")
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
	id, _ := ctx.Params().GetUint("id")
	response := service.DeleteCommentByID(id)
	ctx.Values().Set("response", response)
}

func DeleteCommentByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.DeleteComment(id, uid)
	ctx.Values().Set("response", response)
}
