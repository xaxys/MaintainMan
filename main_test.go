package main

import (
	"encoding/json"
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
	"github.com/spf13/cast"
)

func TestRegisterAndLoginRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("testUser", 10)

	for _, user := range users {
		responseBody := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated).Body()
		fmt.Println(responseBody)
	}
	for _, user := range users {
		responseBody := e.POST("/v1/login").WithJSON(model.LoginRequest{
			Account:  user.Name,
			Password: user.Password,
		}).Expect().Status(httptest.StatusOK).Body()
		fmt.Println(responseBody)
	}
}

func TestUserReNewRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	responseBody := e.GET("/v1/renew").WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK).Body()
	fmt.Println(responseBody)
}

func TestUserViewRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/user").WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK).Body()
	fmt.Println(responseBody)
}

func TestUpdateUserRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("updateUser", 10)
	superAdminToken := getSuperAdminToken()

	for _, user := range users {
		response := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
		fmt.Println(response.Body().Raw())
		u := response.JSON().NotNull().Object().Value("data")
		id := uint(u.Object().Value("id").NotNull().Raw().(float64))

		responseBody := e.PUT("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(model.UpdateUserRequest{
			Name:        user.Name + "_update",
			Password:    user.Password + "_update",
			DisplayName: user.DisplayName + "_update",
			Phone:       "",
			Email:       "",
			RoleName:    "user",
		}).Expect().Status(httptest.StatusNoContent).Body().Raw()
		fmt.Println(responseBody)
	}

}

func TestGetAllUsersRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/user/all").WithHeader("Authorization", "Bearer "+superAdminToken).WithQuery("offset", 0).WithQuery("limit", 50).Expect().Status(httptest.StatusOK).Body().Raw()

	fmt.Println(responseBody)
}

//FIXME:
func TestCreateUser(t *testing.T) {
	app := newApp()
	superAdminToken := getSuperAdminToken()
	e := httptest.New(t, app)
	users := generateRandomUsers("createUser", 10)
	for _, user := range users {
		responseBody := e.POST("/v1/user").WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(model.CreateUserRequest{
			RegisterUserRequest: user,
			RoleName:            "user",
		}).Expect().Status(httptest.StatusCreated).Body().Raw()
		fmt.Println(responseBody)
	}
}

func TestForceDeleteUser(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("deleteUser", 10)
	superAdminToken := getSuperAdminToken()

	for _, user := range users {
		responseBody := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated).Body().Raw()
		fmt.Println(responseBody)
		u := &model.UserJson{}
		_ = json.Unmarshal([]byte(responseBody), u)
		id := u.ID

		responseBody = e.DELETE("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusNoContent).Body().Raw()
		fmt.Println(responseBody)
	}
}

func TestGetUserById(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("getUser", 10)
	superAdminToken := getSuperAdminToken()

	for _, user := range users {
		responseBody := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated).Body().Raw()
		fmt.Println(responseBody)
		u := &model.UserJson{}
		_ = json.Unmarshal([]byte(responseBody), u)
		id := u.ID

		responseBody = e.GET("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK).Body().Raw()
		fmt.Println(responseBody)
	}
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

func TestTagGetByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	ids := []uint{}
	for _, tag := range getTestTags() {
		resp := service.CreateTag(&tag, getSuperAdminAuthInfo())
		ids = append(ids, resp.Data.(*model.TagJson).ID)
	}

	for _, id := range ids {
		url := fmt.Sprintf("/v1/tag/%d", id)
		resp := e.GET(url).
			WithHeader("Authorization", "Bearer "+superAdminToken).
			Expect().Status(httptest.StatusOK)
		t.Log(resp.Body().Raw())
	}
}

func TestTagDeleteRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	ids := []uint{}
	for _, tag := range getTestTags() {
		resp := service.CreateTag(&tag, getSuperAdminAuthInfo())
		ids = append(ids, resp.Data.(*model.TagJson).ID)
	}

	e.DELETE(fmt.Sprintf("/v1/tag/%d", ids[0])).
		Expect().Status(httptest.StatusForbidden)

	for _, id := range ids {
		url := fmt.Sprintf("/v1/tag/%d", id)
		resp := e.DELETE(url).
			WithHeader("Authorization", "Bearer "+superAdminToken).
			Expect().Status(httptest.StatusNoContent)
		t.Log(resp.Body().Raw())
	}
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
