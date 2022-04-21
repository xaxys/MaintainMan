package main

import (
	"fmt"

	"github.com/kataras/iris/v12"

	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/router"
	"github.com/xaxys/maintainman/core/service"
	"github.com/xaxys/maintainman/core/util"
	"github.com/xaxys/maintainman/modules/announce"
	"github.com/xaxys/maintainman/modules/imagehost"
	"github.com/xaxys/maintainman/modules/order"
	"github.com/xaxys/maintainman/modules/role"
	"github.com/xaxys/maintainman/modules/sysinfo"
	"github.com/xaxys/maintainman/modules/user"
	"github.com/xaxys/maintainman/modules/wxnotify"
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

// @title         MaintainMan API
// @version       1.0.0-rc3
// @license.name  MIT
func main() {
	printBanner()
	app := newApp()
	app.Listen(config.AppConfig.GetString("app.listen"))
}

var logLevel = config.AppConfig.GetString("app.loglevel")

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel(logLevel)
	logger.Logger = app.Logger()
	router.Register(app)
	server := module.Server{
		Validator: util.Validator,
		Logger:    app.Logger(),
		Scheduler: service.Scheduler,
		EventBus:  service.Bus,
		Database:  database.DB,
	}
	registry := module.NewRegistry(&server)
	registry.Register(
		&role.Module,
		&user.Module,
		&imagehost.Module,
		&announce.Module,
		&order.Module,
		&wxnotify.Module,
		&sysinfo.Module,
	)
	service.Scheduler.StartAsync()
	return app
}
