package order

import (
	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/util"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// GetSimpleOrderByID return no relative info
func dbGetSimpleOrderByID(id uint) (*Order, error) {
	return txGetSimpleOrderByID(mctx.Database, id)
}

// txGetSimpleOrderByID return no relative info
func txGetSimpleOrderByID(tx *gorm.DB, id uint) (*Order, error) {
	order := &Order{}
	if err := tx.First(order, id).Error; err != nil {
		logger.Logger.Debugf("TxGetOrderByIDErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func dbGetOrderByID(id uint) (*Order, error) {
	return txGetOrderByID(mctx.Database, id)
}

func txGetOrderByID(tx *gorm.DB, id uint) (*Order, error) {
	order := &Order{}
	if err := tx.Preload("Tags").Preload("Comments").First(order, id).Error; err != nil {
		logger.Logger.Debugf("TxGetOrderByIDErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func dbGetAllOrdersWithParam(aul *AllOrderRequest) (orders []*Order, err error) {
	return txGetAllOrdersWithParam(mctx.Database, aul)
}

func txGetAllOrdersWithParam(tx *gorm.DB, aul *AllOrderRequest) (orders []*Order, err error) {
	order := &Order{
		UserID: aul.UserID,
		Status: aul.Status,
	}
	tx = dao.TxPageFilter(tx, &aul.PageParam).Preload("Tags").Where(order)
	if len(aul.Tags) > 0 {
		if aul.Disjunctive {
			tx = tx.Where("id IN (?)", database.DB.Table("order_tags").Select("order_id").Where("tag_id IN (?)", aul.Tags))
		} else {
			for _, tag := range aul.Tags {
				tx = tx.Where("EXISTS (?)", database.DB.Table("order_tags").Select("order_id").Where("tag_id = (?)", tag).Where("order_id = orders.id"))
			}
		}
	}
	if aul.Title != "" {
		tx = tx.Where("title LIKE ?", aul.Title)
	}
	if err = tx.Find(&orders).Error; err != nil {
		logger.Logger.Debugf("GetAllOrdersErr: %v\n", err)
	}
	return
}

func dbGetOrderWithLastStatus(id uint) (*Order, error) {
	return txGetOrderWithLastStatus(mctx.Database, id)
}

func txGetOrderWithLastStatus(tx *gorm.DB, id uint) (*Order, error) {
	order := &Order{}
	if err := tx.Preload("StatusList", "current = TRUE").Model(order).Find(order, id).Error; err != nil {
		logger.Logger.Debugf("GetOrderWithLastStatusErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func dbCreateOrder(aul *CreateOrderRequest, operator uint) (order *Order, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if order, err = txCreateOrder(tx, aul, operator); err != nil {
			logger.Logger.Debugf("CreateOrderErr: %v\n", err)
		}
		return err
	})
	return
}

func txCreateOrder(tx *gorm.DB, aul *CreateOrderRequest, operator uint) (order *Order, err error) {
	order = &Order{}
	copier.Copy(order, aul)
	order.CreatedBy = operator
	order.UserID = operator
	order.Status = StatusWaiting
	tags, err := txGetTagsByIDs(tx, aul.Tags)
	if err != nil {
		return
	}
	if err = dbCheckTagsCongener(tags); err != nil {
		return
	}
	if err = tx.Create(order).Error; err != nil {
		return
	}
	if err = tx.Model(order).Association("Tags").Append(tags); err != nil {
		return
	}
	status := NewStatusWaiting(operator)
	status.SequenceNum = 1
	if err = tx.Model(order).Association("StatusList").Append(status); err != nil {
		return
	}
	return
}

func dbUpdateOrder(id uint, aul *UpdateOrderRequest, operator uint) (order *Order, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if order, err = TxUpdateOrder(tx, id, aul, operator); err != nil {
			logger.Logger.Debugf("UpdateOrderErr: %v\n", err)
		}
		return err
	})
	return
}

func TxUpdateOrder(tx *gorm.DB, id uint, aul *UpdateOrderRequest, operator uint) (order *Order, err error) {
	order = &Order{}
	copier.Copy(order, aul)
	order.ID = id
	order.UpdatedBy = operator

	addTags, err := txGetTagsByIDs(tx, aul.AddTags)
	if err != nil {
		return
	}
	delTags, err := txGetTagsByIDs(tx, aul.DelTags)
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
	if err = dbCheckTagsCongener(order.Tags); err != nil {
		return
	}
	return
}

func dbDeleteOrder(id uint) error {
	return txDeleteOrder(mctx.Database, id)
}

func txDeleteOrder(tx *gorm.DB, id uint) error {
	if err := tx.Delete(Order{}, id).Error; err != nil {
		logger.Logger.Debugf("TxDeleteOrderErr: %v\n", err)
		return err
	}
	return nil
}

func dbChangeOrderStatus(id uint, status *Status) (err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if err = txChangeOrderStatus(tx, id, status); err != nil {
			logger.Logger.Debugf("ChangeOrderStatusErr: %v\n", err)
		}
		return err
	})
	return
}

func txChangeOrderStatus(tx *gorm.DB, id uint, status *Status) error {
	order := &Order{}
	order.ID = id
	order.Status = status.Status
	order.UpdatedBy = status.CreatedBy

	if err := tx.Model(order).Updates(order).Error; err != nil {
		return err
	}
	or, err := txGetOrderWithLastStatus(tx, id)
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

func dbChangeOrderAllowComment(id uint, allow bool) error {
	return txChangeOrderAllowComment(mctx.Database, id, allow)
}

func txChangeOrderAllowComment(tx *gorm.DB, id uint, allow bool) error {
	order := &Order{}
	order.ID = id
	order.AllowComment = util.Tenary[uint](allow, 1, 2)
	if err := tx.Model(order).Updates(order).Error; err != nil {
		logger.Logger.Debugf("TxChangeOrderAllowCommentErr: %v\n", err)
		return err
	}
	return nil
}

func dbAppraiseOrder(id, appraisal, operator uint) (err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if err := txAppraiseOrder(tx, id, appraisal, operator); err != nil {
			logger.Logger.Debugf("AppraiseOrderErr: %v\n", err)
		}
		return err
	})
	return
}

func txAppraiseOrder(tx *gorm.DB, id, appraisal, operator uint) (err error) {
	order := &Order{}
	order.ID = id
	order.Appraisal = appraisal
	order.UpdatedBy = order.UserID
	order.AllowComment = CommentDisallow

	if err = tx.Model(order).Updates(order).Error; err != nil {
		return
	}
	status := NewStatusAppraised(operator)
	if err = txChangeOrderStatus(tx, id, status); err != nil {
		return
	}
	return
}
