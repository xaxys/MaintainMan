package order

import (
	"errors"
	"fmt"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"
	"github.com/xaxys/maintainman/modules/user"

	"gorm.io/gorm"
)

func getOrderByIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || order.ID == 0 {
		return model.ErrorNotFound(err)
	}
	if order.UserID != auth.User {
		repairer := util.LastElem(order.StatusList).RepairerID
		if repairer != nil && *repairer != auth.User {
			return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能查看评论"))
		}
	}
	return forceGetOrderByIDService(id, auth)
}

func forceGetOrderByIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || order.ID == 0 {
		return model.ErrorNotFound(err)
	}
	json := orderToJson(order)
	if rid := util.LastElem(order.StatusList).RepairerID; rid != nil {
		repairer, err := user.GetUserByID(*rid)
		if err != nil {
			mctx.Logger.Warnf("获取订单%d的指派人%d失败: %+v", id, *rid, err)
		}
		json.Repairer = user.UserToJson(repairer)
	}
	return model.Success(json, "获取成功")
}

func getOrderByUserService(aul *UserOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	aul.OrderBy = util.NotEmpty(aul.OrderBy, "id desc")
	allreq := &AllOrderRequest{
		UserID:    auth.User,
		Status:    aul.Status,
		Tags:      aul.Tags,
		PageParam: aul.PageParam,
	}
	return getAllOrdersService(allreq, auth)
}

func getOrderByRepairerService(id uint, aul *RepairerOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	aul.OrderBy = util.NotEmpty(aul.OrderBy, "id desc")
	orders, count, err := dbGetOrderByRepairer(id, aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	os := util.TransSlice(orders, orderToJson)
	return model.SuccessPaged(os, count, "获取成功")
}

func getOrderStatusService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || order.ID == 0 {
		return model.ErrorNotFound(err)
	}
	if order.UserID != auth.User {
		repairer := util.LastElem(order.StatusList).RepairerID
		if repairer != nil && *repairer != auth.User {
			return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能查看评论"))
		}
	}
	return forceGetOrderStatusService(id, auth)
}

func forceGetOrderStatusService(id uint, auth *model.AuthInfo) *model.ApiJson {
	statuses, err := dbGetStatusByOrder(id)
	if err != nil || len(statuses) == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	ss := util.TransSlice(statuses, statusToJson)
	return model.Success(ss, "获取成功")
}

func getAllOrdersService(aul *AllOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	orders, count, err := dbGetAllOrdersWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	os := util.TransSlice(orders, orderToJson)
	return model.SuccessPaged(os, count, "获取成功")
}

func createOrderService(aul *CreateOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
	if errResp := checkTagsService(aul.Tags, "tag.view", role); errResp != nil {
		return errResp
	}
	order, err := dbCreateOrder(aul, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	go mctx.EventBus.Emit("order:create", order.ID)
	return model.SuccessCreate(orderToJson(order), "创建成功")
}

func updateOrderService(id uint, aul *UpdateOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.UserID != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单创建者"))
	}
	role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
	if errResp := checkTagsService(aul.AddTags, "tag.add", role); errResp != nil {
		return errResp
	}
	if errResp := checkTagsService(aul.DelTags, "tag.add", role); errResp != nil {
		return errResp
	}
	return forceUpdateOrderService(id, aul, auth)
}

func forceUpdateOrderService(id uint, aul *UpdateOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dbUpdateOrder(id, aul, auth.User)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	fields := util.NotEmptyFieldName(aul)
	for _, field := range fields {
		event := fmt.Sprintf("order:update:%s", field)
		go mctx.EventBus.Emit(event, order.ID)
	}
	return model.SuccessUpdate(orderToJson(order), "更新成功")
}

func releaseOrderService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于待处理状态"))
	}
	if util.In(order.Status, StatusAppraised, StatusCanceled) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已结束，不能再次维修"))
	}
	status := NewStatusWaiting(auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:waiting", order.ID, StatusWaiting)
	return model.SuccessUpdate(nil, "释放成功")
}

