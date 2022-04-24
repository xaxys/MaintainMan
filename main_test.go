package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"
	"github.com/xaxys/maintainman/modules/announce"
	"github.com/xaxys/maintainman/modules/order"
	"github.com/xaxys/maintainman/modules/user"

	"github.com/kataras/iris/v12/httptest"
	"github.com/spf13/cast"
)

// Test User Router
func TestRegisterAndLoginRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	users := generateRandomUsers("testUser", 10)

	for _, user := range users {
		responseBody := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated).Body()
		t.Log(responseBody)
	}
	for _, u := range users {
		responseBody := e.POST("/v1/login").WithJSON(user.LoginRequest{
			Account:  u.Name,
			Password: u.Password,
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

	responseBody := e.PUT("/v1/user/" + cast.ToString(id)).WithJSON(user.UpdateUserRequest{
		Name:        testUser[0].Name + "_update",
		Password:    testUser[0].Password + "_update",
		DisplayName: testUser[0].DisplayName + "_update",
		Phone:       "",
		Email:       "",
		RoleName:    "user",
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, usr := range users {
		response = e.POST("/v1/register").WithJSON(usr).Expect().Status(httptest.StatusCreated)
		t.Log(response.Body().Raw())
		u = response.JSON().NotNull().Object().Value("data")
		id = uint(u.Object().Value("id").NotNull().Raw().(float64))

		responseBody = e.PUT("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(user.UpdateUserRequest{
			Name:        usr.Name + "_update",
			Password:    usr.Password + "_update",
			DisplayName: usr.DisplayName + "_update",
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
	responseBody := e.POST("/v1/user").WithJSON(user.CreateUserRequest{
		RegisterUserRequest: testUser[0],
		RoleName:            "user",
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	for _, usr := range users {
		responseBody = e.POST("/v1/user").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(user.CreateUserRequest{
				RegisterUserRequest: usr,
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

// Test Tag Router
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
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
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
		response := e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
		id := response.JSON().NotNull().Object().Value("data").Object().Value("id").NotNull().Raw().(float64)
		ids = append(ids, uint(id))
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
	for _, tag := range generateRandomTags("randSorts ", "randName ", 30) {
		response := e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
		id := response.JSON().NotNull().Object().Value("data").Object().Value("id").NotNull().Raw().(float64)
		ids = append(ids, uint(id))
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

// Test Order Router
func TestCreateOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	orders := generateRandomOrders("Add Order", "admin", 6)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	responseBody := e.POST("/v1/order").
		WithJSON(initOrder("Test", "Test", "Earth", "Admin", 5)).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(initWrongOrder("Test", "Test", "Earth", "Admin", 5)).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	for _, order := range orders {
		t.Log(order.Tags)
		responseBody := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
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
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
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
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
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
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/user").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(user.CreateUserRequest{
			RegisterUserRequest: initUser("Test Repairer "+strconv.Itoa(rand.Intn(10000)), "12345678", "Test Repairer "+strconv.Itoa(rand.Intn(10000))),
			RoleName:            "maintainer",
		}).Expect().Status(httptest.StatusCreated)

	t.Log(response.Body().Raw())

	responseBody := e.GET("/v1/order/repairer").WithJSON(order.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/repairer").WithJSON(order.RepairerOrderRequest{
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
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/user").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(user.CreateUserRequest{
			RegisterUserRequest: initUser("Test Repairer "+strconv.Itoa(rand.Intn(10000)), "12345678", "Test Repairer "+strconv.Itoa(rand.Intn(10000))),
			RoleName:            "maintainer",
		}).Expect().Status(httptest.StatusCreated)
	u := response.JSON().NotNull().Object().Value("data")
	id := uint(u.Object().Value("id").NotNull().Raw().(float64))

	t.Log(response.Body().Raw())

	responseBody := e.GET("/v1/order/repairer/" + cast.ToString(id)).WithJSON(order.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/repairer/"+cast.ToString(id)).WithJSON(order.RepairerOrderRequest{
		Current:   true,
		PageParam: model.PageParam{},
	}).WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetOrderByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/order/1").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	responseBody = e.GET("/v1/order/1").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestUpdateOrderByUserRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.PUT("/v1/order/" + cast.ToString(id)).
		WithJSON(order.UpdateOrderRequest{
			Title:        "TestUpdateOrder " + randomNumToString,
			Content:      testOrder.Content + "_updated",
			Address:      testOrder.Address,
			ContactName:  testOrder.ContactName,
			ContactPhone: testOrder.ContactPhone,
			AddTags:      testOrder.Tags,
			DelTags:      testOrder.Tags,
		}).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	addWrongBuildingTag := 4 - testOrder.Tags[0]
	addWrongEmergencyTag := 9 - testOrder.Tags[1]
	responseBody = e.PUT("/v1/order/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.UpdateOrderRequest{
			Title:        "TestUpdateOrder " + randomNumToString,
			Content:      testOrder.Content + "_updated",
			Address:      testOrder.Address,
			ContactName:  testOrder.ContactName,
			ContactPhone: testOrder.ContactPhone,
			AddTags:      []uint{addWrongBuildingTag, addWrongEmergencyTag},
			DelTags:      []uint{},
		}).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/order/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)

	addBuildingTag := 4 - testOrder.Tags[0]
	delBuildingTag := testOrder.Tags[0]
	addEmergencyTag := 9 - testOrder.Tags[1]
	delEmergencyTag := testOrder.Tags[1]
	responseBody = e.PUT("/v1/order/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.UpdateOrderRequest{
			Title:        "TestUpdateOrder " + randomNumToString,
			Content:      testOrder.Content + "_updated",
			Address:      testOrder.Address,
			ContactName:  testOrder.ContactName,
			ContactPhone: testOrder.ContactPhone,
			AddTags:      []uint{addBuildingTag, addEmergencyTag},
			DelTags:      []uint{delBuildingTag, delEmergencyTag},
		}).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestForceUpdateOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.PUT("/v1/order/1/force").
		WithJSON(order.UpdateOrderRequest{
			Title:        "TestUpdateOrder " + randomNumToString,
			Content:      testOrder.Content + "_updated",
			Address:      testOrder.Address,
			ContactName:  testOrder.ContactName,
			ContactPhone: testOrder.ContactPhone,
			AddTags:      testOrder.Tags,
			DelTags:      testOrder.Tags,
		}).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.PUT("/v1/order/"+cast.ToString(id)+"/force").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.UpdateOrderRequest{
			Title:        "TestUpdateOrder " + randomNumToString,
			Content:      testOrder.Content + "_updated",
			Address:      testOrder.Address,
			ContactName:  testOrder.ContactName,
			ContactPhone: testOrder.ContactPhone,
			AddTags:      testOrder.Tags,
			DelTags:      testOrder.Tags,
		}).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestConsumeItemRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	response = e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	itemID := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/item/"+cast.ToString(itemID)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.AddItemRequest{
			ItemID: itemID,
			Num:    100,
			Price:  float64(rand.Intn(100)),
		}).Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/1/consume").
		WithJSON(order.ConsumeItemRequest{
			ItemID:  itemID,
			OrderID: orderID,
			Num:     uint(rand.Intn(99)),
			Price:   float64(rand.Intn(100)),
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(orderID)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(orderID)+"/assign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithQuery("repairer", 1).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(orderID)+"/consume").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.ConsumeItemRequest{
			ItemID:  itemID,
			OrderID: orderID,
			Num:     uint(rand.Intn(99)),
			Price:   float64(rand.Intn(100)),
		}).Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestReleaseOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestReleaseOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/release").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)
}

func TestAssignOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	repairerCreated := initUser("repairerCreated "+randomNumToString, "12345678", "repairer")
	response := e.POST("/v1/register").WithJSON(repairerCreated).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())

	u := response.JSON().NotNull().Object().Value("data")
	repairerId := uint(u.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.PUT("/v1/user/"+cast.ToString(repairerId)).WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(user.UpdateUserRequest{
		Name:        repairerCreated.Name + "_update",
		Password:    repairerCreated.Password + "_update",
		DisplayName: repairerCreated.DisplayName + "_update",
		Phone:       "",
		Email:       "",
		RoleName:    "maintainer",
	}).Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	testOrder := initOrder("TestAssignOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response = e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/assign").
		WithQuery("repairer", repairerId).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/assign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithQuery("repairer", repairerId).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/assign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithQuery("repairer", repairerId).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestSelfAssignOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestSelfAssignOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/selfassign").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestCompleteOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestCompleteOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/complete").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/complete").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestCancelOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestCancelOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/cancel").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/cancel").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestRejectOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestRejectOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/reject").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/reject").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestReportOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestReportOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/report").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/report").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestHoldOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestHoldOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(id) + "/hold").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/report").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/hold").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestAppraiseOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	randomNumToString := cast.ToString(rand.Intn(10000))

	testOrder := initOrder("TestAppriseOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	id := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/"+cast.ToString(id)+"/appraise").
		WithQuery("appraisal", 5).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/release").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusInternalServerError).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/selfassign").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/complete").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(id)+"/appraise").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithQuery("appraisal", 5).
		Expect().Status(httptest.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

// Test Role Router
func TestGetRoleRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/role").
		Expect().Status(httptest.StatusForbidden).
		Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/role").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusOK).
		Body().Raw()
	t.Log(responseBody)
}

