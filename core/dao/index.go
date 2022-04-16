package dao

import (
	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/model"

	"gorm.io/gorm"
)

func TxPageFilter(tx *gorm.DB, param *model.PageParam) (db *gorm.DB) {
	return TxFilter(tx, param.OrderBy, param.Offset, param.Limit)
}

func TxFilter(tx *gorm.DB, orderBy string, uoffset, ulimit uint) (db *gorm.DB) {
	offset := int(uoffset)
	limit := int(ulimit)
	db = tx
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
