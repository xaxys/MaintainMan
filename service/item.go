package service

import (
	"errors"
	"fmt"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"

	"gorm.io/gorm"
)

func GetItemByID(id uint, auth *model.AuthInfo) *model.ApiJson {
	item, err := dao.GetItemByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(ItemToJson(item), "获取成功")
}

func GetItemByName(name string, auth *model.AuthInfo) *model.ApiJson {
	item, err := dao.GetItemByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(ItemToJson(item), "获取成功")
}

func GetItemsByFuzzyName(name string, auth *model.AuthInfo) *model.ApiJson {
	items, err := dao.GetItemsByFuzzyName(name)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	is := util.TransSlice(items, ItemToJson)
	return model.Success(is, "获取成功")
}

func GetAllItems(param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(param); err != nil {
		return model.ErrorValidation(err)
	}
	items, err := dao.GetAllItems(param)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	is := util.TransSlice(items, ItemToJson)
	return model.Success(is, "获取成功")
}

func CreateItem(aul *model.CreateItemRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	item, err := dao.CreateItem(aul, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(ItemToInfoJson(item), "创建成功")
}

func DeleteItem(id uint, auth *model.AuthInfo) *model.ApiJson {
	if err := dao.DeleteItem(id); err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func AddItem(aul *model.AddItemRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	itemlog := dao.ItemLogAdd(aul)
	log, err := dao.AddItem(itemlog, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessUpdate(ItemLogToJson(log), "添加成功")
}

func ConsumeItem(aul *model.ConsumeItemRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dao.GetOrderWithLastStatus(aul.OrderID)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status != model.StatusAssigned {
		return model.ErrorNoPermissions(fmt.Errorf("订单未处于已接单状态"))
	}
	if util.LastElem(order.StatusList).RepairerID != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的当前维修员"))
	}
	itemlog := dao.ItemLogConsume(aul)
	log, err := dao.ConsumeItem(itemlog, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessUpdate(ItemLogToJson(log), "添加成功")
}

func ItemToJson(item *model.Item) *model.ItemJson {
	return util.NilOrValue(item, &model.ItemJson{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Count:       item.Count,
	})
}

func ItemToInfoJson(item *model.Item) *model.ItemInfoJson {
	return util.NilOrValue(item, &model.ItemInfoJson{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Income:      item.Income,
		Count:       item.Count,
		ItemLogs:    util.TransSlice(item.ItemLogs, ItemLogToJson),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
		CreatedBy:   item.CreatedBy,
		UpdatedBy:   item.UpdatedBy,
	})
}

func ItemLogToJson(itemLog *model.ItemLog) *model.ItemLogJson {
	return util.NilOrValue(itemLog, &model.ItemLogJson{
		ID:          itemLog.ID,
		ItemID:      itemLog.ItemID,
		OrderID:     itemLog.OrderID,
		ChangeNum:   itemLog.ChangeNum,
		ChangePrice: itemLog.ChangePrice,
		CreatedAt:   itemLog.CreatedAt.Unix(),
		CreatedBy:   itemLog.CreatedBy,
	})
}
