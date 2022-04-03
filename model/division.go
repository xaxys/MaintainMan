package model

import "gorm.io/gorm"

type Division struct {
	gorm.Model
	Name     string      `gorm:"not null; size:191; unique; comment:分组名称"`
	ParentID uint        `gorm:"not null; default:0; comment:父分组ID"`
	Children []*Division `gorm:"foreignkey:ParentID"`
}

type CreateDivisionRequest struct {
	Name     string `json:"name" validate:"required,lte=191"`
	ParentID uint   `json:"parent_id"`
}

type UpdateDivisionRequest struct {
	Name     string `json:"name" validate:"omitempty,lte=191"`
	ParentID uint   `json:"parent_id"`
}

type DivisionJson struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	ParentID uint            `json:"parent_id"` // 父分组ID
	Children []*DivisionJson `json:"children"`
}
