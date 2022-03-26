package controller

import (
	"maintainman/model"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func ExtractPageParam(ctx iris.Context) *model.PageParam {
	return &model.PageParam{
		OrderBy: ctx.URLParam("order_by"),
		Offset:  util.ToUint(ctx.URLParamIntDefault("offset", 0)),
		Limit:   util.ToUint(ctx.URLParamIntDefault("limit", 0)),
	}
}
