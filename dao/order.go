package dao

import (
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/util"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// GetSimpleOrderByID return no relative info
func GetSimpleOrderByID(id uint) (*model.Order, error) {
	return TxGetSimpleOrderByID(database.DB, id)
}

// TxGetSimpleOrderByID return no relative info
func TxGetSimpleOrderByID(tx *gorm.DB, id uint) (*model.Order, error) {
	order := &model.Order{}
	if err := tx.First(order, id).Error; err != nil {
		logger.Logger.Debugf("TxGetOrderByIDErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func GetOrderByID(id uint) (*model.Order, error) {
	return TxGetOrderByID(database.DB, id)
}

func TxGetOrderByID(tx *gorm.DB, id uint) (*model.Order, error) {
	order := &model.Order{}
	if err := tx.Preload("Tags").Preload("Comments").First(order, id).Error; err != nil {
		logger.Logger.Debugf("TxGetOrderByIDErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func GetAllOrdersWithParam(aul *model.AllOrderRequest) (orders []*model.Order, err error) {
	return TxGetAllOrdersWithParam(database.DB, aul)
}

func TxGetAllOrdersWithParam(tx *gorm.DB, aul *model.AllOrderRequest) (orders []*model.Order, err error) {
	order := &model.Order{
		UserID: aul.UserID,
		Status: aul.Status,
	}
	db := TxPageFilter(tx, &aul.PageParam).Preload("Tags").Where(order)
	if len(aul.Tags) > 0 {
		if aul.Conjunctve {
			for _, tag := range aul.Tags {
				db = db.Where("exist (?)", db.Table("order_tags").Select("order_id").Where("tag_id = (?)", tag).Where("order_id = order.id"))
			}
		} else {
			db = db.Where("id IN (?)", db.Table("order_tags").Select("order_id").Where("tag_id IN (?)", aul.Tags))
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
	return TxGetOrderWithLastStatus(database.DB, id)
}

func TxGetOrderWithLastStatus(tx *gorm.DB, id uint) (*model.Order, error) {
	order := &model.Order{}
	if err := tx.Preload("StatusList", "current = TRUE").Model(order).Find(order, id).Error; err != nil {
		logger.Logger.Debugf("GetOrderWithLastStatusErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func CreateOrder(aul *model.CreateOrderRequest, operator uint) (order *model.Order, err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if order, err = TxCreateOrder(tx, aul, operator); err != nil {
			logger.Logger.Debugf("CreateOrderErr: %v\n", err)
		}
		return err
	})
	return
}

func TxCreateOrder(tx *gorm.DB, aul *model.CreateOrderRequest, operator uint) (order *model.Order, err error) {
	order = &model.Order{}
	copier.Copy(order, aul)
	order.CreatedBy = operator
	order.UserID = operator
	tags, err := TxGetTagsByIDs(tx, aul.Tags)
	if err != nil {
		return
	}
	if err = CheckTagsCongener(tags); err != nil {
		return
	}
	if err = tx.Create(order).Error; err != nil {
		return
	}
	if err = tx.Model(order).Association("Tags").Append(tags); err != nil {
		return
	}
	status := StatusWaiting(operator)
	status.SequenceNum = 1
	if err = tx.Model(order).Association("StatusList").Append(status); err != nil {
		return
	}
	return
}

func UpdateOrder(id uint, aul *model.UpdateOrderRequest, operator uint) (order *model.Order, err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if order, err = TxUpdateOrder(tx, id, aul, operator); err != nil {
			logger.Logger.Debugf("UpdateOrderErr: %v\n", err)
		}
		return err
	})
	return
}

func TxUpdateOrder(tx *gorm.DB, id uint, aul *model.UpdateOrderRequest, operator uint) (order *model.Order, err error) {
	order = &model.Order{}
	copier.Copy(order, aul)
	order.ID = id
	order.UpdatedBy = operator

	addTags, err := TxGetTagsByIDs(tx, aul.AddTags)
	if err != nil {
		return
	}
	delTags, err := TxGetTagsByIDs(tx, aul.DelTags)
	if err != nil {
		return
	}

	if err = tx.Model(order).Updates(order).Error; err != nil {
		return
	}
	if err = tx.Model(order).Association("Tags").Append(addTags); err != nil {
		return
	}
	if err = tx.Model(order).Association("Tags").Delete(delTags); err != nil {
		return
	}
	if err = tx.Preload("Tags").First(order, id).Error; err != nil {
		return
	}
	if err = CheckTagsCongener(order.Tags); err != nil {
		return
	}
	return
}

func DeleteOrder(id uint) error {
	return TxDeleteOrder(database.DB, id)
}

func TxDeleteOrder(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Order{}, id).Error; err != nil {
		logger.Logger.Debugf("TxDeleteOrderErr: %v\n", err)
		return err
	}
	return nil
}

func ChangeOrderStatus(id uint, status *model.Status) (err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if err = TxChangeOrderStatus(tx, id, status); err != nil {
			logger.Logger.Debugf("ChangeOrderStatusErr: %v\n", err)
		}
		return err
	})
	return
}

func TxChangeOrderStatus(tx *gorm.DB, id uint, status *model.Status) error {
	order := &model.Order{}
	order.ID = id
	order.Status = status.Status
	order.UpdatedBy = status.CreatedBy

	if err := tx.Model(order).Updates(order).Error; err != nil {
		return err
	}
	or, err := TxGetOrderWithLastStatus(tx, id)
	if err != nil {
		return err
	}
	statusList := or.StatusList
	lastStatus := util.LastElem(or.StatusList)
	lastStatus.UpdatedBy = status.CreatedBy
	lastStatus.Current = false
	status.SequenceNum = lastStatus.SequenceNum + 1
	statusList = append(statusList, status)
	if err := tx.Model(order).Association("StatusList").Replace(statusList); err != nil {
		return err
	}
	return nil
}

func ChangeOrderAllowComment(id uint, allow bool) error {
	return TxChangeOrderAllowComment(database.DB, id, allow)
}

func TxChangeOrderAllowComment(tx *gorm.DB, id uint, allow bool) error {
	order := &model.Order{}
	order.ID = id
	order.AllowComment = util.Tenary[uint](allow, 1, 2)
	if err := tx.Model(order).Updates(order).Error; err != nil {
		logger.Logger.Debugf("TxChangeOrderAllowCommentErr: %v\n", err)
		return err
	}
	return nil
}

func AppraiseOrder(id, appraisal, operator uint) (err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if err := TxAppraiseOrder(tx, id, appraisal, operator); err != nil {
			logger.Logger.Debugf("AppraiseOrderErr: %v\n", err)
		}
		return err
	})
	return
}

func TxAppraiseOrder(tx *gorm.DB, id, appraisal, operator uint) (err error) {
	order := &model.Order{}
	order.ID = id
	order.Appraisal = appraisal
	order.UpdatedBy = order.UserID
	order.AllowComment = model.CommentDisallow

	if err = tx.Model(order).Updates(order).Error; err != nil {
		return
	}
	status := StatusAppraised(operator)
	if err = TxChangeOrderStatus(tx, id, status); err != nil {
		return
	}
	return
}
