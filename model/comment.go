package model

type Comment struct {
	BaseModel
	OrderID     uint   `gorm:"not null; index:idx_comment_order_seqnum,priority:1; comment:订单ID"`
	UserID      uint   `gorm:"not null; comment:用户ID"`
	UserName    string `gorm:"not null; comment:用户名"`
	User        *User  `gorm:"foreignkey:UserID"`
	SequenceNum uint   `gorm:"not null; index:idx_comment_order_seqnum,priority:2; default:0; comment:发言序号"`
	Content     string `gorm:"not null; comment:内容"`
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required"`
}

type CommentJson struct {
	ID          uint   `json:"id"`
	OrderID     uint   `json:"order_id"`
	UserID      uint   `json:"user_id"`
	UserName    string `json:"user_name"`
	SequenceNum uint   `json:"sequence_num"` // 发言在该订单内的序号
	Content     string `json:"content"`
	CreatedAt   int64  `json:"created_at"` // unix timestamp in seconds (UTC)
}
