package wordcloud

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"
)

func uploadWordsService(id uint, content string) *model.ApiJson {
	wc := NewWordCollectorWithStr(content)
	wordSet := wc.Filter(&LengthFilter{}).ToSlice()
	errs := []error{}
	for _, word := range wordSet {
		err := dbUploadWord(id, &word)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return model.ErrorInternalServer(errs...)
	}
	return model.Success(nil, "上传成功")
}

func getAllWordsService(aul *model.PageParam) *model.ApiJson {
	aul.OrderBy = util.NotEmpty(aul.OrderBy, "count desc")
	words, count, err := dbGetAllWords(aul)
	if err != nil {
		return model.ErrorInternalServer(err)
	}
	ws := util.TransSlice(words, globalWordToJson)
	return model.SuccessPaged(ws, count, "获取成功")
}

func getWordsByOrderService(id uint, aul *model.PageParam) *model.ApiJson {
	aul.OrderBy = util.NotEmpty(aul.OrderBy, "count desc")
	words, count, err := dbGetOrderWords(id, aul)
	if err != nil {
		return model.ErrorInternalServer(err)
	}
	ws := util.TransSlice(words, orderWordToJson)
	return model.SuccessPaged(ws, count, "获取成功")
}

func globalWordToJson(w *GlobalWord) *WordJson {
	return &WordJson{
		Content:   w.Content,
		WordClass: w.WordClass,
		Count:     w.Count,
	}
}

func orderWordToJson(w *OrderWord) *WordJson {
	return &WordJson{
		Content:   w.Content,
		WordClass: w.WordClass,
		Count:     w.Count,
	}
}
