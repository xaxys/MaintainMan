package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetUserOrders godoc
// @Summary 获取当前用户的订单
// @Description 获取当前用户的订单 分页 默认逆序 可按照订单状态过滤
// @Description 状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价
// @Tags order
// @Produce json
// @Param status query int false "订单状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"
// @Param order_by query string false "排序字段 (默认为ID正序) 只接受"{field} {asc|desc}"格式 (e.g. "id desc")"
// @Param offset query uint false "偏移量 (默认为0)"
// @Param limit query uint false "每页数据量 (默认为50)"
// @Success 200 {object} model.ApiJson{data=[]model.OrderJson}  "返回结果 带Tag"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/user [get]
func GetUserOrders(ctx iris.Context) {
	req := &model.UserOrderRequest{}
	if err := ctx.ReadQuery(req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByUser(req, auth)
	ctx.Values().Set("response", response)
}

// GetRepairerOrders godoc
// @Summary 获取当前维修工的订单
// @Description 获取当前维修工的订单 分页 默认逆序 可按照是否本人正在维修过滤
// @Tags order
// @Produce json
// @Param current query bool true "是否本人正在维修"
// @Param order_by query string false "排序字段 (默认为ID正序) 只接受"{field} {asc|desc}"格式 (e.g. "id desc")"
// @Param offset query uint false "偏移量 (默认为0)"
// @Param limit query uint false "每页数据量 (默认为50)"
// @Success 200 {object} model.ApiJson{data=[]model.OrderJson}  "返回结果 带Tag"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/repairer [get]
func GetRepairerOrders(ctx iris.Context) {
	req := &model.RepairerOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByRepairer(auth.User, req, auth)
	ctx.Values().Set("response", response)
}

// ForceGetRepairerOrders godoc
// @Summary 获取某维修工的订单
// @Description 通过维修工ID获取某维修工的订单 分页 默认逆序 可按照是否该人正在维修过滤
// @Tags order
// @Produce json
// @Param id path uint true "维修工ID"
// @Param current query bool true "是否本人正在维修"
// @Param order_by query string false "排序字段 (默认为ID正序) 只接受"{field} {asc|desc}"格式 (e.g. "id desc")"
// @Param offset query uint false "偏移量 (默认为0)"
// @Param limit query uint false "每页数据量 (默认为50)"
// @Success 200 {object} model.ApiJson{data=[]model.OrderJson}  "返回结果 带Tag"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/repairer/{id} [get]
func ForceGetRepairerOrders(ctx iris.Context) {
	req := &model.RepairerOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByRepairer(id, req, auth)
	ctx.Values().Set("response", response)
}

// GetAllOrders godoc
// @Summary 获取所有订单
// @Description 获取所有订单 分页 默认正序 可按照 标题 用户 订单状态 多个Tag(与|或 两种模式)过滤
// @Description 状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价
// @Tags order
// @Produce json
// @Param title query string false "标题"
// @Param user_id query uint false "用户ID"
// @Param status query string false "订单状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"
// @Param tags query []string false "若干 Tag 的 ID"
// @Param conjunctve query bool false "true: 查询包含所有Tag的订单, false: 查询包含任一Tag的订单"
// @Param order_by query string false "排序字段 (默认为ID正序) 只接受"{field} {asc|desc}"格式 (e.g. "id desc")"
// @Param offset query uint false "偏移量 (默认为0)"
// @Param limit query uint false "每页数据量 (默认为50)"
// @Success 200 {object} model.ApiJson{data=[]model.OrderJson}  "返回结果 带Tag"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/all [get]
func GetAllOrders(ctx iris.Context) {
	req := &model.AllOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllOrders(req, auth)
	ctx.Values().Set("response", response)
}

// GetOrderByID GetOrder godoc
// @Summary 获取某个订单
// @Description 通过ID获取某个订单
// @Tags order
// @Produce json
// @Param id path uint true "订单ID"
// @Success 200 {object} model.ApiJson{data=model.OrderJson}  "返回结果 带Tag 带Comment"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id} [get]
func GetOrderByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByID(id, auth)
	ctx.Values().Set("response", response)
}

