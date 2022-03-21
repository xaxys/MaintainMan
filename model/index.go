package model

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	CreatedBy uint
	UpdatedBy uint
}
