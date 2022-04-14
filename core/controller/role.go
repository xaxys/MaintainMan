package controller

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/service"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// GetRole godoc
// @Summary 获取当前用户角色信息
// @Description 获取当前用户角色信息
// @Tags role
// @Produce  json
// @Success 200 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role [get]
func GetRole(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetRoleByName(auth.Role, auth)
	ctx.Values().Set("response", response)
}

// GetRoleByName godoc
// @Summary 获取某角色信息
// @Description 通过角色名获取某角色信息
// @Tags role
// @Produce  json
// @Param name path string true "角色名"
// @Success 200 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role/{name} [get]
func GetRoleByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetRoleByName(name, auth)
	ctx.Values().Set("response", response)
}

// CreateRole godoc
// @Summary 创建角色
// @Description 创建角色
// @Tags role
// @Accept  json
// @Produce  json
// @Param body body model.CreateRoleRequest true "创建角色请求"
// @Success 201 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role [post]
func CreateRole(ctx iris.Context) {
	aul := &model.CreateRoleRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateRole(aul, auth)
	ctx.Values().Set("response", response)
}

// UpdateRole godoc
// @Summary 更新角色
// @Description 更新角色
// @Tags role
// @Accept  json
// @Produce  json
// @Param body body model.UpdateRoleRequest true "更新角色请求"
// @Success 204 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role [put]
func UpdateRole(ctx iris.Context) {
	aul := &model.UpdateRoleRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateRole(name, aul, auth)
	ctx.Values().Set("response", response)
}

// SetDefaultRole godoc
// @Summary 设置默认角色
// @Description 设置默认角色(用户注册时的默认角色)
// @Tags role
// @Accept  json
// @Produce  json
// @Param name path string true "角色名"
// @Success 204 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role/{name}/default [put]
func SetDefaultRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.SetDefaultRole(name, auth)
	ctx.Values().Set("response", response)
}

// SetGuestRole godoc
// @Summary 设置游客角色
// @Description 设置游客角色(用户未登录时的默认角色)
// @Tags role
// @Accept  json
// @Produce  json
// @Param name path string true "角色名"
// @Success 204 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role/{name}/guest [put]
func SetGuestRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.SetGuestRole(name, auth)
	ctx.Values().Set("response", response)
}

// DeleteRole godoc
// @Summary 删除角色
// @Description 删除角色
// @Tags role
// @Accept  json
// @Produce  json
// @Param name path string true "角色名"
// @Success 204 {object} model.ApiJson{data=model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role/{name} [delete]
func DeleteRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteRole(name, auth)
	ctx.Values().Set("response", response)
}

// GetAllRoles godoc
// @Summary 获取所有角色
// @Description 获取所有角色 不分页
// @Tags role
// @Produce  json
// @Success 200 {object} model.ApiJson{data=[]model.RoleJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/role/all [get]
func GetAllRoles(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllRoles(auth)
	ctx.Values().Set("response", response)
}
