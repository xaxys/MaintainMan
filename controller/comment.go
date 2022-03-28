package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetCommentsByOrder godoc
// @Summary 获取订单的评论信息
// @Description 获取订单的评论信息 分页 操作者必须是订单的创建者 或 曾经被分配给该订单的维修工
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "订单id"
// @Param order_by query string false "排序字段"
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=[]model.CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment [get]
func GetCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetCommentsByOrder(id, param, auth)
	ctx.Values().Set("response", response)
}

// ForceGetCommentsByOrder godoc
// @Summary 获取订单的评论信息(管理员)
// @Description 获取任意订单的评论信息 分页
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "订单id"
// @Param order_by query string false "排序字段"
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=[]model.CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment/force [get]
func ForceGetCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	param := ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceGetCommentsByOrder(id, param, auth)
	ctx.Values().Set("response", response)
}

// ForceCreateComment godoc
// @Summary 创建评论(管理员)
// @Description 创建评论
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "订单id"
// @Param body body model.CreateCommentRequest true "评论信息"
// @Success 201 {object} model.ApiJson{data=model.CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment [post]
func ForceCreateComment(ctx iris.Context) {
	aul := &model.CreateCommentRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceCreateComment(id, aul, auth)
	ctx.Values().Set("response", response)
}

// CreateComment godoc
// @Summary 创建评论
// @Description 创建评论 创建者必须是订单的创建者 或 曾经被分配给该订单的维修工
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "订单id"
// @Param body body model.CreateCommentRequest true "评论信息"
// @Success 201 {object} model.ApiJson{data=model.CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment [post]
func CreateComment(ctx iris.Context) {
	aul := &model.CreateCommentRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateComment(id, aul, auth)
	ctx.Values().Set("response", response)
}

// DeleteComment godoc
// @Summary 删除评论
// @Description 删除评论 删除者必须是评论的创建者
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "评论id"
// @Success 204 {object} model.ApiJson{data=model.CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/comment/{id} [delete]
func DeleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceDeleteComment(id, auth)
	ctx.Values().Set("response", response)
}

// ForceDeleteComment godoc
// @Summary 删除评论(管理员)
// @Description 删除评论
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "评论id"
// @Success 204 {object} model.ApiJson{data=model.CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/comment/{id}/force [delete]
func ForceDeleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteComment(id, auth)
	ctx.Values().Set("response", response)
}
