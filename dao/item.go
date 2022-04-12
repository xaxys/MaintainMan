package dao

import (
	"fmt"
	"maintainman/config"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
)

func GetItemByID(id uint) (*model.Item, error) {
	return TxGetItemByID(database.DB, id)
}

func TxGetItemByID(tx *gorm.DB, id uint) (*model.Item, error) {
	item := &model.Item{}
	if err := tx.First(item, id).Error; err != nil {
		logger.Logger.Debugf("GetItemByIDErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func GetItemByName(name string) (*model.Item, error) {
	return TxGetItemByName(database.DB, name)
}

func TxGetItemByName(tx *gorm.DB, name string) (*model.Item, error) {
	item := &model.Item{Name: name}
	if err := tx.Where(item).First(item).Error; err != nil {
		logger.Logger.Debugf("GetItemByNameErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func GetItemsByFuzzyName(name string) (items []*model.Item, err error) {
	return TxGetItemsByFuzzyName(database.DB, name)
}

func TxGetItemsByFuzzyName(tx *gorm.DB, name string) (items []*model.Item, err error) {
	if err = TxFilter(tx, "", 0, 0).Where("name like (?)", name).Find(&items).Error; err != nil {
		logger.Logger.Debugf("GetItemByNameErr: %v\n", err)
		return nil, err
	}
	return
}

func GetAllItems(param *model.PageParam) (items []*model.Item, err error) {
	return TxGetAllItems(database.DB, param)
}

func TxGetAllItems(tx *gorm.DB, param *model.PageParam) (items []*model.Item, err error) {
	if err = TxPageFilter(tx, param).Find(&items).Error; err != nil {
		logger.Logger.Debugf("GetAllItemsErr: %v\n", err)
		return
	}
	return
}

func CreateItem(aul *model.CreateItemRequest, operator uint) (*model.Item, error) {
	return TxCreateItem(database.DB, aul, operator)
}

func TxCreateItem(tx *gorm.DB, aul *model.CreateItemRequest, operator uint) (*model.Item, error) {
	item := JsonToItem(aul)
	item.CreatedBy = operator
	if err := tx.Create(item).Error; err != nil {
		logger.Logger.Debugf("CreateItemErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func DeleteItem(id uint) error {
	return TxDeleteItem(database.DB, id)
}

func TxDeleteItem(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Item{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteItemErr: %v\n", err)
		return err
	}
	return nil
}

func AddItem(itemlog *model.ItemLog, operator uint) (item *model.Item, err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if item, err = TxAddItem(tx, itemlog, operator); err != nil {
			logger.Logger.Debugf("AddItemErr: %v\n", err)
		}
		return err
	})
	return
}

func TxAddItem(tx *gorm.DB, itemlog *model.ItemLog, operator uint) (item *model.Item, err error) {
	itemlog.CreatedBy = operator
	if item, err = GetItemByID(itemlog.ItemID); err != nil {
		return
	}
	item.Count += itemlog.ChangeNum
	item.Price += itemlog.ChangePrice
	item.UpdatedBy = operator
	if err = tx.Create(itemlog).Error; err != nil {
		return
	}
	if err = tx.Save(item).Error; err != nil {
		return
	}
	return
}

func ConsumeItem(itemlog *model.ItemLog, operator uint) (item *model.Item, err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if item, err = TxConsumeItem(tx, itemlog, operator); err != nil {
			logger.Logger.Debugf("ConsumeItemErr: %v\n", err)
		}
		return err
	})
	return
}

func TxConsumeItem(tx *gorm.DB, itemlog *model.ItemLog, operator uint) (item *model.Item, err error) {
	itemlog.CreatedBy = operator
	if item, err = GetItemByID(itemlog.ItemID); err != nil {
		return
	}
	if item.Count < itemlog.ChangeNum && !config.AppConfig.GetBool("app.item_can_negative") {
		return nil, fmt.Errorf("item count is not enough")
	}
	item.Count -= itemlog.ChangeNum
	item.Income += -itemlog.ChangePrice
	item.UpdatedBy = operator
	if err = tx.Create(itemlog).Error; err != nil {
		return
	}
	if err = tx.Save(item).Error; err != nil {
		return
	}
	return
}

func JsonToItem(item *model.CreateItemRequest) *model.Item {
	return &model.Item{
		Name:        item.Name,
		Description: item.Discription,
	}
}
