package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetTagByID godoc
// @Summary 获取某标签信息
// @Description 通过ID获取某标签信息
// @Tags tag
// @Produce  json
// @Param id path uint true "标签ID"
// @Success 200 {object} model.ApiJson{data=model.TagJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/tag/{id} [get]
func GetTagByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetTagByID(id, auth)
	ctx.Values().Set("response", response)
}

// GetAllTagSorts godoc
// @Summary 获取所有标签分类
// @Description 获取所有标签分类
// @Tags tag
// @Produce  json
// @Success 200 {object} model.ApiJson{data=[]string}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/tag/sort [get]
func GetAllTagSorts(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllTagSorts(auth)
	ctx.Values().Set("response", response)
}

// GetAllTagsBySort godoc
// @Summary 获取某分类下的所有标签
// @Description 通过分类名获取某分类下的所有标签
// @Tags tag
// @Produce  json
// @Param name path string true "分类名"
// @Success 200 {object} model.ApiJson{data=[]model.TagJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/tag/sort/{name} [get]
func GetAllTagsBySort(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllTagsBySort(name, auth)
	ctx.Values().Set("response", response)
}

// CreateTag godoc
// @Summary 创建标签
// @Description 创建标签
// @Tags tag
// @Accept  json
// @Produce  json
// @Param body body model.CreateTagRequest true "创建标签请求"
// @Success 201 {object} model.ApiJson{data=model.TagJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/tag [post]
func CreateTag(ctx iris.Context) {
	aul := &model.CreateTagRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateTag(aul, auth)
	ctx.Values().Set("response", response)
}

// DeleteTag godoc
// @Summary 删除标签
// @Description 通过ID删除标签
// @Tags tag
// @Accept  json
// @Produce  json
// @Param id path uint true "标签ID"
// @Success 204 {object} model.ApiJson{data=model.TagJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/tag/{id} [delete]
func DeleteTag(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteTag(id, auth)
	ctx.Values().Set("response", response)
}
