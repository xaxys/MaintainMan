package model

type Comment struct {
	BaseModel
	OrderID     uint   `gorm:"not null; index:idx_comment_order_seqnum,priority:1; comment:订单ID"`
	UserID      uint   `gorm:"not null; comment:用户ID"`
	User        *User  `gorm:"foreignkey:UserID"`
	SequenceNum uint   `gorm:"not null; index:idx_comment_order_seqnum,priority:2; default:0; comment:发言序号"`
	Content     string `gorm:"not null; comment:内容"`
	Photo       string `gorm:"not null; comment:图片"`
}
