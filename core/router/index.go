package router

import (
	"github.com/xaxys/maintainman/core/middleware"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

var APIRoute iris.Party

func Register(app *iris.Application) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORS)
	app.AllowMethods(iris.MethodOptions)

	app.PartyFunc("/", func(home iris.Party) {
		home.HandleDir("/", "./assets")
		home.Get("/", func(ctx iris.Context) {
			ctx.Redirect("/index.html")
		})
	})

	v1 := app.Party("/v1")
	v1.Use(middleware.HeaderExtractor, middleware.TokenValidator)
	v1.Done(middleware.ResponseHandler)
	v1.SetExecutionRules(iris.ExecutionRules{Done: iris.ExecutionOptions{Force: true}})
	APIRoute = v1
}
