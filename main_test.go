package main

import (
	"fmt"
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/kataras/iris/v12/httptest"
	"github.com/spf13/cast"
)

//Test User Router
func TestRegisterAndLoginRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("testUser", 10)

	for _, user := range users {
		responseBody := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated).Body()
		t.Log(responseBody)
	}
	for _, user := range users {
		responseBody := e.POST("/v1/login").WithJSON(model.LoginRequest{
			Account:  user.Name,
			Password: user.Password,
		}).Expect().Status(httptest.StatusOK).Body()
		t.Log(responseBody)
	}
}

func TestUserReNewRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	responseBody := e.GET("/v1/renew").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/renew").WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestUserViewRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	responseBody := e.GET("/v1/user").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/user").WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestUpdateUserRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("updateUser", 10)
	superAdminToken := getSuperAdminToken()

	testUser := generateRandomUsers("forbidUser", 1)
	response := e.POST("/v1/register").WithJSON(testUser[0]).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	u := response.JSON().NotNull().Object().Value("data")
	id := uint(u.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.PUT("/v1/user/" + cast.ToString(id)).WithJSON(model.UpdateUserRequest{
		Name:        testUser[0].Name + "_update",
		Password:    testUser[0].Password + "_update",
		DisplayName: testUser[0].DisplayName + "_update",
		Phone:       "",
		Email:       "",
		RoleName:    "user",
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, user := range users {
		response = e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
		fmt.Println(response.Body().Raw())
		u = response.JSON().NotNull().Object().Value("data")
		id = uint(u.Object().Value("id").NotNull().Raw().(float64))

		responseBody = e.PUT("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(model.UpdateUserRequest{
			Name:        user.Name + "_update",
			Password:    user.Password + "_update",
			DisplayName: user.DisplayName + "_update",
			Phone:       "",
			Email:       "",
			RoleName:    "user",
		}).Expect().Status(httptest.StatusNoContent).Body().Raw()
		t.Log(responseBody)
	}

}

func TestGetAllUsersRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/user/all").WithHeader("Authorization", "Bearer "+superAdminToken).WithQuery("offset", 0).WithQuery("limit", 50).Expect().Status(httptest.StatusOK).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/user/all").WithQuery("offset", 0).WithQuery("limit", 50).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)
}

