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
		responseBody := e.POST("/v1/register").WithJSON(user).Expect().Status(httptest.StatusCreated).Body().Raw()
		fmt.Println(responseBody)
		u := &model.UserJson{}
		_ = json.Unmarshal([]byte(responseBody), u)
		id := u.ID

		responseBody = e.PUT("/v1/user/"+cast.ToString(id)).WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(model.UpdateUserRequest{
			Name:        user.Name + "_update",
			Password:    user.Password,
			DisplayName: user.DisplayName,
			Phone:       user.Phone,
			Email:       user.Email,
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
		responseBody := e.POST("/v1/user/").WithHeader("Authorization", "Bearer "+superAdminToken).WithJSON(model.CreateUserRequest{
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
