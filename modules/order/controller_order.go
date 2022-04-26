package order

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getUserOrders godoc
// @Summary      获取当前用户的订单
// @Description  获取当前用户的订单 分页 默认逆序 可按照订单状态过滤
// @Description  状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价
// @Tags         order
// @Produce      json
// @Param        tags        query     []string                                             false  "若干 Tag 的 ID"
// @Param        disjunctve  query     bool                                                 false  "false: 查询包含所有Tag的订单, true: 查询包含任一Tag的订单"
// @Param        status      query     int                                                  false  "订单状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"
// @Param        order_by    query     string                                               false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset      query     uint                                                 false  "偏移量 (默认为0)"
// @Param        limit       query     uint                                                 false  "每页数据量 (默认为50)"
// @Success      200         {object}  model.ApiJson{data=model.Page{entries=[]OrderJson}}  "返回结果 带Tag"
// @Failure      400         {object}  model.ApiJson{data=[]string}
// @Failure      401         {object}  model.ApiJson{data=[]string}
// @Failure      403         {object}  model.ApiJson{data=[]string}
// @Failure      404         {object}  model.ApiJson{data=[]string}
// @Failure      422         {object}  model.ApiJson{data=[]string}
// @Failure      500         {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/user [get]
func getUserOrders(ctx iris.Context) {
	req := &UserOrderRequest{}
	if err := ctx.ReadQuery(req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getOrderByUserService(req, auth)
	ctx.Values().Set("response", response)
}

// getRepairerOrders godoc
// @Summary      获取当前维修工的订单
// @Description  获取当前维修工的订单 分页 默认逆序 可按照是否本人正在维修过滤
// @Tags         order
// @Produce      json
// @Param        tags        query     []string                                             false  "若干 Tag 的 ID"
// @Param        disjunctve  query     bool                                                 false  "false: 查询包含所有Tag的订单, true: 查询包含任一Tag的订单"
// @Param        status      query     int                                                  false  "订单状态 0:所有 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"
// @Param        current     query     bool                                                 true   "是否本人正在维修"
// @Param        order_by    query     string                                               false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset      query     uint                                                 false  "偏移量 (默认为0)"
// @Param        limit       query     uint                                                 false  "每页数据量 (默认为50)"
// @Success      200         {object}  model.ApiJson{data=model.Page{entries=[]OrderJson}}  "返回结果 带Tag"
// @Failure      400         {object}  model.ApiJson{data=[]string}
// @Failure      401         {object}  model.ApiJson{data=[]string}
// @Failure      403         {object}  model.ApiJson{data=[]string}
// @Failure      404         {object}  model.ApiJson{data=[]string}
// @Failure      422         {object}  model.ApiJson{data=[]string}
// @Failure      500         {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/repairer [get]
func getRepairerOrders(ctx iris.Context) {
	req := &RepairerOrderRequest{}
	if err := ctx.ReadQuery(req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getOrderByRepairerService(auth.User, req, auth)
	ctx.Values().Set("response", response)
}

// forceGetRepairerOrders godoc
// @Summary      获取某维修工的订单 (管理员)
// @Description  通过维修工ID获取某维修工的订单 分页 默认逆序 可按照是否该人正在维修过滤
// @Tags         order
// @Produce      json
// @Param        id          path      uint                                                 true   "维修工ID"
// @Param        tags        query     []string                                             false  "若干 Tag 的 ID"
// @Param        disjunctve  query     bool                                                 false  "false: 查询包含所有Tag的订单, true: 查询包含任一Tag的订单"
// @Param        status      query     int                                                  false  "订单状态 0:所有 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"
// @Param        current     query     bool                                                 true   "是否本人正在维修"
// @Param        order_by    query     string                                               false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset      query     uint                                                 false  "偏移量 (默认为0)"
// @Param        limit       query     uint                                                 false  "每页数据量 (默认为50)"
// @Success      200         {object}  model.ApiJson{data=model.Page{entries=[]OrderJson}}  "返回结果 带Tag"
// @Failure      400         {object}  model.ApiJson{data=[]string}
// @Failure      401         {object}  model.ApiJson{data=[]string}
// @Failure      403         {object}  model.ApiJson{data=[]string}
// @Failure      404         {object}  model.ApiJson{data=[]string}
// @Failure      422         {object}  model.ApiJson{data=[]string}
// @Failure      500         {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/repairer/{id} [get]
func forceGetRepairerOrders(ctx iris.Context) {
	req := &RepairerOrderRequest{}
	if err := ctx.ReadQuery(req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getOrderByRepairerService(id, req, auth)
	ctx.Values().Set("response", response)
}

// getAllOrders godoc
// @Summary      获取所有订单
// @Description  获取所有订单 分页 默认正序 可按照 标题 用户 订单状态 多个Tag(与|或 两种模式)过滤
// @Description  状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价
// @Tags         order
// @Produce      json
// @Param        title       query     string                                               false  "标题"
// @Param        user_id     query     uint                                                 false  "用户ID"
// @Param        status      query     string                                               false  "订单状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"
// @Param        tags        query     []string                                             false  "若干 Tag 的 ID"
// @Param        disjunctve  query     bool                                                 false  "false: 查询包含所有Tag的订单, true: 查询包含任一Tag的订单"
// @Param        order_by    query     string                                               false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset      query     uint                                                 false  "偏移量 (默认为0)"
// @Param        limit       query     uint                                                 false  "每页数据量 (默认为50)"
// @Success      200         {object}  model.ApiJson{data=model.Page{entries=[]OrderJson}}  "返回结果 带Tag"
// @Failure      400         {object}  model.ApiJson{data=[]string}
// @Failure      401         {object}  model.ApiJson{data=[]string}
// @Failure      403         {object}  model.ApiJson{data=[]string}
// @Failure      404         {object}  model.ApiJson{data=[]string}
// @Failure      422         {object}  model.ApiJson{data=[]string}
// @Failure      500         {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/all [get]
func getAllOrders(ctx iris.Context) {
	req := &AllOrderRequest{}
	if err := ctx.ReadQuery(req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getAllOrdersService(req, auth)
	ctx.Values().Set("response", response)
}

// getOrderByID GetOrder godoc
// @Summary      获取某个订单
// @Description  通过ID获取某个订单
// @Tags         order
// @Produce      json
// @Param        id   path      uint                              true  "订单ID"
// @Success      200  {object}  model.ApiJson{data=OrderJson}  "返回结果 带Tag 带Comment 带Repairer"
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id} [get]
func getOrderByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getOrderByIDService(id, auth)
	ctx.Values().Set("response", response)
}

// forceGetOrderByID GetOrder godoc
// @Summary      获取某个订单 (管理员)
// @Description  通过ID获取某个订单 (管理员)
// @Tags         order
// @Produce      json
// @Param        id   path      uint                              true  "订单ID"
// @Success      200  {object}  model.ApiJson{data=OrderJson}  "返回结果 带Tag 带Comment 带Repairer"
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/force [get]
func forceGetOrderByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := forceGetOrderByIDService(id, auth)
	ctx.Values().Set("response", response)
}

// getOrderStatus godoc
// @Summary      获取订单所有历史状态
// @Description  获取订单所有历史状态
// @Tags         order
// @Produce      json
// @Param        id   path      uint                           true  "订单ID"
// @Success      200  {object}  model.ApiJson{data=[]StatusJson}  "返回结果"
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/status [get]
func getOrderStatus(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getOrderStatusService(id, auth)
	ctx.Values().Set("response", response)
}

// forceGetOrderStatus godoc
// @Summary      获取订单所有历史状态 (管理员)
// @Description  通过ID获取订单所有历史状态 (管理员)
// @Tags         order
// @Produce      json
// @Param        id   path      uint                           true  "订单ID"
// @Success      200  {object}  model.ApiJson{data=[]StatusJson}  "返回结果"
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/status/force [get]
func forceGetOrderStatus(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := forceGetOrderStatusService(id, auth)
	ctx.Values().Set("response", response)
}

// createOrder godoc
// @Summary      创建订单
// @Description  创建订单
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        body  body      CreateOrderRequest  true  "请求参数"
// @Success      201   {object}  model.ApiJson{data=OrderJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/order [post]
func createOrder(ctx iris.Context) {
	aul := &CreateOrderRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := createOrderService(aul, auth)
	ctx.Values().Set("response", response)
}

// updateOrder godoc
// @Summary      更新订单
// @Description  更新订单 操作者需为订单创建者
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id    path      uint                true  "订单ID"
// @Param        body  body      UpdateOrderRequest  true  "请求参数"
// @Success      204   {object}  model.ApiJson{data=OrderJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id} [put]
func updateOrder(ctx iris.Context) {
	aul := &UpdateOrderRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := updateOrderService(id, aul, auth)
	ctx.Values().Set("response", response)
}

// forceUpdateOrder godoc
// @Summary      更新订单(管理员)
// @Description  更新订单(管理员)
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id    path      uint                true  "订单ID"
// @Param        body  body      UpdateOrderRequest  true  "请求参数"
// @Success      204   {object}  model.ApiJson{data=OrderJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/force [put]
func forceUpdateOrder(ctx iris.Context) {
	aul := &UpdateOrderRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := forceUpdateOrderService(id, aul, auth)
	ctx.Values().Set("response", response)
}

// change order status

// releaseOrder godoc
// @Summary      释放订单
// @Description  释放订单 从 已接单 已完成 上报中 挂单 已拒绝 到 待处理
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/release [post]
func releaseOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := releaseOrderService(id, auth)
	ctx.Values().Set("response", response)
}

