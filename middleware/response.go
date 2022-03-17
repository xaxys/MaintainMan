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
			ctx.StatusCode(response.Code)
			ctx.JSON(response)
		}
		ctx.Next()
	}
}
