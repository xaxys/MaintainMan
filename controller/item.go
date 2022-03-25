package controller

import (
	"maintainman/model"
	"maintainman/service"

	"github.com/kataras/iris/v12"
)

func GetItemByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.GetItemByID(id)
	ctx.Values().Set("response", response)
}

func GetItemByName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	response := service.GetItemByName(name)
	ctx.Values().Set("response", response)
}

func GetItemsByFuzzyName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	response := service.GetItemsByFuzzyName(name)
	ctx.Values().Set("response", response)
}

func GetAllItems(ctx iris.Context) {
	aul := &model.AllItemJson{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.GetAllItems(aul)
	ctx.Values().Set("response", response)
}

func CreateItem(ctx iris.Context) {
	aul := &model.CreateItemJson{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.CreateItem(aul)
	ctx.Values().Set("response", response)
}

func DeleteItemByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.DeleteItem(id)
	ctx.Values().Set("response", response)
}

func AddItem(ctx iris.Context) {
	aul := &model.AddItemJson{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.AddItem(aul)
	ctx.Values().Set("response", response)
}

func ConsumeItem(ctx iris.Context) {
	aul := &model.ConsumeItemJson{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.ConsumeItem(aul)
	ctx.Values().Set("response", response)
}
