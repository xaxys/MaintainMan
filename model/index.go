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
	CreatedBy string       `gorm:"not null VARCHAR(191)"`
	UpdatedBy string       `gorm:"not null VARCHAR(191)"`
	DeletedBy string       `gorm:"not null VARCHAR(191)"`
}
