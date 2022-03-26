package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetItemByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetItemByID(id, auth)
	ctx.Values().Set("response", response)
}

func GetItemByName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetItemByName(name, auth)
	ctx.Values().Set("response", response)
}

func GetItemsByFuzzyName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetItemsByFuzzyName(name, auth)
	ctx.Values().Set("response", response)
}

func GetAllItems(ctx iris.Context) {
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllItems(param, auth)
	ctx.Values().Set("response", response)
}

func CreateItem(ctx iris.Context) {
	aul := &model.CreateItemRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateItem(aul, auth)
	ctx.Values().Set("response", response)
}

func DeleteItem(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteItem(id, auth)
	ctx.Values().Set("response", response)
}

func AddItem(ctx iris.Context) {
	aul := &model.AddItemRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.ItemID = ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AddItem(aul, auth)
	ctx.Values().Set("response", response)
}

func ConsumeItem(ctx iris.Context) {
	aul := &model.ConsumeItemRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OrderID = ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ConsumeItem(aul, auth)
	ctx.Values().Set("response", response)
}
