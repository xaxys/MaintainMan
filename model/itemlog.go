package model

type ItemLog struct {
	BaseModel
	ItemID    uint   `gorm:"not null; comment:物品ID"`
	Item      *Item  `gorm:"foreignkey:ItemID;"`
	OrderID   uint   `gorm:"not null; comment:订单ID"`
	Order     *Order `gorm:"foreignkey:OrderID;"`
	ChangeNum int    `gorm:"not null; default:0; comment:增加/消耗数量 正:增加 负:减少"`
	ItemCount uint   `gorm:"not null; default:0; comment:增加/消耗后物品数量"`
}
