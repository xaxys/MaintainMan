package order

import (
	"database/sql"
	"time"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func dbGetOrderByRepairer(id uint, json *RepairerOrderRequest) (orders []*Order, err error) {
	return txGetOrderByRepairer(mctx.Database, id, json)
}

func txGetOrderByRepairer(tx *gorm.DB, id uint, json *RepairerOrderRequest) (orders []*Order, err error) {
	status := Status{
		RepairerID: sql.NullInt64{Int64: int64(id), Valid: true},
		Current:    json.Current,
	}
	statuses := []*Status{}
	tx = dao.TxPageFilter(tx, &json.PageParam).Preload("Order.Tags").Where(status)
	if json.Status != 0 {
		tx = tx.Joins("Order", Order{Status: json.Status})
	}
	if len(json.Tags) > 0 {
		if json.Disjunctive {
			tx = tx.Where("id IN (?)", database.DB.Table("order_tags").Select("order_id").Where("tag_id IN (?)", json.Tags))
		} else {
			for _, tag := range json.Tags {
				tx = tx.Where("EXISTS (?)", database.DB.Table("order_tags").Select("order_id").Where("tag_id = (?)", tag).Where("order_id = statuses.order_id"))
			}
		}
	}
	if err = tx.Find(&statuses).Error; err != nil {
		logger.Logger.Debugf("GetOrderByRepairerErr: %v\n", err)
		return
	}
	orders = util.TransSlice(statuses, func(status *Status) *Order { return status.Order })
	return
}

func txGetAppraiseTimeoutOrder(tx *gorm.DB) (ids []uint, err error) {
	status := Status{
		Status:  StatusCompleted,
		Current: true,
	}
	statuses := []*Status{}

	timeout := orderConfig.GetDuration("appraise.timeout")
	exp := time.Now().Add(-timeout)

	if err = tx.Where(status).Where("created_at <= (?)", exp).Find(&statuses).Error; err != nil {
		logger.Logger.Debugf("GetAppraiseTimeoutOrderErr: %v\n", err)
		return
	}

	ids = util.TransSlice(statuses, func(t *Status) uint { return t.OrderID })
	return
}

// NewStatus 创建状态
func NewStatus(status, repairer uint, operator uint) *Status {
	return &Status{
		Status:     status,
		RepairerID: sql.NullInt64{Int64: int64(repairer), Valid: repairer != 0},
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
