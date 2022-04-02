package model

type Tag struct {
	BaseModel
	Sort     string   `gorm:"not null; size:191; index:idx_tag_sort_name,priority:1; comment:分类"`
	Name     string   `gorm:"not null; size:191; index:idx_tag_sort_name,priority:2; comment:标签名称"`
	Level    uint     `gorm:"not null; size:5; default:0; comment:标签等级"`
	Congener uint     `gorm:"not null; default:0; comment:同类型数量"`
	Orders   []*Order `gorm:"many2many:order_tags;"`
}

type CreateTagRequest struct {
	Sort     string `json:"sort" validate:"required,lte=191"`
	Name     string `json:"name" validate:"required,lte=191"`
	Level    uint   `json:"level" validate:"required,gte=0"`
	Congener uint   `json:"congener" validate:"omitempty,gte=0"` // 允许与同Sort的Tag共存的数量 0:不限 n:只允许n个(含自身)
}

type TagJson struct {
	ID       uint   `json:"id"`
	Sort     string `json:"sort"`
	Name     string `json:"name"`
	Level    uint   `json:"level"`
	Congener uint   `json:"congener"` // 允许与同Sort的Tag共存的数量 0:不限 n:只允许n个(含自身)
}
