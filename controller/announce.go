package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetAnnounce godoc
// @Summary 获取公告
// @Description 获取公告
// @Tags announce
// @Produce  json
// @Param id path uint true "公告ID"
// @Success 200 {object} model.ApiJson{data=model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/{id} [get]
func GetAnnounce(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAnnounce(id, auth)
	ctx.Values().Set("response", response)
}

// GetAllAnnounces godoc
// @Summary 获取公告列表
// @Description 获取公告列表 分页 可按标题 开始时间 结束时间 (时间之内|之外 两种模式)过滤
// @Tags announce
// @Produce  json
// @Param title query string false "标题"
// @Param start_time query string false "开始时间; unix timestamp in seconds (UTC); -1代表不限; 含本数"
// @Param end_time query string false "结束时间; unix timestamp in seconds (UTC); -1代表不限; 含本数"
// @Param inclusive query bool false "true: 查询开始时间晚于start,且结束时间早于end的(在某段时间内开始并结束的); false: 查询开始时间早于start,且结束时间晚于end的(在某段时间内都能看到的)"
// @Param order_by query string false "排序字段 (默认为ID正序) 只接受"{field} {asc|desc}"格式 (e.g. "id desc")"
// @Param offset query uint false "偏移量 (默认为0)"
// @Param limit query uint false "每页数据量 (默认为50)"
// @Success 200 {object} model.ApiJson{data=[]model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/all [get]
func GetAllAnnounces(ctx iris.Context) {
	aul := &model.AllAnnounceRequest{}
	if err := ctx.ReadQuery(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllAnnounces(aul, auth)
	ctx.Values().Set("response", response)
}

// GetLatestAnnounce godoc
// @Summary 获取最新公告
// @Description 获取最新公告 分页 强制逆序 开始时间 结束时间 之内
// @Tags announce
// @Produce  json
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=[]model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/ [get]
func GetLatestAnnounces(ctx iris.Context) {
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetLatestAnnounces(param, auth)
	ctx.Values().Set("response", response)
}

// CreateAnnounce godoc
// @Summary 创建公告
// @Description 创建公告
// @Tags announce
// @Accept  json
// @Produce  json
// @Param body body model.CreateAnnounceRequest true "创建公告请求"
// @Success 201 {object} model.ApiJson{data=model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/ [post]
func CreateAnnounce(ctx iris.Context) {
	aul := &model.CreateAnnounceRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateAnnounce(aul, auth)
	ctx.Values().Set("response", response)
}

// UpdateAnnounce godoc
// @Summary 更新公告
// @Description 更新公告
// @Tags announce
// @Accept  json
// @Produce  json
// @Param id path uint true "公告ID"
// @Param body body model.UpdateAnnounceRequest true "更新公告请求"
// @Success 204 {object} model.ApiJson{data=model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/{id} [put]
func UpdateAnnounce(ctx iris.Context) {
	aul := &model.UpdateAnnounceRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateAnnounce(id, aul, auth)
	ctx.Values().Set("response", response)
}

// DeleteAnnounce godoc
// @Summary 删除公告
// @Description 删除公告
// @Tags announce
// @Accept  json
// @Produce  json
// @Param id path uint true "公告ID"
// @Success 204 {object} model.ApiJson{data=model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/{id} [delete]
func DeleteAnnounce(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteAnnounce(id, auth)
	ctx.Values().Set("response", response)
}

// HitAnnounce godoc
// @Summary 点击公告
// @Description 点击公告 增加点击量 默认单个用户单篇文章12h只能点击一次
// @Tags announce
// @Produce  json
// @Param id path uint true "公告ID"
// @Success 204 {object} model.ApiJson{data=model.AnnounceJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/announce/{id}/hit [put]
func HitAnnounce(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.HitAnnounce(id, auth)
	ctx.Values().Set("response", response)
}
