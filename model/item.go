package model

type Item struct {
	BaseModel
	Name        string     `gorm:"not null; size:191; unique; comment:物品名称"`
	Description string     `gorm:"not null; comment:物品描述"`
	Price       float64    `gorm:"not null; default:0; comment:物品总价值"`
	Income      float64    `gorm:"not null; default:0; comment:维修收入"`
	Count       int        `gorm:"not null; default:0; comment:物品数量"`
	ItemLogs    []*ItemLog `gorm:"foreignkey:ItemID"`
}

type CreateItemRequest struct {
	Name        string `json:"name" validate:"required,lte=191"`
	Discription string `json:"discription"`
}

type ItemInfoJson struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"discription"`
	Price       float64        `json:"price"`
	Income      float64        `json:"income"`
	Count       int            `json:"count"`
	ItemLogs    []*ItemLogJson `json:"item_log"`
	CreatedAt   int64          `json:"created_at"` // unix timestamp in seconds (UTC)
	UpdatedAt   int64          `json:"updated_at"` // unix timestamp in seconds (UTC)
	CreatedBy   uint           `json:"created_by"`
	UpdatedBy   uint           `json:"updated_by"`
}

type ItemJson struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"discription"`
	Count       int    `json:"count"`
}
