package controller

import (
	"maintainman/model"
	"maintainman/service"

	"github.com/kataras/iris/v12"
)

func GetUser(ctx iris.Context) {
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.GetUserInfoByID(auth.User, auth)
	ctx.Values().Set("response", response)
}

func GetUserByID(ctx iris.Context) {
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.GetUserByID(auth.User, auth)
	ctx.Values().Set("response", response)
}

func GetAllUsers(ctx iris.Context) {
	aul := &model.AllUserRequest{}
	if err := ctx.ReadQuery(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.GetAllUsers(aul, auth)
	ctx.Values().Set("response", response)
}

func UserLogin(ctx iris.Context) {
	aul := &model.LoginRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.UserLogin(aul, auth)
	ctx.Values().Set("response", response)
}

func WxUserLogin(ctx iris.Context) {
	aul := &model.WxLoginRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.WxUserLogin(aul, auth)
	ctx.Values().Set("response", response)
}

func UserRenew(ctx iris.Context) {
	id := ctx.Values().GetUintDefault("user_id", 0)
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.UserRenew(id, auth)
	ctx.Values().Set("response", response)
}

func UserRegister(ctx iris.Context) {
	aul := &model.RegisterUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.RegisterUser(aul, auth)
	ctx.Values().Set("response", response)
}

func CreateUser(ctx iris.Context) {
	aul := &model.CreateUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.CreateUser(aul, auth)
	ctx.Values().Set("response", response)
}

func UpdateUser(ctx iris.Context) {
	aul := &model.ModifyUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.RoleName = ""
	aul.DivisionID = 0
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.UpdateUser(auth.User, aul, auth)
	ctx.Values().Set("response", response)
}

func ForceUpdateUser(ctx iris.Context) {
	aul := &model.ModifyUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("user_id", 0)
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.UpdateUser(id, aul, auth)
	ctx.Values().Set("response", response)
}

func ForceDeleteUser(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("user_id", 0)
	auth := ctx.Values().Get("auth").(*model.AuthInfo)
	response := service.DeleteUser(id, auth)
	ctx.Values().Set("response", response)
}
