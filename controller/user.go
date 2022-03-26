package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetUser(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetUserInfoByID(auth.User, auth)
	ctx.Values().Set("response", response)
}

func GetUserByID(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetUserByID(auth.User, auth)
	ctx.Values().Set("response", response)
}

func GetAllUsers(ctx iris.Context) {
	aul := &model.AllUserRequest{}
	if err := ctx.ReadQuery(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllUsers(aul, auth)
	ctx.Values().Set("response", response)
}

func UserLogin(ctx iris.Context) {
	aul := &model.LoginRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UserLogin(aul, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}
func WxUserLogin(ctx iris.Context) {
	aul := &model.WxLoginRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.WxUserLogin(aul, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

func UserRenew(ctx iris.Context) {
	id := ctx.Values().GetUintDefault("user_id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UserRenew(id, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

func UserRegister(ctx iris.Context) {
	aul := &model.RegisterUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.RegisterUser(aul, auth)
	ctx.Values().Set("response", response)
}

func CreateUser(ctx iris.Context) {
	aul := &model.CreateUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateUser(aul, auth)
	ctx.Values().Set("response", response)
}

func UpdateUser(ctx iris.Context) {
	aul := &model.UpdateUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.RoleName = ""
	aul.DivisionID = 0
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateUser(auth.User, aul, auth)
	ctx.Values().Set("response", response)
}

func ForceUpdateUser(ctx iris.Context) {
	aul := &model.UpdateUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateUser(id, aul, auth)
	ctx.Values().Set("response", response)
}

func ForceDeleteUser(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteUser(id, auth)
	ctx.Values().Set("response", response)
}
