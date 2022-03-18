package model

type Item struct {
	BaseModel
	Name        string  `gorm:"not null; size:191; unique; comment:物品名称"`
	Discription string  `gorm:"not null; comment:物品描述"`
	Price       float64 `gorm:"not null; comment:物品价格"`
	Count       uint    `gorm:"not null; default:0; comment:物品数量"`
}