// CreateOrder godoc
// @Summary 创建订单
// @Description 创建订单
// @Tags order
// @Accept json
// @Produce json
// @Param body model.CreateOrderRequest true "请求参数"
// @Success 201 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order [post]
func CreateOrder(ctx iris.Context) {
	aul := &model.CreateOrderRequest{}

	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateOrder(aul, auth)
	ctx.Values().Set("response", response)
}

// UpdateOrder godoc
// @Summary 更新订单
// @Description 更新订单 操作者需为订单创建者
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Param body body model.UpdateOrderRequest true "请求参数"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id} [put]
func UpdateOrder(ctx iris.Context) {
	aul := &model.UpdateOrderRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateOrder(id, aul, auth)
	ctx.Values().Set("response", response)
}

// ForceUpdateOrder godoc
// @Summary 更新订单(管理员)
// @Description 更新订单(管理员)
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Param body body model.UpdateOrderRequest true "请求参数"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/force [put]
func ForceUpdateOrder(ctx iris.Context) {
	aul := &model.UpdateOrderRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceUpdateOrder(id, aul, auth)
	ctx.Values().Set("response", response)
}

// change order status

// ReleaseOrder godoc
// @Summary 释放订单
// @Description 释放订单 从 已接单 已完成 上报中 挂单 已拒绝 到 待处理
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/release [post]
func ReleaseOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ReleaseOrder(id, auth)
	ctx.Values().Set("response", response)
}

// AssignOrder godoc
// @Summary 指派订单
// @Description 指派订单 从 待处理 到 已接单
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Param repairer query uint true "维修工ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/assign [post]
func AssignOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	repairer := util.ToUint(ctx.URLParamIntDefault("repairer", 0))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AssignOrder(id, repairer, auth)
	ctx.Values().Set("response", response)
}

// SelfAssignOrder godoc
// @Summary 自指派订单
// @Description 自指派订单 从 待处理 到 已接单
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/selfassign [post]
func SelfAssignOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AssignOrder(id, auth.User, auth)
	ctx.Values().Set("response", response)
}

// CompleteOrder godoc
// @Summary 完成订单
// @Description 完成订单 从 已接单 到 已完成 操作者只能是当前维修工
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/complete [post]
func CompleteOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CompleteOrder(id, auth)
	ctx.Values().Set("response", response)
}

// CancelOrder godoc
// @Summary 取消订单
// @Description 取消订单 从 除已完成 已评价外的状态 到 已取消 操作者只能是订单创建者
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/cancel [post]
func CancelOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CancelOrder(id, auth)
	ctx.Values().Set("response", response)
}

// RejectOrder godoc
// @Summary 拒绝订单
// @Description 拒绝订单 从 待处理 到 已拒绝
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/reject [post]
func RejectOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.RejectOrder(id, auth)
	ctx.Values().Set("response", response)
}

// ReportOrder godoc
// @Summary 上报订单
// @Description 上报订单 从 已接单 到 上报中 操作者只能是当前维修工
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/report [post]
func ReportOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ReportOrder(id, auth)
	ctx.Values().Set("response", response)
}

// HoldOrder godoc
// @Summary 挂起订单
// @Description 挂起订单 从 待处理 到 挂单
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/hold [post]
func HoldOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.HoldOrder(id, auth)
	ctx.Values().Set("response", response)
}

// AppraiseOrder godoc
// @Summary 评价订单
// @Description 评价订单 从 已完成 到 已评价
// @Tags order
// @Accept json
// @Produce json
// @Param id path uint true "订单ID"
// @Param appraisal query uint true "评价分数"
// @Success 204 {object} model.ApiJson{data=model.OrderJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/order/{id}/appraise [post]
func AppraiseOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	appraisal := util.ToUint(ctx.URLParamIntDefault("appraisal", 0))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AppraiseOrder(id, appraisal, auth)
	ctx.Values().Set("response", response)
}
