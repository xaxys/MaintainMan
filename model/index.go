package model

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	CreatedBy uint
	UpdatedBy uint
	DeletedBy uint
}
