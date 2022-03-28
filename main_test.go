package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/spf13/cast"
	"maintainman/config"
	"maintainman/logger"
	"maintainman/model"
	"maintainman/route"
	"maintainman/service"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestRegisterRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("testUser", 10)

	for _, user := range users {
		e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
	}
}

func TestLoginRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("testUser", 10)

	for _, user := range users {
		e.POST("/v1/login").WithJSON(model.LoginRequest{
			Account:  user.Name,
			Password: user.Password,
		}).Expect().Status(httptest.StatusOK)
	}
}

//FIXME: record not found [0.334ms] [rows:0] SELECT * FROM `users` WHERE `users`.`id` = 0 AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
func TestUserReNewRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	e.GET("/v1/renew").WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK)
}

func getSuperAdminToken() string {
	superAdmin := initUser("admin", "12345678", "maintainman default admin")
	apiJson := service.UserLogin(&model.LoginRequest{
		Account:  superAdmin.Name,
		Password: superAdmin.Password,
	}, getMyIPV6(), nil)
	return cast.ToString(apiJson.Data)
}

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel(config.AppConfig.GetString("app.loglevel"))
	logger.Logger = app.Logger()
	route.Route(app)
	return app
}

func generateRandomUsers(prefix string, num uint) (usersRegister []model.RegisterUserRequest) {
	for i := uint(1); i <= num; i++ {
		usersRegister = append(usersRegister, initUser(prefix+strconv.Itoa(int(i)), "12345678", "user"+strconv.Itoa(int(i))))
	}
	return
}

func initUser(name string, password string, displayName string) model.RegisterUserRequest {
	return model.RegisterUserRequest{
		Name:        name,
		Password:    password,
		DisplayName: displayName,
		Phone:       strconv.Itoa(rand.Intn(100000)),
		Email:       strconv.Itoa(rand.Intn(100000)) + "@qq.com",
	}
}

func getMyIPV6() string {
	s, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range s {
		i := regexp.MustCompile(`(\w+:){7}\w+`).FindString(a.String())
		if strings.Count(i, ":") == 7 {
			return i
		}
	}
	return ""
}
