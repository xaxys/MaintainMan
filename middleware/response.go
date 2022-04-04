package middleware

import (
	"maintainman/model"

	"github.com/kataras/iris/v12"
)

var (
	ResponseHandler iris.Handler
)

func init() {
	ResponseHandler = func(ctx iris.Context) {
		response := ctx.Values().Get("response").(*model.ApiJson)
		if response != nil {
			// TODO: Temporary fix for inconsistent status code of log and response
			ctx.StatusCode(response.Code)
			ctx.JSON(response)
			ctx.StatusCode(response.Code)
		}
		ctx.Next()
	}
}
