package wordcloud

import (
	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/model"
	"gorm.io/gorm"
)

func dbUploadWord(id uint, json *WordJson) (err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if err = txUploadWord(tx, id, json); err != nil {
			logger.Logger.Warnf("UploadWordErr: %+v", err)
		}
		return err
	})
	return
}

func txUploadWord(tx *gorm.DB, id uint, word *WordJson) error {
	_, err := txUploadOrderWord(tx, id, word)
	if err != nil {
		return err
	}
	_, err = txUploadGlobalWord(tx, word)
	if err != nil {
		return err
	}
	return nil
}

func dbUploadOrderWord(id uint, json *WordJson) (word *OrderWord, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if word, err = txUploadOrderWord(tx, id, json); err != nil {
			logger.Logger.Warnf("UploadOrderWordErr: %+v", err)
		}
		return err
	})
	return
}

func txUploadOrderWord(tx *gorm.DB, id uint, json *WordJson) (*OrderWord, error) {
	word := &OrderWord{
		OrderID:   id,
		Content:   json.Content,
		WordClass: json.WordClass,
	}

	if err := tx.Where(word).Attrs(word).FirstOrCreate(word).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(word).Update("count", gorm.Expr("count + ?", json.Count)).Error; err != nil {
		return nil, err
	}
	word.Count += json.Count
	return word, nil
}

func dbUploadGlobalWord(json *WordJson) (word *GlobalWord, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if word, err = txUploadGlobalWord(tx, json); err != nil {
			logger.Logger.Warnf("UploadGlobalWordErr: %+v", err)
		}
		return err
	})
	return
}

func txUploadGlobalWord(tx *gorm.DB, json *WordJson) (*GlobalWord, error) {
	word := &GlobalWord{
		Content:   json.Content,
		WordClass: json.WordClass,
	}
	if err := tx.Where(word).Attrs(word).FirstOrCreate(word).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(word).Update("count", gorm.Expr("count + ?", json.Count)).Error; err != nil {
		return nil, err
	}
	word.Count += json.Count
	return word, nil
}

func dbGetAllWords(aul *model.PageParam) (words []*GlobalWord, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if words, count, err = txGetAllWords(tx, aul); err != nil {
			logger.Logger.Warnf("GetAllWordsErr: %+v", err)
		}
		return err
	})
	return
}

func txGetAllWords(tx *gorm.DB, aul *model.PageParam) (words []*GlobalWord, count uint, err error) {
	tx = dao.TxPageFilter(tx, aul)
	if err = tx.Find(&words).Error; err != nil {
		return
	}

	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}

func dbGetOrderWords(id uint, aul *model.PageParam) (words []*OrderWord, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if words, count, err = txGetOrderWords(tx, id, aul); err != nil {
			logger.Logger.Warnf("GetOrderWordsErr: %+v", err)
		}
		return err
	})
	return
}

func txGetOrderWords(tx *gorm.DB, id uint, aul *model.PageParam) (words []*OrderWord, count uint, err error) {
	word := &OrderWord{OrderID: id}
	tx = dao.TxPageFilter(tx, aul).Where(word)
	if err = tx.Find(&words).Error; err != nil {
		return
	}

	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}