func TestCreateUserRouter(t *testing.T) {
	app := newApp()
	superAdminToken := getSuperAdminToken()
	e := httptest.New(t, app)
	users := generateRandomUsers("createUser", 10)
	testUser := generateRandomUsers("forbidUser", 1)
	responseBody := e.POST("/v1/user").WithJSON(model.CreateUserRequest{
		RegisterUserRequest: testUser[0],
		RoleName:            "user",
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, user := range users {
		responseBody = e.POST("/v1/user").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(model.CreateUserRequest{
				RegisterUserRequest: user,
				RoleName:            "user",
			}).Expect().Status(httptest.StatusCreated).Body().Raw()
		t.Log(responseBody)
	}

}

func TestForceDeleteUserRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("deleteUser", 10)
	superAdminToken := getSuperAdminToken()

	testUser := generateRandomUsers("forbidUser", 1)
	response := e.POST("/v1/register").WithJSON(testUser[0]).Expect().Status(httptest.StatusCreated)
	fmt.Println(response.Body().Raw())
	u := response.JSON().NotNull().Object().Value("data")
	id := uint(u.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.DELETE("/v1/user/" + cast.ToString(id)).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, user := range users {
		response = e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
		fmt.Println(response.Body().Raw())
		u = response.JSON().NotNull().Object().Value("data")
		id = uint(u.Object().Value("id").NotNull().Raw().(float64))

		responseBody = e.DELETE("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusNoContent).Body().Raw()
		t.Log(responseBody)
	}
}

func TestGetUserByIdRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("getUser", 10)
	superAdminToken := getSuperAdminToken()

	testUser := generateRandomUsers("forbidUser", 1)
	response := e.POST("/v1/register").WithJSON(testUser[0]).Expect().Status(httptest.StatusCreated)
	fmt.Println(response.Body().Raw())
	u := response.JSON().NotNull().Object().Value("data")
	id := uint(u.Object().Value("id").NotNull().Raw().(float64))
	responseBody := e.GET("/v1/user/" + cast.ToString(id)).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, user := range users {
		response = e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated)
		fmt.Println(response.Body().Raw())
		u = response.JSON().NotNull().Object().Value("data")
		id = uint(u.Object().Value("id").NotNull().Raw().(float64))
		responseBody = e.GET("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).Expect().Status(httptest.StatusOK).Body().Raw()
		t.Log(responseBody)
	}
}

//Test Tag Router
func TestTagCreateRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	tags := getTestTags()

	responseBody := e.POST("/v1/tag").WithJSON(tags[0]).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, tag := range tags {
		responseBody = e.POST("/v1/tag").
			WithJSON(tag).
			WithHeader("Authorization", "Bearer "+superAdminToken).
			Expect().Status(httptest.StatusCreated).Body().Raw()
		t.Log(responseBody)
	}
}

func TestTagSortsRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	for _, tag := range getTestTags() {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	responseBody := e.GET("/v1/tag/sort").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

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

//Test Order Router
func TestCreateOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	orders := generateRandomOrders("Add Order", "admin", 6)
	tags := getTestTags()
	for _, tag := range tags {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	responseBody := e.POST("/v1/order").
		WithJSON(initOrder("Test", "Test", "Earth", "Admin", 5)).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, order := range orders {
		responseBody = e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(order).Expect().Status(httptest.StatusCreated).Body().Raw()
		t.Log(responseBody)
	}

}

func TestGetUserOrdersRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}
	responseBody := e.GET("/v1/order/user").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/user").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetAllOrdersRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}
	responseBody := e.GET("/v1/order/all").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/all").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetRepairerOrdersRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	response := e.POST("/v1/user").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(model.CreateUserRequest{
			RegisterUserRequest: initUser("Test Repairer "+strconv.Itoa(rand.Intn(10000)), "12345678", "Test Repairer "+strconv.Itoa(rand.Intn(10000))),
			RoleName:            "maintainer",
		}).Expect().Status(httptest.StatusCreated)

	t.Log(response.Body().Raw())

	responseBody := e.GET("/v1/order/repairer").WithJSON(model.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/repairer").WithJSON(model.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetRepairerOrdersByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	response := e.POST("/v1/user").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(model.CreateUserRequest{
			RegisterUserRequest: initUser("Test Repairer "+strconv.Itoa(rand.Intn(10000)), "12345678", "Test Repairer "+strconv.Itoa(rand.Intn(10000))),
			RoleName:            "maintainer",
		}).Expect().Status(httptest.StatusCreated)
	u := response.JSON().NotNull().Object().Value("data")
	id := uint(u.Object().Value("id").NotNull().Raw().(float64))

	t.Log(response.Body().Raw())

	responseBody := e.GET("/v1/order/repairer/" + cast.ToString(id)).WithJSON(model.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/repairer/"+cast.ToString(id)).WithJSON(model.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetOrderByIdRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/order/1").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)
	tags := getTestTags()
	for _, tag := range tags {
		service.CreateTag(&tag, getSuperAdminAuthInfo())
	}

	responseBody = e.GET("/v1/order/1").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetOrderByUserRouter(t *testing.T) {

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
		usersRegister = append(usersRegister, initUser(prefix+strconv.Itoa(int(i))+util.RandomString(5), "12345678", "Random name user"+strconv.Itoa(int(i))))
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

func generateRandomOrders(baseTitle string, baseName string, num uint) (orders []model.CreateOrderRequest) {
	for i := uint(1); i <= num; i++ {
		orders = append(orders, initOrder(baseTitle, "Content:"+strconv.Itoa(rand.Intn(100000)), "Address:"+strconv.Itoa(rand.Intn(100000)), baseName, uint(rand.Int63n(7))))
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

func initOrder(title string, content string, address string, name string, maxTagID uint) model.CreateOrderRequest {
	tags := make([]uint, 0)
	for i := uint(1); i <= maxTagID; i++ {
		tags = append(tags, i)
	}
	return model.CreateOrderRequest{
		Title:        title,
		Content:      content,
		Address:      address,
		ContactName:  name,
		ContactPhone: strconv.Itoa(rand.Intn(100000)),
		Tags:         tags,
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
