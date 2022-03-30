package main

import (
	"fmt"
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

// Tag Router

func TestTagCreateRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	tags := getTestTags()

	e.POST("/v1/tag").
		WithJSON(tags[0]).
		Expect().Status(httptest.StatusForbidden)

	for _, tag := range tags {
		e.POST("/v1/tag").
			WithJSON(tag).
			WithHeader("Authorization", "Bearer "+superAdminToken).
			Expect().Status(httptest.StatusCreated)
	}
}

func TestTagSortsRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	for _, tag := range getTestTags() {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	e.GET("/v1/tag/sort").
		Expect().Status(httptest.StatusForbidden)

	response := e.GET("/v1/tag/sort").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusOK)
	t.Log(response.Body().Raw())

	sorts := response.JSON().NotNull().Object().Value("data").Array().NotEmpty()
	for _, sort := range sorts.Iter() {
		url := fmt.Sprintf("/v1/tag/sort/%s", sort.String().NotEmpty().Raw())
		resp := e.GET(url).
			WithHeader("Authorization", "Bearer "+superAdminToken).
			Expect().Status(httptest.StatusOK)
		t.Log(resp.Body().Raw())
	}
}

func TestTagGetBySortRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	for _, tag := range getTestTags() {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	e.GET("/v1/tag/sort/楼名").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusOK)
}

// Test Utils

func getSuperAdminToken() string {
	token, _ := util.GetJwtString(1, "super_admin")
	return token
}

func getSuperAdminAuthInfo() *model.AuthInfo {
	return &model.AuthInfo{
		User: 1,
		Role: "super_admin",
		IP:   getMyIPV6(),
	}
}

func generateRandomUsers(prefix string, num uint) (usersRegister []model.RegisterUserRequest) {
	for i := uint(1); i <= num; i++ {
		usersRegister = append(usersRegister, initUser(prefix+strconv.Itoa(int(i))+util.RandomString(5), "12345678", "disp name user"+strconv.Itoa(int(i))))
	}
	return
}

func getTestTags() []model.CreateTagRequest {
	return []model.CreateTagRequest{
		{
			Sort:  "楼名",
			Name:  "一舍",
			Level: 1,
		},
		{
			Sort:  "楼名",
			Name:  "二舍",
			Level: 1,
		},
		{
			Sort:  "楼名",
			Name:  "三舍",
			Level: 1,
		},
		{
			Sort:  "紧急程度",
			Name:  "一般",
			Level: 1,
		},
		{
			Sort:  "紧急程度",
			Name:  "紧急",
			Level: 1,
		},
		{
			Sort:  "故障类型",
			Name:  "漏水",
			Level: 2,
		},
		{
			Sort:  "故障类型",
			Name:  "电线",
			Level: 2,
		},
	}
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
