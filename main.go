package main

import (
	"github.com/kataras/iris/v12"

	"maintainman/config"
	"maintainman/initialize"
	"maintainman/route"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel(config.AppConfig.GetString("app.loglevel"))
	initialize.InitDefaultData()
	route.Route(app)
	app.Listen(config.AppConfig.GetString("app.listen"))
}
