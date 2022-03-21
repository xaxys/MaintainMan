package model

type Tag struct {
	BaseModel
	Sort   string   `gorm:"not null; size:191; unique; index:idx_tag_sort_name,priority:1; comment:分类"`
	Name   string   `gorm:"not null; size:191; unique; index:idx_tag_sort_name,priority:2; comment:标签名称"`
	Level  uint     `gorm:"not null; size:5; default:0; comment:标签等级"`
	Orders []*Order `gorm:"many2many:order_tags;"`
}

type TagJson struct {
	Sort  string `json:"sort"`
	Name  string `json:"name"`
	Level uint   `json:"level"`
}
