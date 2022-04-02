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
	OrderBy string `json:"order_by" url:"order_by" validate:"omitempty,order_by"` // 排序字段 (默认为ID正序) 只接受`{field} {asc|desc}`格式 (e.g. `id desc`)
	Offset  uint   `json:"offset" url:"offset"`                                   // 偏移量 (默认为0)
	Limit   uint   `json:"limit" url:"limit"`                                     // 每页记录数 (默认为50)
}

type AuthInfo struct {
	User uint
	Role string
	IP   string
}
