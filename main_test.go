package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"maintainman/config"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/route"
	"maintainman/service"
	"maintainman/util"
	"math/rand"
	"strconv"
	"testing"
)

func TestRegisterRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("man", 10)

	e.GET("/").Expect().Status(httptest.StatusOK)

	for _, user := range users {
		e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
	}
}

func TestLoginRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("man", 10)

	for _, user := range users {
		e.POST("/v1/login").WithJSON(model.LoginJson{
			Account:  user.Name,
			Password: user.Password,
		}).Expect().Status(httptest.StatusOK)
	}
}

func TestUserReNewRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	e.GET("/v1/renew").WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK)
}

func getSuperAdminToken() string {
	superAdmin := initUser("超管", "12345678", "super_admin", "super_admin")
	apiJson := service.UserLogin(&model.LoginJson{
		Account:  superAdmin.Name,
		Password: superAdmin.Password,
	})
	return util.Strval(apiJson.Data)
}

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel(config.AppConfig.GetString("app.loglevel"))
	logger.Logger = app.Logger()
	route.Route(app)
	return app
}

func generateRandomUsers(prefix string, num uint) (usersRegister []model.ModifyUserJson) {
	for i := uint(1); i <= num; i++ {
		usersRegister = append(usersRegister, initUser(prefix+strconv.Itoa(int(i)), "12345678", "user"+strconv.Itoa(int(i)), "user"))
	}
	return
}

func initUser(name string, password string, displayName string, roleName string) model.ModifyUserJson {
	return model.ModifyUserJson{
		Name:        name,
		Password:    password,
		DisplayName: displayName,
		RoleName:    roleName,
		Phone:       strconv.Itoa(rand.Intn(100000)),
		Email:       strconv.Itoa(rand.Intn(100000)) + "@qq.com",
	}
}
