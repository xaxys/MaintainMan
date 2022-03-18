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
	CreatedBy string       `gorm:"size:191"`
	UpdatedBy string       `gorm:"size:191"`
	DeletedBy string       `gorm:"size:191"`
}
