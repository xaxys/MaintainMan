package dao

import (
	"database/sql"
	"maintainman/config"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/util"
	"time"

	"gorm.io/gorm"
)

func GetOrderByRepairer(id uint, json *model.RepairerOrderRequest) (orders []*model.Order, err error) {
	return TxGetOrderByRepairer(database.DB, id, json)
}

func TxGetOrderByRepairer(tx *gorm.DB, id uint, json *model.RepairerOrderRequest) (orders []*model.Order, err error) {
	status := &model.Status{
		RepairerID: sql.NullInt64{Int64: int64(id), Valid: true},
		Current:    json.Current,
	}
	statuses := []*model.Status{}
	tx = TxPageFilter(tx, &json.PageParam).Preload("Order.Tags").Where(status)
	if json.Status != 0 {
		tx = tx.Joins("Order", &model.Order{Status: json.Status})
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
	orders = util.TransSlice(statuses, func(status *model.Status) *model.Order { return status.Order })
	return
}

func TxGetAppraiseTimeoutOrder(tx *gorm.DB) (ids []uint, err error) {
	status := &model.Status{
		Status:  model.StatusCompleted,
		Current: true,
	}
	statuses := []*model.Status{}

	timeout := config.AppConfig.GetDuration("app.appraise.timeout")
	exp := time.Now().Add(-timeout)

	if err = tx.Where(status).Where("created_at <= (?)", exp).Find(&statuses).Error; err != nil {
		logger.Logger.Debugf("GetAppraiseTimeoutOrderErr: %v\n", err)
		return
	}

	ids = util.TransSlice(statuses, func(t *model.Status) uint { return t.OrderID })
	return
}

func NewStatus(status, repairer uint, operator uint) *model.Status {
	return &model.Status{
		Status:     status,
		RepairerID: sql.NullInt64{Int64: int64(repairer), Valid: repairer != 0},
		Current:    true,
		BaseModel: model.BaseModel{
			CreatedBy: operator,
		},
	}
}

// StatusWaiting 待处理
func StatusWaiting(operator uint) *model.Status {
	return NewStatus(1, 0, operator)
}

// StatusAccepted 已接单
func StatusAssigned(repairer, operator uint) *model.Status {
	return NewStatus(2, repairer, operator)
}

// StatusCompleted 已完成
func StatusCompleted(operator uint) *model.Status {
	return NewStatus(3, 0, operator)
}

// StatusReported 上报中
func StatusReported(operator uint) *model.Status {
	return NewStatus(4, 0, operator)
}

// StatusHold 挂单
func StatusHold(operator uint) *model.Status {
	return NewStatus(5, 0, operator)
}

// StatusCanceled 已取消
func StatusCanceled(operator uint) *model.Status {
	return NewStatus(6, 0, operator)
}

// StatusRejected 已拒绝
func StatusRejected(operator uint) *model.Status {
	return NewStatus(7, 0, operator)
}

// StatusAppraised 已评价
func StatusAppraised(operator uint) *model.Status {
	return NewStatus(8, 0, operator)
}
