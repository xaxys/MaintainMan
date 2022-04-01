package dao

import (
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
		RepairerID: id,
		Current:    json.Current,
	}
	statuses := []*model.Status{}
	if err = TxPageFilter(tx, &json.PageParam).Preload("Order").Where(status).Find(&statuses).Error; err != nil {
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
		RepairerID: repairer,
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