// TODO: Allow Create Role Router
//func TestCreateRoleRouter(t *testing.T) {
//	app := newApp()
//	e := httptest.New(t, app)
//	randomNumToString := cast.ToString(rand.Intn(10000))
//	superAdminToken := getSuperAdminToken()
//	responseBody := e.POST("/v1/role").
//		WithJSON(model.CreateRoleRequest{
//			Name:        "Test Role " + randomNumToString,
//			DisplayName: "test_role",
//			Permissions: []string{
//				"order.create",
//			},
//			Inheritance: []string{
//				"admin",
//			},
//		}).
//		Expect().Status(httptest.StatusForbidden).
//		Body().Raw()
//	t.Log(responseBody)
//
//	responseBody = e.POST("/v1/role").
//		WithHeader("Authorization", "Bearer "+superAdminToken).
//		WithJSON(model.CreateRoleRequest{
//			Name:        "Test Role " + randomNumToString,
//			DisplayName: "test_role",
//			Permissions: []string{
//				"order.create",
//			},
//			Inheritance: []string{
//				"repairer",
//			},
//		}).
//		Expect().Status(httptest.StatusCreated).
//		Body().Raw()
//	t.Log(responseBody)
//}

func TestGetAllRoles(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	responseBody := e.GET("/v1/role/all").
		Expect().Status(httptest.StatusForbidden).
		Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/role/all").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(httptest.StatusOK).
		Body().Raw()
	t.Log(responseBody)
}

