package dao

import (
	"maintainman/config"
	"maintainman/database"

	"gorm.io/gorm"
)

func Filter(orderBy string, uoffset, ulimit uint) (db *gorm.DB) {
	offset := int(uoffset)
	limit := int(ulimit)
	db = database.DB
	if len(orderBy) > 0 {
		db = db.Order(orderBy)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	PageLimit := config.AppConfig.GetInt("app.page_limit")
	if limit > PageLimit {
		limit = PageLimit
	}
	if limit <= 0 {
		limit = config.AppConfig.GetInt("app.page_limit_default")
	}
	db = db.Limit(limit)
	return
}
