package wordcloud

import (
	"github.com/xaxys/maintainman/core/model"
)

type LengthFilter struct {

}

func (lf *LengthFilter) IsLegal(word WordJson) bool {
	return len(word.Content) >= 4
}

func uploadWordsService(aul *UploadWordsRequest) *model.ApiJson {
	wc := NewWordCollectorWithStr(aul.Content)
	wordSet := wc.Filter(&LengthFilter{}).ToSlice()
	words := make([]WordJson, 0)
	for _, word := range wordSet {
		word, err := dbUploadWord(&word, aul.OrderId)
		if err != nil {
			return model.ErrorInternalServer(err)
		}

		words = append(words, *word)
	
	}
	return &model.ApiJson {
		Code: 200,
		Status: true,
		Msg: "Upload words Successfully.",
		Data: words,
	}
}

func getAllWordsService(aul *GetAllWordsRequest) *model.ApiJson {
	words, _, err := dbGetAllWords(aul)
	if err != nil {
		return model.ErrorInternalServer(err)
	}
    
    return &model.ApiJson {
		Code: 200,
		Status: true,
		Msg: "Get words Successfully.",
		Data: accumulateWords(words),
	}
}

func GetWordsByOrderIdService(aul *GetWordsByOrderIdRequest) *model.ApiJson {
	words, _, err := dbGetWordsByOrderId(aul)
    if err != nil {
		return model.ErrorInternalServer(err)
	}

	return &model.ApiJson {
		Code: 200,
		Status: true,
		Msg: "Get words Successfully.",
		Data: accumulateWords(words),
	}
}