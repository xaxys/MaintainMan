package model

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

// immutable
type Status struct {
	BaseModel
	OrderID     uint   `gorm:"not null; comment:订单ID"`
	Order       *Order `gorm:"foreignkey:OrderID;"`
	Status      uint   `gorm:"not null; size:5; default:0; comment:状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"`
	Current     bool   `gorm:"not null; index:idx_status_repairer_current,priority:2; size:1; default:0; comment:是否最新状态"`
	RepairerID  uint   `gorm:"not null; index:idx_status_repairer_current,priority:1; comment:维修员ID"`
	Repairer    *User  `gorm:"foreignkey:RepairerID;"`
	SequenceNum uint   `gorm:"not null; default:0; comment:状态序号"`
}

type StatusJson struct {
	Status      uint   `json:"status"`
	RepairerID  uint   `json:"repairer_id"`
	Repairer    *User  `json:"repairer"`
	SequenceNum uint   `json:"sequence_num"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}
