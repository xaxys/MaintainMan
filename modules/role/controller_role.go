package role

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/rbac"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getRole godoc
// @Summary      获取当前用户角色信息
// @Description  获取当前用户角色信息
// @Tags         role
// @Produce      json
// @Success      200  {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/role [get]
func getRole(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getRoleByNameService(auth.Role, auth)
	ctx.Values().Set("response", response)
}

// getRoleByName godoc
// @Summary      获取某角色信息
// @Description  通过角色名获取某角色信息
// @Tags         role
// @Produce      json
// @Param        name  path      string  true  "角色名"
// @Success      200   {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/role/{name} [get]
func getRoleByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getRoleByNameService(name, auth)
	ctx.Values().Set("response", response)
}

// createRole godoc
// @Summary      创建角色
// @Description  创建角色
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        body  body      rbac.CreateRoleRequest  true  "创建角色请求"
// @Success      201   {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/role [post]
func createRole(ctx iris.Context) {
	aul := &rbac.CreateRoleRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := createRoleService(aul, auth)
	ctx.Values().Set("response", response)
}

// updateRole godoc
// @Summary      更新角色
// @Description  更新角色
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        body  body      rbac.UpdateRoleRequest  true  "更新角色请求"
// @Success      204   {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/role [put]
func updateRole(ctx iris.Context) {
	aul := &rbac.UpdateRoleRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := updateRoleService(name, aul, auth)
	ctx.Values().Set("response", response)
}

// setDefaultRole godoc
// @Summary      设置默认角色
// @Description  设置默认角色(用户注册时的默认角色)
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "角色名"
// @Success      204   {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/role/{name}/default [put]
func setDefaultRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := setDefaultRoleService(name, auth)
	ctx.Values().Set("response", response)
}

// setGuestRole godoc
// @Summary      设置游客角色
// @Description  设置游客角色(用户未登录时的默认角色)
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "角色名"
// @Success      204   {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/role/{name}/guest [put]
func setGuestRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := setGuestRoleService(name, auth)
	ctx.Values().Set("response", response)
}

// deleteRole godoc
// @Summary      删除角色
// @Description  删除角色
// @Tags         role
// @Accept       json
// @Produce      json
// @Param        name  path      string  true  "角色名"
// @Success      204   {object}  model.ApiJson{data=rbac.RoleJson}
// @Failure      400   {object}  model.ApiJson{data=[]string}
// @Failure      401   {object}  model.ApiJson{data=[]string}
// @Failure      403   {object}  model.ApiJson{data=[]string}
// @Failure      404   {object}  model.ApiJson{data=[]string}
// @Failure      422   {object}  model.ApiJson{data=[]string}
// @Failure      500   {object}  model.ApiJson{data=[]string}
// @Router       /v1/role/{name} [delete]
func deleteRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := deleteRoleService(name, auth)
	ctx.Values().Set("response", response)
}

// getAllRoles godoc
// @Summary      获取所有角色
// @Description  获取所有角色 不分页
// @Tags         role
// @Produce      json
// @Success      200  {object}  model.ApiJson{data=[]rbac.RoleJson}
// @Failure      400  {object}  model.ApiJson{data=[]string}
// @Failure      401  {object}  model.ApiJson{data=[]string}
// @Failure      403  {object}  model.ApiJson{data=[]string}
// @Failure      404  {object}  model.ApiJson{data=[]string}
// @Failure      422  {object}  model.ApiJson{data=[]string}
// @Failure      500  {object}  model.ApiJson{data=[]string}
// @Router       /v1/role/all [get]
func getAllRoles(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := getAllRolesService(auth)
	ctx.Values().Set("response", response)
}
