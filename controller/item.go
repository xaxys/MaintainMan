package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetItemByID godoc
// @Summary 获取某ID物品信息
// @Description 通过ID获取某物品信息
// @Tags item
// @Produce  json
// @Param id path uint true "物品ID"
// @Success 200 {object} model.ApiJson{data=model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{id} [get]
func GetItemByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetItemByID(id, auth)
	ctx.Values().Set("response", response)
}

// GetItemByName godoc
// @Summary 获取某名称物品信息
// @Description 通过名称获取某物品信息
// @Tags item
// @Produce  json
// @Param name path string true "物品名称"
// @Success 200 {object} model.ApiJson{data=model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{name} [get]
func GetItemByName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetItemByName(name, auth)
	ctx.Values().Set("response", response)
}

// GetItemsByFuzzyName godoc
// @Summary 获取大概是某些名称的物品们的信息
// @Description 通过名称获取大概是某些名称的物品们的信息
// @Tags item
// @Produce  json
// @Param name path string true "物品名称"
// @Success 200 {object} model.ApiJson{data=[]model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{name}/fuzzy [get]
func GetItemsByFuzzyName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetItemsByFuzzyName(name, auth)
	ctx.Values().Set("response", response)
}

// GetAllItems godoc
// @Summary 获取所有物品信息
// @Description 获取所有物品信息 分页
// @Tags item
// @Produce  json
// @Param order_by query string false "排序字段"
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=[]model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item [get]
func GetAllItems(ctx iris.Context) {
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllItems(param, auth)
	ctx.Values().Set("response", response)
}

// CreateItem godoc
// @Summary 创建物品
// @Description 创建物品
// @Tags item
// @Accept  json
// @Produce  json
// @Param item body model.CreateItemRequest true "物品信息"
// @Success 201 {object} model.ApiJson{data=model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item [post]
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

// DeleteItem godoc
// @Summary 删除物品
// @Description 删除物品
// @Tags item
// @Accept  json
// @Produce  json
// @Param id path uint true "物品ID"
// @Success 204 {object} model.ApiJson{data=model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{id} [delete]
func DeleteItem(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteItem(id, auth)
	ctx.Values().Set("response", response)
}

// AddItem godoc
// @Summary 添加物品数量(进货)
// @Description 添加物品数量(进货)
// @Tags item
// @Accept  json
// @Produce  json
// @Param id path uint true "物品ID"
// @Param body body model.AddItemRequest true "物品数量"
// @Success 204 {object} model.ApiJson{data=model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{id} [post]
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

// ConsumeItem godoc
// @Summary 消耗物品数量(订单消耗)
// @Description 消耗物品数量(订单消耗)
// @Tags item
// @Accept  json
// @Produce  json
// @Param id path uint true "订单ID"
// @Param body body model.ConsumeItemRequest true "物品数量"
// @Success 204 {object} model.ApiJson{data=model.ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/consume [post]
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
