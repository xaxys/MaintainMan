package dao

import (
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
)

func GetItemByID(id uint) (*model.Item, error) {
	item := &model.Item{}
	if err := database.DB.First(item, id).Error; err != nil {
		logger.Logger.Debugf("GetItemByIDErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func GetItemByName(name string) (*model.Item, error) {
	item := &model.Item{Name: name}
	if err := database.DB.Where(item).First(item).Error; err != nil {
		logger.Logger.Debugf("GetItemByNameErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func GetItemsByFuzzyName(name string) (items []*model.Item, err error) {
	if err = Filter("", 0, 0).Where("name like (?)", name).Find(&items).Error; err != nil {
		logger.Logger.Debugf("GetItemByNameErr: %v\n", err)
		return nil, err
	}
	return
}

func GetAllItems(param *model.PageParam) (items []*model.Item, err error) {
	if err = PageFilter(param).Find(&items).Error; err != nil {
		logger.Logger.Debugf("GetAllItemsErr: %v\n", err)
		return
	}
	return
}

func CreateItem(aul *model.CreateItemRequest, operator uint) (*model.Item, error) {
	item := JsonToItem(aul)
	item.CreatedBy = operator
	if err := database.DB.Create(item).Error; err != nil {
		logger.Logger.Debugf("CreateItemErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func DeleteItem(id uint) error {
	if err := database.DB.Delete(&model.Item{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteItemErr: %v\n", err)
		return err
	}
	return nil
}

func AddItem(itemlog *model.ItemLog, operator uint) (*model.ItemLog, error) {
	itemlog.CreatedBy = operator
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		item, err := GetItemByID(itemlog.ItemID)
		if err != nil {
			return err
		}
		item.Count += itemlog.ChangeNum
		item.Price += itemlog.ChangePrice
		item.UpdatedBy = operator
		if err := tx.Create(itemlog).Error; err != nil {
			return err
		}
		if err := tx.Save(item).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("AddItemErr: %v\n", err)
		return nil, err
	}
	return itemlog, nil
}

func ConsumeItem(itemlog *model.ItemLog, operator uint) (*model.ItemLog, error) {
	itemlog.CreatedBy = operator
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		item, err := GetItemByID(itemlog.ItemID)
		if err != nil {
			return err
		}
		// TODO: 考虑到实际情况，是否应该允许消耗的数量大于库存数量?
		// if item.Count < itemlog.ChangeNum {
		// 	return fmt.Errorf("item count is not enough")
		// }
		item.Count -= itemlog.ChangeNum
		item.Income += -itemlog.ChangePrice
		item.UpdatedBy = operator
		if err := tx.Create(itemlog).Error; err != nil {
			return err
		}
		if err := tx.Save(item).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("ConsumeItemErr: %v\n", err)
		return nil, err
	}
	return itemlog, nil
}

func JsonToItem(item *model.CreateItemRequest) *model.Item {
	return &model.Item{
		Name:        item.Name,
		Description: item.Discription,
	}
}
