package model

import (
	"time"
)

type Announce struct {
	BaseModel
	Title     string     `gorm:"not null; unique; size:191; comment:标题"`
	Content   string     `gorm:"not null; comment:内容"`
	StartTime *time.Time `gorm:"not null; size:8; index:idx_announce_start_end,priority:1; index:idx_announce_end_start,priority:2; comment:开始时间"`
	EndTime   *time.Time `gorm:"not null; size:8; index:idx_announce_start_end,priority:2; index:idx_announce_end_start,priority:1; comment:结束时间"`
	Hits      uint       `gorm:"not null; default:0; index; comment:点击数"`
}

type CreateAnnounceRequest struct {
	Title     string `json:"title" validate:"required,lte=191"`
	Content   string `json:"content" validate:"required"`
	StartTime int64  `json:"start_time" validate:"required,gte=-1,lte=253370764799"`                // unix timestamp in seconds (UTC); -1代表不限; 含本数
	EndTime   int64  `json:"end_time" validate:"required,eq=-1|gtfield=StartTime,lte=253370764799"` // unix timestamp in seconds (UTC); -1代表不限; 含本数; 必须大于start
}

type UpdateAnnounceRequest struct {
	Title     string `json:"title" validate:"lte=191"`
	Content   string `json:"content"`
	StartTime int64  `json:"start_time" validate:"required,gte=-1,lte=253370764799"`                // unix timestamp in seconds (UTC); -1代表不限; 含本数
	EndTime   int64  `json:"end_time" validate:"required,eq=-1|gtfield=StartTime,lte=253370764799"` // unix timestamp in seconds (UTC); -1代表不限; 含本数
}

type ModifyAnnounceRequest struct {
	Title     string
	Content   string
	StartTime int64
	EndTime   int64
}

type AllAnnounceRequest struct {
	Title     string `json:"title" url:"title" validate:"lte=191"`
	StartTime int64  `json:"start_time" url:"start_time" validate:"gte=-1,lte=253370764799"` // unix timestamp in seconds (UTC); -1代表不限; 含本数
	EndTime   int64  `json:"end_time" url:"end_time" validate:"gte=-1,lte=253370764799"`     // unix timestamp in seconds (UTC); -1代表不限; 含本数
	Inclusive bool   `json:"inclusive" url:"inclusive"`                                      // true: 查询开始时间晚于start,且结束时间早于end的(在某段时间内开始并结束的); false: 查询开始时间早于start,且结束时间晚于end的(在某段时间内都能看到的)
	PageParam
}

type AnnounceJson struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	StartTime int64  `json:"start_time"` // unix timestamp in seconds (UTC)
	EndTime   int64  `json:"end_time"`   // unix timestamp in seconds (UTC)
	Hits      uint   `json:"hits"`       // 点击数
	CreatedAt int64  `json:"created_at"` // unix timestamp in seconds (UTC)
	UpdatedAt int64  `json:"updated_at"` // unix timestamp in seconds (UTC)
}
