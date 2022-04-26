package order

import (
	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/util"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func dbGetOrderCount() (count uint, err error) {
	return txGetOrderCount(mctx.Database)
}

func txGetOrderCount(tx *gorm.DB) (uint, error) {
	count := int64(0)
	if err := tx.Model(&Order{}).Count(&count).Error; err != nil {
		mctx.Logger.Warnf("GetOrderCountErr: %v\n", err)
		return 0, err
	}
	return uint(count), nil
}

// GetSimpleOrderByID return no relative info
func dbGetSimpleOrderByID(id uint) (*Order, error) {
	return txGetSimpleOrderByID(mctx.Database, id)
}

// txGetSimpleOrderByID return no relative info
func txGetSimpleOrderByID(tx *gorm.DB, id uint) (*Order, error) {
	order := &Order{}
	if err := tx.First(order, id).Error; err != nil {
		mctx.Logger.Warnf("TxGetOrderByIDErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func dbGetOrderByID(id uint) (*Order, error) {
	return txGetOrderByID(mctx.Database, id)
}

func txGetOrderByID(tx *gorm.DB, id uint) (*Order, error) {
	order := &Order{}
	if err := tx.Preload("Tags").Preload("Comments").Preload("StatusList", "current = true").First(order, id).Error; err != nil {
		mctx.Logger.Warnf("TxGetOrderByIDErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func dbGetAllOrdersWithParam(aul *AllOrderRequest) (orders []*Order, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if orders, count, err = txGetAllOrdersWithParam(tx, aul); err != nil {
			mctx.Logger.Warnf("GetAllOrdersWithParam: %v\n", err)
		}
		return err
	})
	return
}

func txGetAllOrdersWithParam(tx *gorm.DB, aul *AllOrderRequest) (orders []*Order, count uint, err error) {
	order := &Order{
		UserID: aul.UserID,
		Status: aul.Status,
	}
	tx = dao.TxPageFilter(tx, &aul.PageParam).Model(order).Where(order)
	if len(aul.Tags) > 0 {
		if aul.Disjunctive {
			tx = tx.Where("id IN (?)", mctx.Database.Table("order_tags").Select("order_id").Where("tag_id IN (?)", aul.Tags))
		} else {
			for _, tag := range aul.Tags {
				tx = tx.Where("EXISTS (?)", mctx.Database.Table("order_tags").Select("order_id").Where("tag_id = (?)", tag).Where("order_id = orders.id"))
			}
		}
	}
	if aul.Title != "" {
		tx = tx.Where("title LIKE ?", aul.Title)
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil || cnt == 0 {
		return
	}
	count = uint(cnt)
	if err = tx.Preload("Tags").Find(&orders).Error; err != nil {
		return
	}
	return
}

func dbGetOrderWithLastStatus(id uint) (*Order, error) {
	return txGetOrderWithLastStatus(mctx.Database, id)
}

func txGetOrderWithLastStatus(tx *gorm.DB, id uint) (*Order, error) {
	order := &Order{}
	if err := tx.Preload("StatusList", "current = true").Model(order).Find(order, id).Error; err != nil {
		mctx.Logger.Warnf("GetOrderWithLastStatusErr: %v\n", err)
		return nil, err
	}
	return order, nil
}

func dbCreateOrder(aul *CreateOrderRequest, operator uint) (order *Order, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if order, err = txCreateOrder(tx, aul, operator); err != nil {
			mctx.Logger.Warnf("CreateOrderErr: %v\n", err)
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
			mctx.Logger.Warnf("UpdateOrderErr: %v\n", err)
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
		mctx.Logger.Warnf("TxDeleteOrderErr: %v\n", err)
		return err
	}
	return nil
}

func dbChangeOrderStatus(id uint, status *Status) (err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if err = txChangeOrderStatus(tx, id, status); err != nil {
			mctx.Logger.Warnf("ChangeOrderStatusErr: %v\n", err)
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
	lastStatus := util.LastElem(or.StatusList)

	updateStatus := map[string]any{
		"current":    false,
		"updated_by": status.CreatedBy,
	}
	if err := tx.Model(lastStatus).Select("current", "updated_by").Updates(updateStatus).Error; err != nil {
		return err
	}

	status.SequenceNum = lastStatus.SequenceNum + 1
	if err := tx.Model(order).Association("StatusList").Append(status); err != nil {
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
		mctx.Logger.Warnf("TxChangeOrderAllowCommentErr: %v\n", err)
		return err
	}
	return nil
}

func dbAppraiseOrder(id, appraisal, operator uint) (err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if err := txAppraiseOrder(tx, id, appraisal, operator); err != nil {
			mctx.Logger.Warnf("AppraiseOrderErr: %v\n", err)
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
