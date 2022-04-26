package order

import "github.com/xaxys/maintainman/core/model"

type ItemLog struct {
	model.BaseModel
	ItemID      uint    `gorm:"not null; comment:物品ID"`
	Item        *Item   `gorm:"foreignkey:ItemID;"`
	OrderID     *uint   `gorm:"comment:订单ID"`
	Order       *Order  `gorm:"foreignkey:OrderID;"`
	ChangeNum   int     `gorm:"not null; default:0; comment:增加/消耗数量 正:增加 负:减少"`
	ChangePrice float64 `gorm:"not null; default:0; comment:开销 正:进货 负:订单收费"`
}

type AddItemRequest struct {
	ItemID uint    `json:"-"`
	Num    uint    `json:"num"`
	Price  float64 `json:"price"`
}

type ConsumeItemRequest struct {
	ItemID  uint    `json:"item_id"`
	OrderID uint    `json:"order_id"`
	Num     uint    `json:"num"`
	Price   float64 `json:"price"`
}

type ItemLogJson struct {
	ID          uint    `json:"id"`
	ItemID      uint    `json:"item_id"`
	OrderID     uint    `json:"order_id"`
	ChangeNum   int     `json:"change_num"`   // 增加/消耗数量 正:增加 负:减少
	ChangePrice float64 `json:"change_price"` // 开销 正:进货 负:订单收费
	CreatedAt   int64   `json:"created_at"`   // unix timestamp in seconds (UTC)
	CreatedBy   uint    `json:"created_by"`   // unix timestamp in seconds (UTC)
}
