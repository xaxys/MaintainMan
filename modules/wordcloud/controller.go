package wordcloud

import (
	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/model"
)

// getAllWords godoc
// @Summary      查询所有词云
// @Description  查询所有词云
// @Tags         word
// @Accept       json
// @Produce      json
// @Param        order_by  query     string                                              false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset    query     uint                                                false  "偏移量 (默认为0)"
// @Param        limit     query     uint                                                false  "每页数据量 (默认为50)"
// @Success      200       {object}  model.ApiJson{data=model.Page{entries=[]WordJson}}  "返回结果"
// @Failure      400       {object}  model.ApiJson{data=[]string}
// @Failure      401       {object}  model.ApiJson{data=[]string}
// @Failure      500       {object}  model.ApiJson{data=[]string}
// @Router       /v1/word [get]
func getAllWords(ctx iris.Context) {
	param := &model.PageParam{}
	if err := ctx.ReadQuery(param); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := getAllWordsService(param)
	ctx.Values().Set("response", response)
}

// getWordsByOrder godoc
// @Summary      根据订单ID查询词云
// @Description  根据订单ID查询词云
// @Tags         word
// @Accept       json
// @Produce      json
// @Param        id        path      uint                                                true   "订单ID"
// @Param        order_by  query     string                                              false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset    query     uint                                                false  "偏移量 (默认为0)"
// @Param        limit     query     uint                                                false  "每页数据量 (默认为50)"
// @Success      200       {object}  model.ApiJson{data=model.Page{entries=[]WordJson}}  "返回结果"
// @Failure      400       {object}  model.ApiJson{data=[]string}
// @Failure      401       {object}  model.ApiJson{data=[]string}
// @Failure      500       {object}  model.ApiJson{data=[]string}
// @Router       /v1/word/{id} [get]
func getWordsByOrder(ctx iris.Context) {
	param := &model.PageParam{}
	if err := ctx.ReadQuery(param); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id, _ := ctx.Params().GetUint("id")
	response := getWordsByOrderService(id, param)
	ctx.Values().Set("response", response)
}
