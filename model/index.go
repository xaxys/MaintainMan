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
	OrderBy string `json:"order_by" url:"order_by" validate:"omitempty,order_by"`
	Offset  uint   `json:"offset" url:"offset"`
	Limit   uint   `json:"limit" url:"limit"`
}

type AuthInfo struct {
	User uint
	Role string
	IP   string
}