// assignOrder godoc
// @Summary      指派订单
// @Description  指派订单 从 待处理 到 已接单
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id        path      uint  true  "订单ID"
// @Param        repairer  query     uint  true  "维修工ID"
// @Success      204       {object}  model.ApiJson{data=OrderJson}
// @Failure      400       {object}  model.ApiJson{data=[]string}
// @Failure      401       {object}  model.ApiJson{data=[]string}
// @Failure      403       {object}  model.ApiJson{data=[]string}
// @Failure      404       {object}  model.ApiJson{data=[]string}
// @Failure      422       {object}  model.ApiJson{data=[]string}
// @Failure      500       {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/assign [post]
func assignOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	repairer := util.ToUint(ctx.URLParamIntDefault("repairer", 0))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := assignOrderService(id, repairer, auth)
	ctx.Values().Set("response", response)
}

// selfAssignOrder godoc
// @Summary      自指派订单
// @Description  自指派订单 从 待处理 到 已接单
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/selfassign [post]
func selfAssignOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := assignOrderService(id, auth.User, auth)
	ctx.Values().Set("response", response)
}

// completeOrder godoc
// @Summary      完成订单
// @Description  完成订单 从 已接单 到 已完成 操作者只能是当前维修工
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/complete [post]
func completeOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := completeOrderService(id, auth)
	ctx.Values().Set("response", response)
}

