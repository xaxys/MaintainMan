package user

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getDivision godoc
// @Summary      获取某分组信息
// @Description  通过ID获取某分组信息
// @Tags         division
// @Produce      json
// @Param        id   path      uint  true  "分组ID"
// @Success      200  {object}  model.ApiJson{data=DivisionJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/division/{id} [get]
func getDivision(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getDivisionService(id, auth)
	ctx.Values().Set("response", response)
}

// getDivisionsByParentID godoc
// @Summary      获取某分组下的子分组
// @Description  通过父分组ID获取某分组下的子分组
// @Tags         division
// @Produce      json
// @Param        id   path      uint  true  "父分组ID"
// @Success      200  {object}  model.ApiJson{data=[]DivisionJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/division/{id}/children [get]
func getDivisionsByParentID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getDivisionsByParentIDService(id, auth)
	ctx.Values().Set("response", response)
}

// createDivision godoc
// @Summary      创建分组
// @Description  创建分组
// @Tags         division
// @Accept       json
// @Produce      json
// @Param        body  body      CreateDivisionRequest  true  "创建分组请求"
// @Success      201   {object}  model.ApiJson{data=DivisionJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/division [post]
func createDivision(ctx iris.Context) {
	aul := &CreateDivisionRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := createDivisionService(aul, auth)
	ctx.Values().Set("response", response)
}

// updateDivision godoc
// @Summary      更新分组
// @Description  更新分组
// @Tags         division
// @Accept       json
// @Produce      json
// @Param        body  body      UpdateDivisionRequest  true  "更新分组请求"
// @Success      204   {object}  model.ApiJson{data=DivisionJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/division/{id} [put]
func updateDivision(ctx iris.Context) {
	aul := &UpdateDivisionRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := updateDivisionService(id, aul, auth)
	ctx.Values().Set("response", response)
}

// deleteDivision godoc
// @Summary      删除分组
// @Description  删除分组
// @Tags         division
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "分组ID"
// @Success      204  {object}  model.ApiJson
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/division/{id} [delete]
func deleteDivision(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := deleteDivisionService(id, auth)
	ctx.Values().Set("response", response)
}
