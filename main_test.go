package main

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/kataras/iris/v12/httptest"
	"github.com/spf13/cast"
)

func TestRegisterAndLoginRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("testUser", 10)

	for _, user := range users {
		e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
	}
	for _, user := range users {
		e.POST("/v1/login").WithJSON(model.LoginRequest{
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
	superAdmin := initUser("admin", "12345678", "maintainman default admin")
	apiJson := service.UserLogin(&model.LoginRequest{
		Account:  superAdmin.Name,
		Password: superAdmin.Password,
	}, getMyIPV6(), nil)
	return cast.ToString(apiJson.Data)
}

func generateRandomUsers(prefix string, num uint) (usersRegister []model.RegisterUserRequest) {
	for i := uint(1); i <= num; i++ {
		usersRegister = append(usersRegister, initUser(prefix+strconv.Itoa(int(i))+util.RandomString(5), "12345678", "disp name user"+strconv.Itoa(int(i))))
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
