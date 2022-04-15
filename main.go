package main

import (
	"fmt"

	"github.com/kataras/iris/v12"

	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/initialize"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/router"
	"github.com/xaxys/maintainman/core/service"
	"github.com/xaxys/maintainman/core/util"
	"github.com/xaxys/maintainman/module"
	"github.com/xaxys/maintainman/modules/announce"
	"github.com/xaxys/maintainman/modules/imagehost"
	"github.com/xaxys/maintainman/modules/order"
	"github.com/xaxys/maintainman/modules/sysinfo"
)

var (
	BuildTags = "unknown"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GoVersion = "unknown"
)

// Font: smslant
// http://www.network-science.de/ascii/
func printBanner() {
	fmt.Println()
	fmt.Println("   __  ___       _       __         _        __  ___          ")
	fmt.Println("  /  |/  /___ _ (_)___  / /_ ___ _ (_)___   /  |/  /___ _ ___ ")
	fmt.Println(" / /|_/ // _ `// // _ \\/ __// _ `// // _ \\ / /|_/ // _ `// _ \\")
	fmt.Println("/_/  /_/ \\_,_//_//_//_/\\__/ \\_,_//_//_//_//_/  /_/ \\_,_//_//_/")
	fmt.Println()
	fmt.Println("Welcome to use MaintainMan!")
	fmt.Println("Version:   " + BuildTags)
	fmt.Println("Built:     " + BuildTime)
	fmt.Println("GitCommit: " + GitCommit)
	fmt.Println("GoVersion: " + GoVersion)
	fmt.Println()
}

func main() {
	printBanner()
	app := newApp()
	app.Logger().SetLevel(config.AppConfig.GetString("app.loglevel"))
	app.Listen(config.AppConfig.GetString("app.listen"))
}

func newApp() *iris.Application {
	app := iris.New()
	logger.Logger = app.Logger()
	initialize.InitDefaultData()
	router.Register(app)
	server := module.Server{
		Validator: util.Validator,
		Logger:    app.Logger(),
		Scheduler: service.Scheduler,
		Database:  database.DB,
	}
	registry := module.NewRegistry(&server)
	registry.Register(
		&imagehost.Module,
		&announce.Module,
		&order.Module,
		&sysinfo.Module,
	)
	service.Scheduler.StartAsync()
	return app
}
