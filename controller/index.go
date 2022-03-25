package controller

import (
	"maintainman/model"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func ExtractPageParam(ctx iris.Context) (p model.PageParam) {
	p.OrderBy = ctx.URLParam("order_by")
	p.Offset = util.ToUint(ctx.URLParamIntDefault("offset", 0))
	p.Limit = util.ToUint(ctx.URLParamIntDefault("limit", 0))
	return
}
