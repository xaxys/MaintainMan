package order

import (
	"github.com/xaxys/maintainman/core/controller"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getItemByID godoc
// @Summary 获取某ID物品信息
// @Description 通过ID获取某物品信息
// @Tags item
// @Produce  json
// @Param id path uint true "物品ID"
// @Success 200 {object} model.ApiJson{data=ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{id} [get]
func getItemByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getItemByIDService(id, auth)
	ctx.Values().Set("response", response)
}

// getItemByName godoc
// @Summary 获取某名称物品信息
// @Description 通过名称获取某物品信息
// @Tags item
// @Produce  json
// @Param name path string true "物品名称"
// @Success 200 {object} model.ApiJson{data=ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{name} [get]
func getItemByName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getItemByNameService(name, auth)
	ctx.Values().Set("response", response)
}

// getItemsByFuzzyName godoc
// @Summary 获取大概是某些名称的物品们的信息
// @Description 通过名称获取大概是某些名称的物品们的信息
// @Tags item
// @Produce  json
// @Param name path string true "物品名称"
// @Success 200 {object} model.ApiJson{data=[]ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{name}/fuzzy [get]
func getItemsByFuzzyName(ctx iris.Context) {
	name := ctx.Params().Get("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getItemsByFuzzyNameService(name, auth)
	ctx.Values().Set("response", response)
}

// getAllItems godoc
// @Summary 获取所有物品信息
// @Description 获取所有物品信息 分页
// @Tags item
// @Produce  json
// @Param order_by query string false "排序字段"
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=[]ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item [get]
func getAllItems(ctx iris.Context) {
	param := controller.ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getAllItemsService(param, auth)
	ctx.Values().Set("response", response)
}

// createItem godoc
// @Summary 创建物品
// @Description 创建物品
// @Tags item
// @Accept  json
// @Produce  json
// @Param item body CreateItemRequest true "物品信息"
// @Success 201 {object} model.ApiJson{data=ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item [post]
func createItem(ctx iris.Context) {
	aul := &CreateItemRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := createItemService(aul, auth)
	ctx.Values().Set("response", response)
}

// deleteItem godoc
// @Summary 删除物品
// @Description 删除物品
// @Tags item
// @Accept  json
// @Produce  json
// @Param id path uint true "物品ID"
// @Success 204 {object} model.ApiJson{data=ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{id} [delete]
func deleteItem(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := deleteItemService(id, auth)
	ctx.Values().Set("response", response)
}

// addItem godoc
// @Summary 添加物品数量(进货)
// @Description 添加物品数量(进货)
// @Tags item
// @Accept  json
// @Produce  json
// @Param id path uint true "物品ID"
// @Param body body AddItemRequest true "物品数量"
// @Success 204 {object} model.ApiJson{data=ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/item/{id} [post]
func addItem(ctx iris.Context) {
	aul := &AddItemRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.ItemID = ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := addItemService(aul, auth)
	ctx.Values().Set("response", response)
}

// consumeItem godoc
// @Summary 消耗物品数量(订单消耗)
// @Description 消耗物品数量(订单消耗)
// @Tags item
// @Accept  json
// @Produce  json
// @Param id path uint true "订单ID"
// @Param body body ConsumeItemRequest true "物品数量"
// @Success 204 {object} model.ApiJson{data=ItemJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/consume [post]
func consumeItem(ctx iris.Context) {
	aul := &ConsumeItemRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OrderID = ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := consumeItemService(aul, auth)
	ctx.Values().Set("response", response)
}