// cancelOrder godoc
// @Summary      取消订单
// @Description  取消订单 从 除已完成 已评价外的状态 到 已取消 操作者只能是订单创建者
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/cancel [post]
func cancelOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := cancelOrderService(id, auth)
	ctx.Values().Set("response", response)
}

// rejectOrder godoc
// @Summary      拒绝订单
// @Description  拒绝订单 从 待处理 到 已拒绝
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/reject [post]
func rejectOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := rejectOrderService(id, auth)
	ctx.Values().Set("response", response)
}

// reportOrder godoc
// @Summary      上报订单
// @Description  上报订单 从 已接单 到 上报中 操作者只能是当前维修工
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/report [post]
func reportOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := reportOrderService(id, auth)
	ctx.Values().Set("response", response)
}

// holdOrder godoc
// @Summary      挂起订单
// @Description  挂起订单 从 待处理 到 挂单
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "订单ID"
// @Success      204  {object}  model.ApiJson{data=OrderJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/hold [post]
func holdOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := holdOrderService(id, auth)
	ctx.Values().Set("response", response)
}

// appraiseOrder godoc
// @Summary      评价订单
// @Description  评价订单 从 已完成 到 已评价
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true  "订单ID"
// @Param        appraisal  query     uint  true  "评价分数"
// @Success      204        {object}  model.ApiJson{data=OrderJson}
// @Failure      400        {object}  model.ApiJson{data=[]string}
// @Failure      401        {object}  model.ApiJson{data=[]string}
// @Failure      403        {object}  model.ApiJson{data=[]string}
// @Failure      404        {object}  model.ApiJson{data=[]string}
// @Failure      422        {object}  model.ApiJson{data=[]string}
// @Failure      500        {object}  model.ApiJson{data=[]string}
// @Router       /v1/order/{id}/appraise [post]
func appraiseOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	appraisal := util.ToUint(ctx.URLParamIntDefault("appraisal", 0))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := appraiseOrderService(id, appraisal, auth)
	ctx.Values().Set("response", response)
}
