package order

import (
	"github.com/xaxys/maintainman/modules/user"

	"github.com/xaxys/maintainman/core/model"
)

const (
	StatusIllegal = iota
	StatusWaiting
	StatusAssigned
	StatusCompleted
	StatusReported
	StatusHold
	StatusCanceled
	StatusRejected
	StatusAppraised
)

func StatusName(status int) string {
	switch status {
	case StatusIllegal:
		return "非法状态"
	case StatusWaiting:
		return "待处理"
	case StatusAssigned:
		return "已接单"
	case StatusCompleted:
		return "已完成"
	case StatusReported:
		return "上报中"
	case StatusHold:
		return "挂单中"
	case StatusCanceled:
		return "已取消"
	case StatusRejected:
		return "已拒绝"
	case StatusAppraised:
		return "已评价"
	default:
		return "未知状态"
	}
}

// immutable
type Status struct {
	model.BaseModel
	OrderID     uint       `gorm:"not null; comment:订单ID"`
	Order       *Order     `gorm:"foreignkey:OrderID;"`
	Status      uint       `gorm:"not null; size:5; default:0; comment:状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"`
	Current     bool       `gorm:"not null; index:idx_status_repairer_current,priority:2; size:1; default:0; comment:是否最新状态"`
	RepairerID  *uint      `gorm:"index:idx_status_repairer_current,priority:1; comment:维修员ID"`
	Repairer    *user.User `gorm:"foreignkey:RepairerID;"`
	SequenceNum uint       `gorm:"not null; default:0; comment:状态序号"`
}

type StatusJson struct {
	SequenceNum uint           `json:"sequence_num"` // 状态序号
	Status      uint           `json:"status"`       // 状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价
	RepairerID  uint           `json:"repairer_id"`
	Repairer    *user.UserJson `json:"repairer"`
	CreatedAt   int64          `json:"created_at"` // unix timestamp in seconds (UTC)
	UpdatedAt   int64          `json:"updated_at"` // unix timestamp in seconds (UTC)
	CreatedBy   uint           `json:"created_by"`
	UpdatedBy   uint           `json:"updated_by"`
}
