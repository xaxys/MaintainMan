package dao

import (
	"maintainman/config"
	"maintainman/database"
	"maintainman/model"

	"gorm.io/gorm"
)

func PageFilter(param *model.PageParam) *gorm.DB {
	return Filter(param.OrderBy, param.Offset, param.Limit)
}

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
	PageLimit := config.AppConfig.GetInt("app.page.limit")
	if limit > PageLimit {
		limit = PageLimit
	}
	if limit <= 0 {
		limit = config.AppConfig.GetInt("app.page.default")
	}
	db = db.Limit(limit)
	return
}
