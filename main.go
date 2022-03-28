package main

import (
	"github.com/kataras/iris/v12"

	"maintainman/config"
	"maintainman/initialize"
	"maintainman/logger"
	"maintainman/route"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel(config.AppConfig.GetString("app.loglevel"))
	logger.Logger = app.Logger()
	initialize.InitDefaultData()
	route.Route(app)
	err := app.Listen(config.AppConfig.GetString("app.listen"))
	if err != nil {
		return
	}
}
