package model

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	CreatedBy uint
	UpdatedBy uint
}

type PageParam struct {
	OrderBy string `url:"order_by" validate:"omitempty,order_by"`
	Offset  uint   `url:"offset"`
	Limit   uint   `url:"limit"`
}

type AuthInfo struct {
	User uint
	Role string
	IP   string
}
