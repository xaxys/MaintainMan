package order

import (
	"time"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func dbGetOrderByRepairer(id uint, json *RepairerOrderRequest) (orders []*Order, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if orders, count, err = txGetOrderByRepairer(tx, id, json); err != nil {
			mctx.Logger.Warnf("GetOrderByRepairer: %v\n", err)
		}
		return err
	})
	return
}

func txGetOrderByRepairer(tx *gorm.DB, id uint, json *RepairerOrderRequest) (orders []*Order, count uint, err error) {
	status := &Status{
		RepairerID: &id,
		Current:    json.Current,
	}
	statuses := []*Status{}
	tx = dao.TxPageFilter(tx, &json.PageParam).Model(status).Where(status)
	if json.Status != StatusIllegal {
		tx = tx.Joins("INNER JOIN orders ON orders.id = statuses.order_id AND orders.status = (?)", json.Status)
	}
	if len(json.Tags) > 0 {
		if json.Disjunctive {
			tx = tx.Where("id IN (?)", mctx.Database.Table("order_tags").Select("order_id").Where("tag_id IN (?)", json.Tags))
		} else {
			for _, tag := range json.Tags {
				tx = tx.Where("EXISTS (?)", mctx.Database.Table("order_tags").Select("order_id").Where("tag_id = (?)", tag).Where("order_id = statuses.order_id"))
			}
		}
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil || cnt == 0 {
		return
	}
	count = uint(cnt)
	if err = tx.Preload("Order.Tags").Find(&statuses).Error; err != nil {
		return
	}
	orders = util.TransSlice(statuses, func(status *Status) *Order { return status.Order })
	return
}

func dbGetStatusByOrder(id uint) (statuses []*Status, err error) {
	return txGetStatusByOrder(mctx.Database, id)
}

func txGetStatusByOrder(tx *gorm.DB, id uint) (statuses []*Status, err error) {
	status := &Status{
		OrderID: id,
	}
	if err = tx.Preload("Repairer").Where(status).Find(&statuses).Order("sequence_num").Error; err != nil {
		mctx.Logger.Warnf("GetStatusByOrderErr: %v\n", err)
		return
	}
	return
}

func txGetAppraiseTimeoutOrder(tx *gorm.DB) (ids []uint, err error) {
	status := &Status{
		Status:  StatusCompleted,
		Current: true,
	}
	statuses := []*Status{}

	timeout := orderConfig.GetDuration("appraise.timeout")
	exp := time.Now().Add(-timeout)

	if err = tx.Where(status).Where("created_at <= (?)", exp).Find(&statuses).Error; err != nil {
		mctx.Logger.Warnf("GetAppraiseTimeoutOrderErr: %v\n", err)
		return
	}

	ids = util.TransSlice(statuses, func(t *Status) uint { return t.OrderID })
	return
}

// NewStatus 创建状态
func NewStatus(status, repairer uint, operator uint) *Status {
	return &Status{
		Status:     status,
		RepairerID: util.Tenary(repairer != 0, &repairer, nil),
		Current:    true,
		BaseModel: model.BaseModel{
			CreatedBy: operator,
		},
	}
}

// StatusWaiting 待处理
func NewStatusWaiting(operator uint) *Status {
	return NewStatus(1, 0, operator)
}

// StatusAccepted 已接单
func NewStatusAssigned(repairer, operator uint) *Status {
	return NewStatus(2, repairer, operator)
}

// StatusCompleted 已完成
func NewStatusCompleted(operator uint) *Status {
	return NewStatus(3, 0, operator)
}

// StatusReported 上报中
func NewStatusReported(operator uint) *Status {
	return NewStatus(4, 0, operator)
}

// StatusHold 挂单
func NewStatusHold(operator uint) *Status {
	return NewStatus(5, 0, operator)
}

// StatusCanceled 已取消
func NewStatusCanceled(operator uint) *Status {
	return NewStatus(6, 0, operator)
}

// StatusRejected 已拒绝
func NewStatusRejected(operator uint) *Status {
	return NewStatus(7, 0, operator)
}

// StatusAppraised 已评价
func NewStatusAppraised(operator uint) *Status {
	return NewStatus(8, 0, operator)
}
