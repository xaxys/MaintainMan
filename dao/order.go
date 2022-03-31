package dao

import (
	"fmt"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/util"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetOrderByID(id uint) (*model.Order, error) {
	order := &model.Order{}

	if err := database.DB.Preload("Tags").Preload("Comments").First(order, id).Error; err != nil {
		logger.Logger.Debugf("GetOrderByIDErr: %v\n", err)
		return nil, err
	}

	return order, nil
}

func GetAllOrdersWithParam(aul *model.AllOrderRequest) (orders []*model.Order, err error) {
	order := &model.Order{
		UserID: aul.UserID,
		Status: aul.Status,
	}
	db := Filter(aul.OrderBy, aul.Offset, aul.Limit).Preload("Tags").Where(order)
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

func GetOrderWithLastStatus(id uint) (*model.Order, error) {
	order := &model.Order{}
	order.ID = id
	if err := database.DB.Preload("StatusList", "current = TRUE").Find(order).Error; err != nil {
		logger.Logger.Debugf("GetOrderWithLastStatusErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func CreateOrder(aul *model.CreateOrderRequest, operator uint) (*model.Order, error) {
	order := &model.Order{}
	copier.Copy(order, aul)
	order.CreatedBy = operator

	//XXX: 此处增加将auth中的用户id赋值给user_id的逻辑
	order.UserID = operator
	tags, errs := GetTagsByIDs(aul.Tags)
	if len(errs) > 0 {
		logger.Logger.Debugf("CreateOrderErr: %v\n", errs)
		return nil, fmt.Errorf("Several tags not found")
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if err := tx.Model(order).Association("Tags").Append(tags); err != nil {
			return err
		}
		status := StatusWaiting(operator)
		status.SequenceNum = 1
		if err := tx.Model(order).Association("StatusList").Append(status); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("CreateOrderErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func UpdateOrder(id uint, aul *model.UpdateOrderRequest, operator uint) (*model.Order, error) {
	order := &model.Order{}
	copier.Copy(order, aul)
	order.ID = id
	order.UpdatedBy = operator

	addTags, errs := GetTagsByIDs(aul.AddTags)
	if len(errs) > 0 {
		logger.Logger.Debugf("UpdateOrderErr: %v\n", errs)
		return nil, fmt.Errorf("Several tags not found")
	}
	delTags, errs := GetTagsByIDs(aul.DelTags)
	if len(errs) > 0 {
		logger.Logger.Debugf("UpdateOrderErr: %v\n", errs)
		return nil, fmt.Errorf("Several tags not found")
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(order).Updates(order).Error; err != nil {
			return err
		}
		if err := tx.Model(order).Association("Tags").Append(addTags); err != nil {
			return err
		}
		if err := tx.Model(order).Association("Tags").Delete(delTags); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("UpdateOrderErr: %v\n", err)
		return nil, err
	}

	return order, nil
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
		if err := tx.Model(order).Updates(order).Error; err != nil {
			return err
		}
		or := &model.Order{}
		if err := tx.Preload("StatusList").Find(or, order.ID).Error; err != nil {
			return err
		}
		statusList := or.StatusList
		lastStatus := statusList[len(or.StatusList)-1]
		lastStatus.UpdatedBy = status.CreatedBy
		lastStatus.Current = false
		status.SequenceNum = lastStatus.SequenceNum + 1
		statusList = append(statusList, status)
		if err := tx.Model(order).Association("StatusList").Replace(statusList); err != nil {
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
	order.AllowComment = util.Tenary[uint](allow, 1, 2)

	if err := database.DB.Model(order).Updates(order).Error; err != nil {
		logger.Logger.Debugf("ChangeOrderAllowCommentErr: %v\n", err)
		return err
	}
	return nil
}

func AppraiseOrder(id, appraisal uint) error {
	order := &model.Order{}
	order.ID = id
	order.Appraisal = appraisal
	order.UpdatedBy = order.UserID

	if err := database.DB.Model(order).Updates(order).Error; err != nil {
		logger.Logger.Debugf("AppraiseOrderErr: %v\n", err)
		return err
	}

	return ChangeOrderStatus(id, StatusAppraised(order.UserID))
}
