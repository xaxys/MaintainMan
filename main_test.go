package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"maintainman/config"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/route"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.POST("/v1/register").WithJSON(model.ModifyUserJson{
		Name:        "test1",
		Password:    "12345678",
		DisplayName: "admin",
	}).Expect().Status(httptest.StatusOK)

	e.POST("/v1/login").WithJSON(model.LoginJson{
		Account:  "张城玮",
		Password: "12345678",
	}).Expect().Status(httptest.StatusOK)
}

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel(config.AppConfig.GetString("app.loglevel"))
	logger.Logger = app.Logger()
	route.Route(app)
	return app
}
