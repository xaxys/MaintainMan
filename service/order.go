package service

import (
	"errors"
	"fmt"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"

	"gorm.io/gorm"
)

func GetOrderByID(id uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(OrderToJson(order), "获取成功")
}

func GetOrderByUser(id, status, offset uint) *model.ApiJson {
	orders, err := dao.GetOrderByUser(id, status, offset)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	os := util.TransSlice(orders, OrderToJson)
	return model.Success(os, "获取成功")
}

func GetOrderByRepairer(id uint, current bool, offset uint) *model.ApiJson {
	orders, err := dao.GetOrderByRepairer(id, current, offset)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	os := util.TransSlice(orders, OrderToJson)
	return model.Success(os, "获取成功")
}

func GetAllOrders(aul *model.AllOrderJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	orders, err := dao.GetAllOrdersWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	os := util.TransSlice(orders, OrderToJson)
	return model.Success(os, "获取成功")
}

func CreateOrder(aul *model.ModifyOrderJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dao.CreateOrder(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(OrderToJson(order), "创建成功")
}

func UpdateOrder(id uint, aul *model.ModifyOrderJson) *model.ApiJson {
	or, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if or.UserID != aul.OperatorID {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单创建者"))
	}
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dao.UpdateOrder(id, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(OrderToJson(order), "更新成功")
}

func UpdateOrderByID(id uint, aul *model.ModifyOrderJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dao.UpdateOrder(id, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(OrderToJson(order), "更新成功")
}

// TODO: 订单应不允许删除?
// func DeleteOrder(id uint) *model.ApiJson {
// 	if err := dao.DeleteOrder(id); err != nil {
// 		return model.ErrorDeleteDatabase(err)
// 	}
// 	return model.SuccessUpdate(nil, "删除成功")
// }

func ReleaseOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于待处理状态"))
	}
	if util.In(order.Status, model.StatusAppraised, model.StatusCanceled) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已结束，不能再次维修"))
	}
	status := dao.StatusWaiting(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "释放成功")
}

func AssignOrder(id, uid, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已接单状态"))
	}
	if order.Status != model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能指派"))
	}
	status := dao.StatusAssigned(uid, operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "指派成功")
}

func SelfAssignOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已接单状态"))
	}
	if order.Status != model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能自我指派"))
	}
	status := dao.StatusAssigned(operator, operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "自我指派成功")
}

func CompleteOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已完成状态"))
	}
	if order.Status != model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于已指派状态，不能完成"))
	}
	if order.UserID != operator {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单当前指派人"))
	}
	status := dao.StatusCompleted(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "结单成功")
}

func CancelOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusCanceled {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已取消状态"))
	}
	if util.In(order.Status, model.StatusCompleted, model.StatusAppraised) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已完成，不能取消"))
	}
	status := dao.StatusCanceled(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "取消成功")
}

func RejectOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusRejected {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已拒绝状态"))
	}
	if order.Status != model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能拒绝"))
	}
	status := dao.StatusRejected(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "拒绝成功")
}

func AppraiseOrder(id, appraisal, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusAppraised {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已评价状态"))
	}
	if order.Status != model.StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未完成，不能评价"))
	}
	if order.UserID != operator {
		return model.ErrorUpdateDatabase(fmt.Errorf("您不是订单的创建者，不能评价"))
	}
	status := dao.StatusAppraised(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	if err := dao.AppraiseOrder(id, appraisal); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "评价成功")
}

func ReportOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusReported {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已上报状态"))
	}
	if order.Status != model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未指派，不能上报"))
	}
	if order.StatusList[len(order.StatusList)-1].RepairerID != operator {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单指派人，不能上报"))
	}
	status := dao.StatusReported(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "上报成功")
}

func HoldOrder(id, operator uint) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusHold {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于挂单状态"))
	}
	if !util.In(order.Status, model.StatusReported, model.StatusWaiting) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理或已上报状态，不能挂单"))
	}
	status := dao.StatusHold(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "挂单成功")
}

func OrderToJson(order *model.Order) *model.OrderJson {
	return &model.OrderJson{
		ID:           order.ID,
		UserID:       order.UserID,
		User:         UserToJson(order.User),
		Title:        order.Title,
		Content:      order.Content,
		Address:      order.Address,
		ContactName:  order.ContactName,
		ContactPhone: order.ContactPhone,
		Status:       order.Status,
		CreatedAt:    order.CreatedAt.Unix(),
		UpdatedAt:    order.UpdatedAt.Unix(),
		Appraisal:    order.Appraisal,
		Tags:         util.TransSlice(order.Tags, TagToJson),
	}
}
