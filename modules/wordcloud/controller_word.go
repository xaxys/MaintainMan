package wordcloud

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/kataras/iris/v12"
)

func uploadWords(ctx iris.Context) {
	aul := UploadWordsRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return 
	}

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return 
	}
	aul.OrderId = id

	response := uploadWordsService(&aul)
	ctx.JSON(response)
}

func getAllWords(ctx iris.Context) {
	aul := &GetAllWordsRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return 
	}
	
	response := getAllWordsService(aul)

    ctx.JSON(response)
}

func getWordByOrderId(ctx iris.Context) {
	aul := &GetAllWordsRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return 
	}

	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return 
	}

	req := GetWordsByOrderIdRequest {
		PageParam: aul.PageParam,
	}
	req.OrderId = id

}