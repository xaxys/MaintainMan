package model

// immutable
type Status struct {
	BaseModel
	OrderID     uint  `gorm:"not null; comment:订单ID"`
	Status      uint  `gorm:"not null; size:5; default:0; comment:状态 0:非法 1:待维修 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝"`
	RepairerID  uint  `gorm:"not null; comment:维修员ID"`
	Repairer    *User `gorm:"foreignkey:RepairerID;"`
	SequenceNum uint  `gorm:"not null; default:0; comment:状态序号"`
}
