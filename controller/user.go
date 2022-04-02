package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

// GetUser godoc
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户信息 附带角色和权限信息
// @Tags user
// @Produce  json
// @Success 200 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user [get]
func GetUser(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetUserInfoByID(auth.User, auth)
	ctx.Values().Set("response", response)
}

// GetUserByID godoc
// @Summary 获取某用户信息
// @Description 通过ID获取某用户信息
// @Tags user
// @Produce  json
// @Param id path uint true "用户ID"
// @Success 200 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user/{id} [get]
func GetUserByID(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetUserByID(auth.User, auth)
	ctx.Values().Set("response", response)
}

// GetAllUsers godoc
// @Summary 获取所有用户信息
// @Description 获取所有用户信息 用户名 昵称查找 分页
// @Tags user
// @Produce  json
// @Param name query string false "用户名"
// @Param display_name query string false "昵称"
// @Param order_by query string false "排序字段 (默认为ID正序) 只接受"{field} {asc|desc}"格式 (e.g. "id desc")"
// @Param offset query uint false "偏移量 (默认为0)"
// @Param limit query uint false "每页数据量 (默认为50)"
// @Success 200 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user/all [get]
func GetAllUsers(ctx iris.Context) {
	aul := &model.AllUserRequest{}
	if err := ctx.ReadQuery(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllUsers(aul, auth)
	ctx.Values().Set("response", response)
}

// UserLogin godoc
// @Summary 用户登录
// @Description 用户登录
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body model.LoginRequest true "登录信息"
// @Success 200 {object} model.ApiJson{data=string} "JWT Token"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/login [post]
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

// WxUserLogin godoc
// @Summary 微信登录
// @Description 微信登录
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body model.WxLoginRequest true "登录信息"
// @Success 200 {object} model.ApiJson{data=string} "JWT Token"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/wxlogin [post]
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

// WxUserRegister godoc
// @Summary 微信注册并登陆
// @Description 微信注册并登陆
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body model.WxRegisterRequest true "登录信息"
// @Success 200 {object} model.ApiJson{data=string} "JWT Token"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/wxregister [post]
func WxUserRegister(ctx iris.Context) {
	aul := &model.WxRegisterRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.WxUserRegister(aul, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

// UserRenew godoc
// @Summary 用户登录续期
// @Description 用户登录续期
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ApiJson{data=string} "JWT Token"
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/renew [get]
func UserRenew(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	id := util.NilOrBaseValue(auth, func(v *model.AuthInfo) uint { return v.User }, 0)
	response := service.UserRenew(id, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

// UserRegister godoc
// @Summary 用户注册
// @Description 用户注册
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body model.RegisterUserRequest true "注册信息"
// @Success 201 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/register [post]
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

// CreateUser godoc
// @Summary 创建用户(管理员)
// @Description 创建用户 所有字段都可设置 普通用户应使用注册，而不是这个创建
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body model.CreateUserRequest true "创建信息"
// @Success 201 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user [post]
func CreateUser(ctx iris.Context) {
	aul := &model.CreateUserRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateUser(aul, auth)
	ctx.Values().Set("response", response)
}

// UpdateUser godoc
// @Summary 更新当前用户
// @Description 更新当前用户 除角色和分组外其他字段可更新
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body model.UpdateUserRequest true "更新信息"
// @Success 204 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user [put]
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

// ForceUpdateUser godoc
// @Summary 更新用户(管理员)
// @Description 通过ID更新用户 所有字段都可更新
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "用户ID"
// @Param body body model.UpdateUserRequest true "更新信息"
// @Success 204 {object} model.ApiJson{data=model.UserJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user/{id} [put]
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

// ForceDeleteUser godoc
// @Summary 删除用户(管理员)
// @Description 通过ID删除用户
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "用户ID"
// @Success 204 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/user/{id} [delete]
func ForceDeleteUser(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteUser(id, auth)
	ctx.Values().Set("response", response)
}
