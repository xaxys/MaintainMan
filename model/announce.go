package model

import (
	"database/sql"
	"time"
)

type Announce struct {
	BaseModel
	Title     string       `gorm:"not null; unique; size:191; comment:标题"`
	Content   string       `gorm:"not null; comment:内容"`
	StartTime sql.NullTime `gorm:"index:idx_announce_time,priority:1; default:0000-00-00 00:00:00; comment:开始时间"`
	EndTime   sql.NullTime `gorm:"index:idx_announce_time,priority:2; default:0000-00-00 00:00:00; comment:结束时间"`
	Hits      uint         `gorm:"not null; default:0; comment:点击数"`
}

type ModifyAnnounceJson struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type AllAnnounceJson struct {
	Title   string `json:"title" validate:"lte=191"`
	OrderBy string `json:"order_by"`
	Limit   int    `json:"limit" validate:"number"`
	Offset  int    `json:"offset" validate:"number"`
}

type AnnounceJson struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	RoleName    string    `json:"user_role"`
	Role        *RoleJson `json:"role,omitempty"`
}