func assignOrderService(id, repairer uint, auth *model.AuthInfo) *model.ApiJson {
	if repairer == 0 {
		return model.ErrorUpdateDatabase(fmt.Errorf("维修人不能为空"))
	}
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已接单状态"))
	}
	if order.Status != StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能指派"))
	}
	status := NewStatusAssigned(repairer, auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:assigned", order.ID, StatusAssigned, repairer)
	return model.SuccessUpdate(nil, "指派成功")
}

func completeOrderService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已完成状态"))
	}
	if order.Status != StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于已指派状态，不能完成"))
	}
	if order.UserID != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单当前指派人"))
	}
	status := NewStatusCompleted(auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:completed", order.ID, StatusCompleted)
	return model.SuccessUpdate(nil, "结单成功")
}

func cancelOrderService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusCanceled {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已取消状态"))
	}
	if util.In(order.Status, StatusCompleted, StatusAppraised) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已完成，不能取消"))
	}
	status := NewStatusCanceled(auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:canceled", order.ID, StatusCanceled)
	return model.SuccessUpdate(nil, "取消成功")
}

func rejectOrderService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusRejected {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已拒绝状态"))
	}
	if order.Status != StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能拒绝"))
	}
	status := NewStatusRejected(auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:rejected", order.ID, StatusRejected)
	return model.SuccessUpdate(nil, "拒绝成功")
}

func appraiseOrderService(id, appraisal uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusAppraised {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已评价状态"))
	}
	if order.Status != StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未完成，不能评价"))
	}
	if order.UserID != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("您不是订单的创建者，不能评价"))
	}
	if err := dbAppraiseOrder(id, appraisal, auth.User); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:appraised", order.ID, StatusAppraised)
	return model.SuccessUpdate(nil, "评价成功")
}

func reportOrderService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetOrderWithLastStatus(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusReported {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已上报状态"))
	}
	if order.Status != StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未指派，不能上报"))
	}
	repairer := util.LastElem(order.StatusList).RepairerID
	if repairer != nil && *repairer != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单指派人，不能上报"))
	}
	status := NewStatusReported(auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:reported", order.ID, StatusReported)
	return model.SuccessUpdate(nil, "上报成功")
}

func holdOrderService(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == StatusHold {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于挂单状态"))
	}
	if !util.In(order.Status, StatusReported, StatusWaiting) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理或已上报状态，不能挂单"))
	}
	status := NewStatusHold(auth.User)
	if err := dbChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	go mctx.EventBus.Emit("order:update:status:hold", order.ID, StatusHold)
	return model.SuccessUpdate(nil, "挂单成功")
}

func autoAppraiseOrderService() {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		orders, err := txGetAppraiseTimeoutOrder(tx)
		if err != nil {
			return err
		}
		for _, order := range orders {
			def := util.ToUint(orderConfig.GetInt("appraise.default"))
			_ = dbAppraiseOrder(order, def, 0)
			go mctx.EventBus.Emit("order:update:status:appraised", order, StatusAppraised)
		}
		return nil
	})
}

func orderToJson(order *Order) *OrderJson {
	return &OrderJson{
		ID:           order.ID,
		UserID:       order.UserID,
		Title:        order.Title,
		Content:      order.Content,
		Address:      order.Address,
		ContactName:  order.ContactName,
		ContactPhone: order.ContactPhone,
		Status:       order.Status,
		CreatedAt:    order.CreatedAt.Unix(),
		UpdatedAt:    order.UpdatedAt.Unix(),
		Appraisal:    order.Appraisal,
		Tags:         util.TransSlice(order.Tags, tagToJson),
		AllowComment: order.AllowComment == CommentAllow,
		Comments:     util.TransSlice(order.Comments, commentToJson),
	}
}

func statusToJson(status *Status) *StatusJson {
	return &StatusJson{
		Status:      status.Status,
		SequenceNum: status.SequenceNum,
		RepairerID:  util.NilOrBaseValue(status.RepairerID, func(t *uint) uint { return *t }, 0),
		Repairer:    user.UserToJson(status.Repairer),
		CreatedAt:   status.CreatedAt.Unix(),
		CreatedBy:   status.CreatedBy,
		UpdatedAt:   status.UpdatedAt.Unix(),
		UpdatedBy:   status.UpdatedBy,
	}
}
