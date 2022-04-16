package middleware

import (
	"github.com/xaxys/maintainman/core/model"

	"github.com/kataras/iris/v12"
)

var (
	ResponseHandler iris.Handler
)

func init() {
	ResponseHandler = func(ctx iris.Context) {
		response := ctx.Values().Get("response")
		if response == nil {
			ctx.Next()
			return
		}
		apiJson, ok := response.(*model.ApiJson)
		if !ok {
			ctx.Next()
			return
		}
		if apiJson != nil {
			// TODO: Temporary fix for inconsistent status code of log and response
			ctx.StatusCode(apiJson.Code)
			ctx.JSON(apiJson)
			ctx.StatusCode(apiJson.Code)
		}
		ctx.Next()
	}
}
