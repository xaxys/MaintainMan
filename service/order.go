package service

import (
	"errors"
	"fmt"
	"maintainman/config"
	"maintainman/dao"
	"maintainman/database"
	"maintainman/model"
	"maintainman/util"

	"gorm.io/gorm"
)

func GetOrderByID(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(OrderToJson(order), "获取成功")
}

func GetOrderByUser(aul *model.UserOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	aul.OrderBy = util.NotEmpty(aul.OrderBy, "id desc")
	allreq := &model.AllOrderRequest{
		UserID:    auth.User,
		Status:    aul.Status,
		Tags:      aul.Tags,
		PageParam: aul.PageParam,
	}
	return GetAllOrders(allreq, auth)
}

func GetOrderByRepairer(id uint, aul *model.RepairerOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	aul.OrderBy = util.NotEmpty(aul.OrderBy, "id desc")
	orders, err := dao.GetOrderByRepairer(id, aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	os := util.TransSlice(orders, OrderToJson)
	return model.Success(os, "获取成功")
}

func GetAllOrders(aul *model.AllOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	orders, err := dao.GetAllOrdersWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	os := util.TransSlice(orders, OrderToJson)
	return model.Success(os, "获取成功")
}

func CreateOrder(aul *model.CreateOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
	if errResp := CheckTags(aul.Tags, "tag.view", role); errResp != nil {
		return errResp
	}
	order, err := dao.CreateOrder(aul, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(OrderToJson(order), "创建成功")
}

func UpdateOrder(id uint, aul *model.UpdateOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
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
	if errResp := CheckTags(aul.AddTags, "tag.add", role); errResp != nil {
		return errResp
	}
	if errResp := CheckTags(aul.DelTags, "tag.add", role); errResp != nil {
		return errResp
	}
	return ForceUpdateOrder(id, aul, auth)
}

func ForceUpdateOrder(id uint, aul *model.UpdateOrderRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dao.UpdateOrder(id, aul, auth.User)
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

func ReleaseOrder(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于待处理状态"))
	}
	if util.In(order.Status, model.StatusAppraised, model.StatusCanceled) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已结束，不能再次维修"))
	}
	status := dao.StatusWaiting(auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "释放成功")
}

func AssignOrder(id, repairer uint, auth *model.AuthInfo) *model.ApiJson {
	if repairer == 0 {
		return model.ErrorUpdateDatabase(fmt.Errorf("维修人不能为空"))
	}
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已接单状态"))
	}
	if order.Status != model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能指派"))
	}
	status := dao.StatusAssigned(repairer, auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "指派成功")
}

func CompleteOrder(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已完成状态"))
	}
	if order.Status != model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于已指派状态，不能完成"))
	}
	if order.UserID != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单当前指派人"))
	}
	status := dao.StatusCompleted(auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "结单成功")
}

func CancelOrder(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusCanceled {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已取消状态"))
	}
	if util.In(order.Status, model.StatusCompleted, model.StatusAppraised) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已完成，不能取消"))
	}
	status := dao.StatusCanceled(auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "取消成功")
}

func RejectOrder(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusRejected {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已拒绝状态"))
	}
	if order.Status != model.StatusWaiting {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理状态，不能拒绝"))
	}
	status := dao.StatusRejected(auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "拒绝成功")
}

func AppraiseOrder(id, appraisal uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusAppraised {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已评价状态"))
	}
	if order.Status != model.StatusCompleted {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未完成，不能评价"))
	}
	if order.UserID != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("您不是订单的创建者，不能评价"))
	}
	if err := dao.AppraiseOrder(id, appraisal, auth.User); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "评价成功")
}

func ReportOrder(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetOrderWithLastStatus(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusReported {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于已上报状态"))
	}
	if order.Status != model.StatusAssigned {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单未指派，不能上报"))
	}
	if uint(util.LastElem(order.StatusList).RepairerID.Int64) != auth.User {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是订单指派人，不能上报"))
	}
	status := dao.StatusReported(auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "上报成功")
}

func HoldOrder(id uint, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetSimpleOrderByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if order.Status == model.StatusHold {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单已处于挂单状态"))
	}
	if !util.In(order.Status, model.StatusReported, model.StatusWaiting) {
		return model.ErrorUpdateDatabase(fmt.Errorf("订单不处于待处理或已上报状态，不能挂单"))
	}
	status := dao.StatusHold(auth.User)
	if err := dao.ChangeOrderStatus(id, status); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "挂单成功")
}

func AutoAppraiseOrder() {
	database.DB.Transaction(func(tx *gorm.DB) error {
		orders, err := dao.TxGetAppraiseTimeoutOrder(tx)
		if err != nil {
			return err
		}
		for _, order := range orders {
			def := util.ToUint(config.AppConfig.GetInt("app.appraise.default"))
			_ = dao.AppraiseOrder(order, def, 0)
		}
		return nil
	})
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
		AllowComment: order.AllowComment == model.CommentAllow,
		Comments:     util.TransSlice(order.Comments, CommentToJson),
	}
}