// Test Permission Router
func TestGetAllPermissionRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	responseBody := e.GET("/v1/permission/all").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/permission/all").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetPermissionRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	responseBody := e.GET("/v1/permission/role.viewall").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/permission/viewall").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

// Test Item Router
func TestCreateItemRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	responseBody := e.POST("/v1/item").
		WithJSON(itemTest).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated).Body().Raw()
	t.Log(responseBody)
}

func TestGetAllItemRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	responseBody := e.GET("/v1/item/all").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/item/all").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetItemByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	response := e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.GET("/v1/item/" + cast.ToString(id)).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/item/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestAddItemByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	response := e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/item/" + cast.ToString(id)).
		WithJSON(order.AddItemRequest{
			ItemID: id,
			Num:    uint(rand.Intn(100)),
			Price:  float64(rand.Intn(100)),
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/item/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.AddItemRequest{
			ItemID: id,
			Num:    uint(rand.Intn(100)),
			Price:  float64(rand.Intn(100)),
		}).Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestDeleteItemByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	response := e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.DELETE("/v1/item/" + cast.ToString(id)).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.DELETE("/v1/item/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestGetItemByNameRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	response := e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := item.Object().Value("name").NotNull().Raw().(string)

	responseBody := e.GET("/v1/item/name/" + id).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/item/name/"+id).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestGetItemByNameFuzzyRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	itemTest := order.CreateItemRequest{
		Name:        "test_item" + randomNumToString,
		Discription: "test_item",
	}

	response := e.POST("/v1/item").
		WithJSON(itemTest).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := item.Object().Value("name").NotNull().Raw().(string)

	responseBody := e.GET("/v1/item/name/" + id + "/fuzzy").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/item/name/"+id+"/fuzzy").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

// Test Comment Router
func TestCreateCommentRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(orderID) + "/comment").
		WithJSON(order.CreateCommentRequest{
			Content: "comment " + randomNumToString,
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(orderID)+"/comment").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.CreateCommentRequest{
			Content: "comment " + randomNumToString,
		}).Expect().Status(httptest.StatusCreated).Body().Raw()
	t.Log(responseBody)
}

func TestGetCommentsByOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	comments := generateRandomComments("test ", 10)

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	for _, comment := range comments {
		responseBody := e.POST("/v1/order/"+cast.ToString(orderID)+"/comment").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(comment).Expect().Status(httptest.StatusCreated).Body().Raw()
		t.Log(responseBody)
	}

	response = e.GET("/v1/order/" + cast.ToString(orderID) + "/comment").
		WithJSON(testOrder).Expect().Status(httptest.StatusForbidden)
	t.Log(response.Body().Raw())

	response = e.GET("/v1/order/"+cast.ToString(orderID)+"/comment").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusOK)
	t.Log(response.Body().Raw())
}

func TestForceCreateCommentRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.POST("/v1/order/" + cast.ToString(orderID) + "/comment/force").
		WithJSON(order.CreateCommentRequest{
			Content: "comment " + randomNumToString,
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/order/"+cast.ToString(orderID)+"/comment/force").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.CreateCommentRequest{
			Content: "comment " + randomNumToString,
		}).Expect().Status(httptest.StatusCreated).Body().Raw()
	t.Log(responseBody)
}

func TestForceGetCommentsByOrderRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}
	comments := generateRandomComments("test ", 10)

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	for _, comment := range comments {
		responseBody := e.POST("/v1/order/"+cast.ToString(orderID)+"/comment").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(comment).Expect().Status(httptest.StatusCreated).Body().Raw()
		t.Log(responseBody)
	}

	response = e.GET("/v1/order/" + cast.ToString(orderID) + "/comment/force").
		WithJSON(testOrder).Expect().Status(httptest.StatusForbidden)
	t.Log(response.Body().Raw())

	response = e.GET("/v1/order/"+cast.ToString(orderID)+"/comment/force").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusOK)
	t.Log(response.Body().Raw())
}

func TestDeleteCommentRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	response = e.POST("/v1/order/"+cast.ToString(orderID)+"/comment").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.CreateCommentRequest{
			Content: "comment " + randomNumToString,
		}).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	commentCreated := response.JSON().NotNull().Object().Value("data")
	commentID := uint(commentCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.DELETE("/v1/comment/" + cast.ToString(commentID)).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.DELETE("/v1/comment/"+cast.ToString(commentID)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestForceDeleteCommentRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	testOrder := initOrder("TestUpdateOrder "+randomNumToString, "Test", "Earth", "Admin", 5)
	tags := getTestTags()
	for _, tag := range tags {
		e.POST("/v1/tag").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(tag).
			Expect().Status(httptest.StatusCreated)
	}

	response := e.POST("/v1/order").WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(testOrder).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	orderCreated := response.JSON().NotNull().Object().Value("data")
	orderID := uint(orderCreated.Object().Value("id").NotNull().Raw().(float64))

	response = e.POST("/v1/order/"+cast.ToString(orderID)+"/comment").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(order.CreateCommentRequest{
			Content: "comment " + randomNumToString,
		}).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body().Raw())
	commentCreated := response.JSON().NotNull().Object().Value("data")
	commentID := uint(commentCreated.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.DELETE("/v1/comment/" + cast.ToString(commentID) + "/force").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.DELETE("/v1/comment/"+cast.ToString(commentID)+"/force").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

// Test Announce Router

func TestGetLatestAnnouncesRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	responseBody := e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusCreated).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestCreateAnnounceRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	responseBody := e.POST("/v1/announce").
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusCreated).Body().Raw()
	t.Log(responseBody)
}

func TestGetAnnounceByIDRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	response := e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.GET("/v1/announce/" + cast.ToString(id)).
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)

}

func TestGetAllAnnouncesRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))
	responseBody := e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 200000,
		}).Expect().Status(http.StatusCreated).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce/all").
		WithQueryObject(announce.AllAnnounceRequest{
			Title:     randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()) - 100000,
			EndTime:   cast.ToInt64(time.Now().Unix()) + 100000,
			Inclusive: true,
			PageParam: model.PageParam{},
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce/all").
		WithQueryObject(announce.AllAnnounceRequest{
			Title:     "",
			StartTime: cast.ToInt64(time.Now().Unix()) - 100000,
			EndTime:   cast.ToInt64(time.Now().Unix()) + 100000,
			Inclusive: true,
			PageParam: model.PageParam{},
		}).WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce/all").
		WithQueryObject(announce.AllAnnounceRequest{
			Title:     "",
			StartTime: -1,
			EndTime:   -1,
			Inclusive: true,
			PageParam: model.PageParam{},
		}).WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestUpdateAnnounceRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	response := e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.PUT("/v1/announce/" + cast.ToString(id)).
		WithJSON(announce.UpdateAnnounceRequest{
			Title:     randomNumToString + "_updated",
			Content:   randomNumToString + "_updated",
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.PUT("/v1/announce/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.UpdateAnnounceRequest{
			Title:     randomNumToString + "_updated",
			Content:   randomNumToString + "_updated",
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestDeleteAnnounceRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	response := e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.DELETE("/v1/announce/" + cast.ToString(id)).
		WithJSON(announce.UpdateAnnounceRequest{
			Title:     randomNumToString + "_updated",
			Content:   randomNumToString + "_updated",
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.DELETE("/v1/announce/"+cast.ToString(id)).
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.UpdateAnnounceRequest{
			Title:     randomNumToString + "_updated",
			Content:   randomNumToString + "_updated",
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusNoContent).Body().Raw()
	t.Log(responseBody)
}

func TestHitAnnounceRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()
	randomNumToString := cast.ToString(rand.Intn(10000))

	response := e.POST("/v1/announce").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		WithJSON(announce.CreateAnnounceRequest{
			Title:     randomNumToString,
			Content:   randomNumToString,
			StartTime: cast.ToInt64(time.Now().Unix()),
			EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
		}).Expect().Status(http.StatusCreated)
	t.Log(response.Body().Raw())

	item := response.JSON().NotNull().Object().Value("data")
	id := uint(item.Object().Value("id").NotNull().Raw().(float64))

	responseBody := e.GET("/v1/announce/" + cast.ToString(id) + "/hit").
		Expect().Status(httptest.StatusForbidden).Body().Raw()
	t.Log(responseBody)

	responseBody = e.GET("/v1/announce/"+cast.ToString(id)+"/hit").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusNoContent).Body().Raw()

	responseBody = e.GET("/v1/announce/"+cast.ToString(id)+"/hit").
		WithHeader("Authorization", "Bearer "+superAdminToken).
		Expect().Status(http.StatusOK).Body().Raw()
	t.Log(responseBody)
}

func TestMultiHitAnnounceRouter(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	superAdminToken := getSuperAdminToken()

	aids := []uint{}
	for i := 0; i < 1000; i++ {
		response := e.POST("/v1/announce").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(announce.CreateAnnounceRequest{
				Title:     util.RandomString(100),
				Content:   util.RandomString(100),
				StartTime: cast.ToInt64(time.Now().Unix()),
				EndTime:   cast.ToInt64(time.Now().Unix()) + 10000,
			}).Expect().Status(httptest.StatusCreated)
		id := uint(response.JSON().Object().Value("data").Object().Value("id").NotNull().Raw().(float64))
		aids = append(aids, uint(id))
		if i%10 == 0 {
			t.Log("create announce: ", i)
		}
	}
	t.Log("create 1000 announces")

	uids := []uint{}
	users := generateRandomUsers("AnnounceHitTest", 25)
	for i, user := range users {
		response := e.POST("/v1/user").
			WithHeader("Authorization", "Bearer "+superAdminToken).
			WithJSON(user).Expect().Status(httptest.StatusCreated)
		u := response.JSON().NotNull().Object().Value("data")
		id := uint(u.Object().Value("id").NotNull().Raw().(float64))
		uids = append(uids, id)
		if i%10 == 0 {
			t.Log("create user: ", i)
		}
	}
	t.Log("create 25 users")

	for i, uid := range uids {
		token, err := util.GetJwtString(uid, "", "user")
		if err != nil {
			t.Error(err)
		}
		for j, aid := range aids {
			e.GET("/v1/announce/"+cast.ToString(aid)+"/hit").
				WithHeader("Authorization", "Bearer "+token).
				Expect().Status(http.StatusNoContent).Body().Raw()

			if j%50 == 0 {
				t.Logf("%d/%d", i, j)
			}
		}
	}
	t.Log("all users hit all announces")

	for i, uid := range uids {
		token, err := util.GetJwtString(uid, "", "user")
		if err != nil {
			t.Error(err)
		}
		for j, aid := range aids {
			e.GET("/v1/announce/"+cast.ToString(aid)+"/hit").
				WithHeader("Authorization", "Bearer "+token).
				Expect().Status(http.StatusOK).Body().Raw()

			if j%50 == 0 {
				t.Logf("%d/%d", i, j)
			}
		}
	}
	t.Log("all users hit all announces again")
}

// Test Utils
func getSuperAdminToken() string {
	token, _ := util.GetJwtString(1, "fake super admin", "super_admin")
	return token
}

func getSuperAdminAuthInfo() *model.AuthInfo {
	return &model.AuthInfo{
		User: 1,
		Role: "super_admin",
		IP:   getMyIPV6(),
	}
}

func getTestTags() []order.CreateTagRequest {
	return []order.CreateTagRequest{
		{
			Sort:     "楼名",
			Name:     "一舍",
			Level:    1,
			Congener: 1,
		},
		{
			Sort:     "楼名",
			Name:     "二舍",
			Level:    1,
			Congener: 1,
		},
		{
			Sort:     "楼名",
			Name:     "三舍",
			Level:    1,
			Congener: 1,
		},
		{
			Sort:     "紧急程度",
			Name:     "一般",
			Level:    1,
			Congener: 1,
		},
		{
			Sort:     "紧急程度",
			Name:     "紧急",
			Level:    1,
			Congener: 1,
		},
		{
			Sort:     "故障类型",
			Name:     "漏水",
			Level:    2,
			Congener: 0,
		},
		{
			Sort:     "故障类型",
			Name:     "漏电",
			Level:    2,
			Congener: 0,
		},
		{
			Sort:     "故障类型",
			Name:     "外墙脱落",
			Level:    2,
			Congener: 0,
		},
		{
			Sort:     "故障类型",
			Name:     "锁坏了",
			Level:    2,
			Congener: 0,
		},
		{
			Sort:     "故障类型",
			Name:     "灯坏了",
			Level:    2,
			Congener: 0,
		},
		{
			Sort:     "故障类型",
			Name:     "插头没电",
			Level:    2,
			Congener: 0,
		},
		{
			Sort:     "故障类型",
			Name:     "撤硕没水",
			Level:    2,
			Congener: 0,
		},
	}
}

func generateRandomUsers(prefix string, num uint) (usersRegister []user.RegisterUserRequest) {
	for i := uint(1); i <= num; i++ {
		usersRegister = append(usersRegister, initUser(prefix+strconv.Itoa(int(i))+util.RandomString(5), "12345678", "Random name user"+strconv.Itoa(int(i))))
	}
	return
}

func generateRandomComments(prefix string, num uint) (comments []order.CreateCommentRequest) {
	for i := uint(1); i <= num; i++ {
		comments = append(comments, initComment(prefix))
	}
	return
}

func generateRandomTags(baseSort, baseName string, num uint) (tags []order.CreateTagRequest) {
	for i := uint(1); i <= num; i++ {
		tags = append(tags, initTag(baseSort+util.RandomString(1), baseName+util.RandomString(5), uint(rand.Int()%64)))
	}
	return
}

func generateRandomOrders(baseTitle string, baseName string, num uint) (orders []order.CreateOrderRequest) {
	for i := uint(1); i <= num; i++ {
		orders = append(orders, initOrder(baseTitle, "Content:"+strconv.Itoa(rand.Intn(100000)), "Address:"+strconv.Itoa(rand.Intn(100000)), baseName, uint(rand.Int63n(7))))
	}
	return
}

func initUser(name string, password string, displayName string) user.RegisterUserRequest {
	return user.RegisterUserRequest{
		Name:        name,
		Password:    password,
		DisplayName: displayName,
		Phone:       strconv.Itoa(rand.Intn(100000)),
		Email:       strconv.Itoa(rand.Intn(100000)) + "@qq.com",
		RealName:    name,
	}
}

func initTag(sort, name string, level uint) order.CreateTagRequest {
	return order.CreateTagRequest{
		Sort:  sort,
		Name:  name,
		Level: level,
	}
}

func initWrongOrder(title string, content string, address string, name string, maxTagID uint) order.CreateOrderRequest {
	tags := make([]uint, 0)
	tags = append(tags, 1, 2)
	tags = append(tags, 4, 5)
	for i := uint(1 + 5); i <= maxTagID+5; i++ {
		tags = append(tags, i)
	}
	return order.CreateOrderRequest{
		Title:        title,
		Content:      content,
		Address:      address,
		ContactName:  name,
		ContactPhone: strconv.Itoa(rand.Intn(100000)),
		Tags:         tags,
	}
}

func initOrder(title string, content string, address string, name string, maxTagID uint) order.CreateOrderRequest {
	tags := make([]uint, 0)
	building := rand.Intn(3) + 1
	tags = append(tags, uint(building))
	emergency := rand.Intn(2) + 4
	tags = append(tags, uint(emergency))
	for i := uint(1 + 5); i <= maxTagID+5; i++ {
		tags = append(tags, i)
	}
	return order.CreateOrderRequest{
		Title:        title,
		Content:      content,
		Address:      address,
		ContactName:  name,
		ContactPhone: strconv.Itoa(rand.Intn(100000)),
		Tags:         tags,
	}
}

func initComment(baseContent string) order.CreateCommentRequest {
	return order.CreateCommentRequest{
		Content: baseContent + strconv.Itoa(rand.Intn(100000)),
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
