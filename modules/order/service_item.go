package order

import (
	"errors"
	"fmt"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func getItemByIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	item, err := dbGetItemByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(itemToJson(item), "获取成功")
}

func getItemByNameService(name string, auth *model.AuthInfo) *model.ApiJson {
	item, err := dbGetItemByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(itemToJson(item), "获取成功")
}

func getItemsByFuzzyNameService(name string, auth *model.AuthInfo) *model.ApiJson {
	items, err := dbGetItemsByFuzzyName(name)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	is := util.TransSlice(items, itemToJson)
	return model.Success(is, "获取成功")
}

func getAllItemsService(param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(param); err != nil {
		return model.ErrorValidation(err)
	}
	items, err := dbGetAllItems(param)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	is := util.TransSlice(items, itemToJson)
	return model.Success(is, "获取成功")
}

func createItemService(aul *CreateItemRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	item, err := dbCreateItem(aul, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(itemToInfoJson(item), "创建成功")
}

func deleteItemService(id uint, auth *model.AuthInfo) *model.ApiJson {
	if err := dbDeleteItem(id); err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func addItemService(aul *AddItemRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	itemlog := dbItemLogAdd(aul)
	log, err := dbAddItem(itemlog, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessUpdate(itemToJson(log), "添加成功")
}

func consumeItemService(aul *ConsumeItemRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	order, err := dbGetOrderWithLastStatus(aul.OrderID)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	if order.Status != StatusAssigned {
		return model.ErrorNoPermissions(fmt.Errorf("订单未处于已接单状态"))
	}
	if uint(util.LastElem(order.StatusList).RepairerID.Int64) != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的当前维修员"))
	}
	itemlog := dbItemLogConsume(aul)
	log, err := dbConsumeItem(itemlog, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessUpdate(itemToJson(log), "添加成功")
}

func itemToJson(item *Item) *ItemJson {
	if item == nil {
		return nil
	} else {
		return &ItemJson{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Count:       item.Count,
		}
	}

}

func itemToInfoJson(item *Item) *ItemInfoJson {
	if item == nil {
		return nil
	} else {
		return &ItemInfoJson{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Income:      item.Income,
			Count:       item.Count,
			ItemLogs:    util.TransSlice(item.ItemLogs, itemLogToJson),
			CreatedAt:   item.CreatedAt.Unix(),
			UpdatedAt:   item.UpdatedAt.Unix(),
			CreatedBy:   item.CreatedBy,
			UpdatedBy:   item.UpdatedBy,
		}
	}

}

func itemLogToJson(itemLog *ItemLog) *ItemLogJson {
	if itemLog == nil {
		return nil
	} else {
		return &ItemLogJson{
			ID:          itemLog.ID,
			ItemID:      itemLog.ItemID,
			OrderID:     itemLog.OrderID,
			ChangeNum:   itemLog.ChangeNum,
			ChangePrice: itemLog.ChangePrice,
			CreatedAt:   itemLog.CreatedAt.Unix(),
			CreatedBy:   itemLog.CreatedBy,
		}
	}

}
