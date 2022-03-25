package controller

import (
	"maintainman/model"
	"maintainman/service"

	"github.com/kataras/iris/v12"
)

func GetTagByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.GetTagByID(id)
	ctx.Values().Set("response", response)
}

func GetAllTagSorts(ctx iris.Context) {
	response := service.GetAllTagSorts()
	ctx.Values().Set("response", response)
}

func GetAllTagsBySort(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	response := service.GetAllTagsBySort(name)
	ctx.Values().Set("response", response)
}

func CreateTag(ctx iris.Context) {
	aul := &model.CreateTagJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.CreateTag(aul)
	ctx.Values().Set("response", response)
}

func DeleteTagByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	response := service.DeleteTag(id)
	ctx.Values().Set("response", response)
}
