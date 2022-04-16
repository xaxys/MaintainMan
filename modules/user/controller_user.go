package user

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getUser godoc
// @Summary      获取当前登录用户信息
// @Description  获取当前登录用户信息 附带角色和权限信息
// @Tags         user
// @Produce      json
// @Success      200  {object}  model.ApiJson{data=user.UserJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/user [get]
func getUser(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getUserInfoByIDService(auth.User, auth)
	ctx.Values().Set("response", response)
}

// getUserByID godoc
// @Summary      获取某用户信息
// @Description  通过ID获取某用户信息
// @Tags         user
// @Produce      json
// @Param        id   path      uint  true  "用户ID"
// @Success      200  {object}  model.ApiJson{data=user.UserJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/user/{id} [get]
func getUserByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getUserInfoByIDService(id, auth)
	ctx.Values().Set("response", response)
}

// getUsersByDivision godoc
// @Summary      获取某分组下的所有用户信息
// @Description  获取某分组下的所有用户信息 分页
// @Tags         user
// @Produce      json
// @Param        id        path      uint    true   "分组ID"
// @Param        order_by  query     string  false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset    query     uint    false  "偏移量 (默认为0)"
// @Param        limit     query     uint    false  "每页数据量 (默认为50)"
// @Success      200       {object}  model.ApiJson{data=model.Page{entries=[]user.UserJson}}
// @Failure      400       {object}  model.ApiJson{data=[]string}
// @Failure      401       {object}  model.ApiJson{data=[]string}
// @Failure      403       {object}  model.ApiJson{data=[]string}
// @Failure      404       {object}  model.ApiJson{data=[]string}
// @Failure      422       {object}  model.ApiJson{data=[]string}
// @Failure      500       {object}  model.ApiJson{data=[]string}
// @Router       /v1/user/division/{id} [get]
func getUsersByDivision(ctx iris.Context) {
	param := &model.PageParam{}
	if err := ctx.ReadQuery(param); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getUsersByDivisionService(id, param, auth)
	ctx.Values().Set("response", response)
}

// getAllUsers godoc
// @Summary      获取所有用户信息
// @Description  获取所有用户信息 用户名 昵称查找 分页
// @Tags         user
// @Produce      json
// @Param        name          query     string  false  "用户名"
// @Param        display_name  query     string  false  "昵称"
// @Param        order_by      query     string  false  "排序字段 (默认为ID正序)  只接受  {field}  {asc|desc}  格式  (e.g. id desc)"
// @Param        offset        query     uint    false  "偏移量 (默认为0)"
// @Param        limit         query     uint    false  "每页数据量 (默认为50)"
// @Success      200           {object}  model.ApiJson{data=model.Page{entries=[]user.UserJson}}
// @Failure      400           {object}  model.ApiJson{data=[]string}
// @Failure      401           {object}  model.ApiJson{data=[]string}
// @Failure      403           {object}  model.ApiJson{data=[]string}
// @Failure      404           {object}  model.ApiJson{data=[]string}
// @Failure      422           {object}  model.ApiJson{data=[]string}
// @Failure      500           {object}  model.ApiJson{data=[]string}
// @Router       /v1/user/all [get]
func getAllUsers(ctx iris.Context) {
	aul := &AllUserRequest{}
	if err := ctx.ReadQuery(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getAllUsersService(aul, auth)
	ctx.Values().Set("response", response)
}

// userLogin godoc
// @Summary      用户登录
// @Description  用户登录
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest                true  "登录信息"
// @Success      200   {object}  model.ApiJson{data=string}  "JWT Token"
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/login [post]
func userLogin(ctx iris.Context) {
	aul := &LoginRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := userLoginService(aul, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

// wxUserLogin godoc
// @Summary      微信登录
// @Description  微信登录
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      WxLoginRequest              true  "登录信息"
// @Success      200   {object}  model.ApiJson{data=string}  "JWT Token"
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/wxlogin [post]
func wxUserLogin(ctx iris.Context) {
	aul := &WxLoginRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := wxUserLoginService(aul, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

// wxUserRegister godoc
// @Summary      微信注册并登陆
// @Description  微信注册并登陆
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      WxRegisterRequest           true  "登录信息"
// @Success      200   {object}  model.ApiJson{data=string}  "JWT Token"
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/wxregister [post]
func wxUserRegister(ctx iris.Context) {
	aul := &WxRegisterRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := wxUserRegisterService(aul, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

// userRenew godoc
// @Summary      用户登录续期
// @Description  用户登录续期
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.ApiJson{data=string}  "JWT Token"
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/renew [get]
func userRenew(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	id := util.NilOrBaseValue(auth, func(v *model.AuthInfo) uint { return v.User }, 0)
	response := userRenewService(id, ctx.Request().RemoteAddr, auth)
	ctx.Values().Set("response", response)
}

// userRegister godoc
// @Summary      用户注册
// @Description  用户注册
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterUserRequest  true  "注册信息"
// @Success      201   {object}  model.ApiJson{data=user.UserJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/register [post]
func userRegister(ctx iris.Context) {
	aul := &RegisterUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := registerUserService(aul, auth)
	ctx.Values().Set("response", response)
}

// createUser godoc
// @Summary      创建用户(管理员)
// @Description  创建用户 所有字段都可设置 普通用户应使用注册，而不是这个创建
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      CreateUserRequest  true  "创建信息"
// @Success      201   {object}  model.ApiJson{data=UserJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/user [post]
func createUser(ctx iris.Context) {
	aul := &CreateUserRequest{}
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := createUserService(aul, auth)
	ctx.Values().Set("response", response)
}

// updateUser godoc
// @Summary      更新当前用户
// @Description  更新当前用户 除角色和分组外其他字段可更新
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      UpdateUserRequest  true  "更新信息"
// @Success      204   {object}  model.ApiJson{data=UserJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/user [put]
func updateUser(ctx iris.Context) {
	aul := &UpdateUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.RoleName = ""
	aul.DivisionID = 0
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := updateUserService(auth.User, aul, auth)
	ctx.Values().Set("response", response)
}

// forceUpdateUser godoc
// @Summary      更新用户(管理员)
// @Description  通过ID更新用户 所有字段都可更新
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "用户ID"
// @Param        body  body      UpdateUserRequest  true  "更新信息"
// @Success      204   {object}  model.ApiJson{data=UserJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/user/{id} [put]
func forceUpdateUser(ctx iris.Context) {
	aul := &UpdateUserRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := updateUserService(id, aul, auth)
	ctx.Values().Set("response", response)
}

// forceDeleteUser godoc
// @Summary      删除用户(管理员)
// @Description  通过ID删除用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "用户ID"
// @Success      204  {object}  model.ApiJson
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/user/{id} [delete]
func forceDeleteUser(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := deleteUserService(id, auth)
	ctx.Values().Set("response", response)
}
