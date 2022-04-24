package wordcloud

import (
	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/logger"
	"gorm.io/gorm"
)

//TODO:Lacking tests.
func dbUploadWord(word *WordJson, id uint) (*WordJson, error) {
	return txUploadWord(mctx.Database, word, id)
}

func txUploadWord(tx *gorm.DB, word *WordJson, id uint) (*WordJson, error) {
	wordModel := &Word{}

	if err := tx.Model(wordModel).Where("order_id = ?", id).First(wordModel).Error; err == nil {
		if err := tx.Model(wordModel).Where("content = ?", word.Content).First(wordModel).Error; err != nil {
			if err := tx.Model(&Word{}).Create(&Word {
				Content: word.Content,
				WordClass: word.WordClass,
				Count: word.Count,
			}).Error; err != nil {
				logger.Logger.Fatalf("Create word Error: %v", err)
				return nil, err
			}
			
			return &WordJson{
				Content: word.Content,
				WordClass: word.WordClass,
				Count: word.Count,
			}, nil
		}
	    
		if err := tx.Model(wordModel).Update("count", wordModel.Count + word.Count).Error; err != nil {
			logger.Logger.Fatalf("Update word Error: %v", err)
			return nil, err
	    }

		return &WordJson{
			Content: word.Content,
			WordClass: word.WordClass,
			Count: word.Count,
		}, nil

	} else {
		logger.Logger.Fatalf("Find Order Error: %v", err)
		return nil, err
	}
    
}

func dbGetAllWords(aul *GetAllWordsRequest) (wordJson []*WordJson, count uint, err error) {
	return txGetAllWords(mctx.Database, aul)
}

func txGetAllWords(tx *gorm.DB, aul *GetAllWordsRequest) (wordJson []*WordJson, count uint, err error) {
	wordModels := make([]Word, 0)
	wordModel := &Word{}

	tx = dao.TxPageFilter(tx, &aul.PageParam).Where(wordModel)
	if err = tx.Model(wordModel).Find(&wordModels).Error; err != nil {
		return
	}
    
	cnt := int64(0)
	if err = tx.Model(wordModel).Count(&cnt).Error; err != nil {
		return
	}
    
	for _, word:= range wordModels {
		wordJson = append(wordJson, &WordJson {
			Content: word.Content,
			WordClass: word.WordClass,
            Count: word.Count,
		})
	}
	count = uint(cnt)
	return 
}

func dbGetWordsByOrderId(aul *GetWordsByOrderIdRequest) (wordJson []*WordJson, count uint, err error) {
	return txGetWordsByOrderId(mctx.Database, aul)
}

func txGetWordsByOrderId(tx *gorm.DB, aul *GetWordsByOrderIdRequest) (wordJson []*WordJson, count uint, err error) {
	wordModel := &Word{}
	wordModels := make([]Word, 0)

	if err := tx.Model(wordModel).Where("order_id = ?", aul.OrderId).First(wordModel).Error; err != nil {
		logger.Logger.Fatalf("Find Order Error: %v", err)
		return nil, 0, err
	} 
   
	tx = dao.TxPageFilter(tx, &aul.PageParam).Where(wordModel)
	if err = tx.Model(wordModel).Where("order_id = ?", aul.OrderId).Find(&wordModels).Error; err != nil {
		return
	}
    
	cnt := int64(0)
	if err = tx.Model(wordModel).Where("order_id = ?", aul.OrderId).Count(&cnt).Error; err != nil {
		return
	}
    
	for _, word:= range wordModels {
		wordJson = append(wordJson, &WordJson {
			Content: word.Content,
			WordClass: word.WordClass,
            Count: word.Count,
		})
	}
	count = uint(cnt)
	return 
}
