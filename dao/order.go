package dao

import (
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
)

func GetOrderByID(id uint) (*model.Order, error) {
	order := &model.Order{}

	if err := database.DB.Preload("Tags").First(order, id).Error; err != nil {
		logger.Logger.Debugf("GetOrderByIDErr: %v\n", err)
		return nil, err
	}

	return order, nil
}

func GetOrderByUser(userID, status, offset uint) (*model.Order, error) {
	order := &model.Order{
		UserID: userID, Status: status}

	if err := Filter("id", offset, 0).Where(order).First(order).Error; err != nil {
		logger.Logger.Debugf("GetOrderByUserErr: %v\n", err)
		return nil, err
	}

	return order, nil
}

func GetAllOrdersWithParam(aul *model.AllOrderJson) (orders []*model.Order, err error) {
	order := &model.Order{
		UserID: aul.UserID,
		Status: aul.Status,
	}
	db := Filter(aul.OrderBy, aul.Offset, aul.Limit).Where(order)
	if len(aul.Tags) > 0 {
		if aul.Conjunctve {
			for _, tag := range aul.Tags {
				db = db.Where("exist (?)", db.Table("order_tags").Select("order_id").Where("tag_id = ?", tag).Where("order_id = order.id"))
			}
		} else {
			db = db.Where("id IN (?)", db.Table("order_tags").Select("order_id").Where("tag_id IN ?", aul.Tags))
		}
	}
	if aul.Title != "" {
		db = db.Where("title like ?", aul.Title)
	}
	if err = db.Find(&orders).Error; err != nil {
		logger.Logger.Debugf("GetAllOrdersErr: %v\n", err)
	}
	return
}

func CreateOrder(aul *model.ModifyOrderJson) (*model.Order, error) {
	order := JsonToOrder(aul)
	order.CreatedBy = aul.OperatorID

	tags, err := GetTagsByIDs(aul.AddTags)
	if err != nil {
		logger.Logger.Debugf("CreateOrderErr: %v\n", err)
		return nil, err
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Create(order).Error; err != nil {
			return err
		}
		if err := database.DB.Model(order).Association("Tags").Append(tags); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("CreateOrderErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func UpdateOrder(id uint, aul *model.ModifyOrderJson) (*model.Order, error) {
	order := JsonToOrder(aul)
	order.ID = id
	order.UpdatedBy = aul.OperatorID

	addTags, err := GetTagsByIDs(aul.AddTags)
	if err != nil {
		logger.Logger.Debugf("UpdateOrderErr: %v\n", err)
		return nil, err
	}
	delTags, err := GetTagsByIDs(aul.DelTags)
	if err != nil {
		logger.Logger.Debugf("UpdateOrderErr: %v\n", err)
		return nil, err
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Model(order).Updates(order).Error; err != nil {
			return err
		}
		if err := database.DB.Model(order).Association("Tags").Append(addTags); err != nil {
			return err
		}
		if err := database.DB.Model(order).Association("Tags").Delete(delTags); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("UpdateOrderErr: %v\n", err)
		return nil, err
	}

	return order, err
}

func DeleteOrder(id uint) error {
	order := &model.Order{}
	order.ID = id

	if err := database.DB.Delete(order).Error; err != nil {
		logger.Logger.Debugf("DeleteOrderErr: %v\n", err)
		return err
	}

	return nil
}

func ChangeOrderStatus(id uint, status *model.Status) error {
	order := &model.Order{}
	order.ID = id
	order.Status = status.Status
	order.UpdatedBy = status.CreatedBy

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Model(order).Updates(order).Error; err != nil {
			return err
		}
		if err := database.DB.Model(order).Association("StatusList").Append(status); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("ChangeOrderStatusErr: %v\n", err)
		return err
	}
	return nil
}

func ChangeOrderAllowComment(id uint, allow bool) error {
	order := &model.Order{}
	order.ID = id
	if allow {
		order.AllowComment = 1
	} else {
		order.AllowComment = 2
	}

	if err := database.DB.Model(order).Updates(order).Error; err != nil {
		logger.Logger.Debugf("ChangeOrderAllowCommentErr: %v\n", err)
		return err
	}
	return nil
}

func AddOrderItemLog(id uint, itemLog *model.ItemLog) error {
	order := &model.Order{}
	order.ID = id
	if err := database.DB.Model(order).Association("ItemLogs").Append(itemLog); err != nil {
		logger.Logger.Debugf("AppendOrderItemLogErr: %v\n", err)
		return err
	}
	return nil
}

func JsonToOrder(json *model.ModifyOrderJson) *model.Order {
	order := &model.Order{
		Title:        json.Title,
		Content:      json.Content,
		Address:      json.Address,
		ContactName:  json.ContactName,
		ContactPhone: json.ContactPhone,
	}
	return order
}
