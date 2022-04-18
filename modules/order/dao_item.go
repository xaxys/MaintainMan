package order

import (
	"fmt"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"

	"gorm.io/gorm"
)

func dbGetItemCount() (uint, error) {
	return txGetItemCount(mctx.Database)
}

func txGetItemCount(tx *gorm.DB) (uint, error) {
	count := int64(0)
	if err := tx.Model(&Item{}).Count(&count).Error; err != nil {
		mctx.Logger.Warnf("GetItemCountErr: %v\n", err)
		return 0, err
	}
	return uint(count), nil
}

func dbGetItemByID(id uint) (*Item, error) {
	return txGetItemByID(mctx.Database, id)
}

func txGetItemByID(tx *gorm.DB, id uint) (*Item, error) {
	item := &Item{}
	if err := tx.First(item, id).Error; err != nil {
		mctx.Logger.Warnf("GetItemByIDErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func dbGetItemByName(name string) (*Item, error) {
	return txGetItemByName(mctx.Database, name)
}

func txGetItemByName(tx *gorm.DB, name string) (*Item, error) {
	item := &Item{Name: name}
	if err := tx.Where(item).First(item).Error; err != nil {
		mctx.Logger.Warnf("GetItemByNameErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func dbGetItemsByFuzzyName(name string) (items []*Item, err error) {
	return TxGetItemsByFuzzyName(mctx.Database, name)
}

func TxGetItemsByFuzzyName(tx *gorm.DB, name string) (items []*Item, err error) {
	if err = dao.TxFilter(tx, "", 0, 0).Where("name like (?)", name).Find(&items).Error; err != nil {
		mctx.Logger.Warnf("GetItemByNameErr: %v\n", err)
		return nil, err
	}
	return
}

func dbGetAllItems(param *model.PageParam) (items []*Item, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if items, count, err = txGetAllItems(tx, param); err != nil {
			mctx.Logger.Warnf("GetAllItemsErr: %v\n", err)
		}
		return err
	})
	return
}

func txGetAllItems(tx *gorm.DB, param *model.PageParam) (items []*Item, count uint, err error) {
	tx = dao.TxPageFilter(tx, param)
	if err = tx.Find(&items).Error; err != nil {
		return
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}

func dbCreateItem(aul *CreateItemRequest, operator uint) (*Item, error) {
	return TxCreateItem(mctx.Database, aul, operator)
}

func TxCreateItem(tx *gorm.DB, aul *CreateItemRequest, operator uint) (*Item, error) {
	item := jsonToItem(aul)
	item.CreatedBy = operator
	if err := tx.Create(item).Error; err != nil {
		mctx.Logger.Warnf("CreateItemErr: %v\n", err)
		return nil, err
	}
	return item, nil
}

func dbDeleteItem(id uint) error {
	return TxDeleteItem(mctx.Database, id)
}

func TxDeleteItem(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&Item{}, id).Error; err != nil {
		mctx.Logger.Warnf("DeleteItemErr: %v\n", err)
		return err
	}
	return nil
}

func dbAddItem(itemlog *ItemLog, operator uint) (item *Item, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if item, err = txAddItem(tx, itemlog, operator); err != nil {
			mctx.Logger.Warnf("AddItemErr: %v\n", err)
		}
		return err
	})
	return
}

func txAddItem(tx *gorm.DB, itemlog *ItemLog, operator uint) (item *Item, err error) {
	itemlog.CreatedBy = operator
	if item, err = dbGetItemByID(itemlog.ItemID); err != nil {
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

func dbConsumeItem(itemlog *ItemLog, operator uint) (item *Item, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if item, err = txConsumeItem(tx, itemlog, operator); err != nil {
			mctx.Logger.Warnf("ConsumeItemErr: %v\n", err)
		}
		return err
	})
	return
}

func txConsumeItem(tx *gorm.DB, itemlog *ItemLog, operator uint) (item *Item, err error) {
	itemlog.CreatedBy = operator
	if item, err = dbGetItemByID(itemlog.ItemID); err != nil {
		return
	}
	if item.Count < itemlog.ChangeNum && !orderConfig.GetBool("item_can_negative") {
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

func jsonToItem(item *CreateItemRequest) *Item {
	return &Item{
		Name:        item.Name,
		Description: item.Discription,
	}
}
