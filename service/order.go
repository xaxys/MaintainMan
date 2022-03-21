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
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dao.UpdateOrder(id, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(OrderToJson(order), "更新成功")
}

func DeleteOrder(id uint) *model.ApiJson {
	if err := dao.DeleteOrder(id); err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func ReleaseOrder(id, operator uint) *model.ApiJson {
	status := dao.StatusWaiting(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "释放成功")
}

func AssignOrder(id, uid, operator uint) *model.ApiJson {
	status := dao.StatusAssigned(uid, operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "指派成功")
}

func CompleteOrder(id, operator uint) *model.ApiJson {
	status := dao.StatusCompleted(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "结单成功")
}

func CancelOrder(id, operator uint) *model.ApiJson {
	status := dao.StatusCanceled(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "取消成功")
}

func RejectOrder(id, operator uint) *model.ApiJson {
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
	if order.Status != model.StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未完成，不能评价"))
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
	status := dao.StatusReported(operator)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "上报成功")
}

func HoldOrder(id, operator uint) *model.ApiJson {
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
