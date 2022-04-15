package order

import (
	"github.com/xaxys/maintainman/core/controller"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getCommentsByOrder godoc
// @Summary 获取订单的评论信息
// @Description 获取订单的评论信息 分页 操作者必须是订单的创建者 或 曾经被分配给该订单的维修工
// @Tags comment
// @Produce  json
// @Param id path uint true "订单id"
// @Param order_by query string false "排序字段"
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=model.Page{entries=[]CommentJson}}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment [get]
func getCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	param := controller.ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getCommentsByOrderService(id, param, auth)
	ctx.Values().Set("response", response)
}

// forceGetCommentsByOrder godoc
// @Summary 获取订单的评论信息(管理员)
// @Description 获取任意订单的评论信息 分页
// @Tags comment
// @Produce  json
// @Param id path uint true "订单id"
// @Param order_by query string false "排序字段"
// @Param offset query uint false "偏移量"
// @Param limit query uint false "每页数据量"
// @Success 200 {object} model.ApiJson{data=model.Page{entries=[]CommentJson}}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment/force [get]
func forceGetCommentsByOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	param := controller.ExtractPageParam(ctx)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := forceGetCommentsByOrderService(id, param, auth)
	ctx.Values().Set("response", response)
}

// forceCreateComment godoc
// @Summary 创建评论(管理员)
// @Description 创建评论
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "订单id"
// @Param body body CreateCommentRequest true "评论信息"
// @Success 201 {object} model.ApiJson{data=CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment/force [post]
func forceCreateComment(ctx iris.Context) {
	aul := &CreateCommentRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := forceCreateCommentService(id, aul, auth)
	ctx.Values().Set("response", response)
}

// createComment godoc
// @Summary 创建评论
// @Description 创建评论 创建者必须是订单的创建者 或 曾经被分配给该订单的维修工
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "订单id"
// @Param body body CreateCommentRequest true "评论信息"
// @Success 201 {object} model.ApiJson{data=CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/comment [post]
func createComment(ctx iris.Context) {
	aul := &CreateCommentRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := createCommentService(id, aul, auth)
	ctx.Values().Set("response", response)
}

// deleteComment godoc
// @Summary 删除评论
// @Description 删除评论 删除者必须是评论的创建者
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "评论id"
// @Success 204 {object} model.ApiJson{data=CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/comment/{id} [delete]
func deleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := forceDeleteCommentService(id, auth)
	ctx.Values().Set("response", response)
}

// forceDeleteComment godoc
// @Summary 删除评论(管理员)
// @Description 删除评论
// @Tags comment
// @Accept  json
// @Produce  json
// @Param id path uint true "评论id"
// @Success 204 {object} model.ApiJson{data=CommentJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/comment/{id}/force [delete]
func forceDeleteComment(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := DeleteCommentService(id, auth)
	ctx.Values().Set("response", response)
}
