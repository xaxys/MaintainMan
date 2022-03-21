package model

type Order struct {
	BaseModel
	UserID       uint       `gorm:"not null; index:idx_order_user_status,priority:1; comment:用户ID"`
	User         *User      `gorm:"foreignkey:UserID"`
	Title        string     `gorm:"not null; unique; size:191; comment:标题"`
	Content      string     `gorm:"not null; comment:内容"`
	Address      string     `gorm:"not null; comment:地址"`
	ContactName  string     `gorm:"not null; size:191; comment:联系人"`
	ContactPhone string     `gorm:"not null; size:191; comment:联系电话"`
	Status       uint       `gorm:"not null; size:5; default:0; index:idx_order_user_status,priority:2; comment:状态 0:非法 1:待维修 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝"`
	StatusList   []*Status  `gorm:"foreignkey:OrderID"`
	AllowComment uint       `gorm:"not null; size:2 default:1; comment:是否允许评论 1:允许 2:不允许"`
	Comments     []*Comment `gorm:"foreignkey:OrderID"`
	ItemLogs     []*ItemLog `gorm:"foreignkey:OrderID"`
	Tags         []*Tag     `gorm:"many2many:order_tags;"`
	// Use Tag to solve this
	// DefectID     uint       `gorm:"not null; comment:故障分类ID"`
	// Defect       *Defect    `gorm:"foreignkey:DefectID"`
	// PlaceID      uint       `gorm:"not null; comment:地点ID"`
	// Place        *Place     `gorm:"foreignkey:PlaceID"`
}

type ModifyOrderJson struct {
	Title        string `json:"title" validate:"lte=191"`
	Content      string `json:"content"`
	Address      string `json:"address"`
	ContactName  string `json:"contact_name" validate:"lte=191"`
	ContactPhone string `json:"contact_phone" validate:"lte=191"`
	AddTags      []uint `json:"add_tags"`
	DelTags      []uint `json:"del_tags"`
	OperatorID   uint   // Filled by system
}

type AllOrderJson struct {
	Title      string `json:"title" validate:"lte=191"`
	UserID     uint   `json:"user_id"`
	Status     uint   `json:"status"`
	Tags       []uint `json:"tags"`
	Conjunctve bool   `json:"conjunctve"`
	OrderBy    string `json:"order_by"`
	Offset     uint   `json:"offset"`
	Limit      uint   `json:"limit"`
}
