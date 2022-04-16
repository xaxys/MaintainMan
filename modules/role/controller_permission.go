package role

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

// getPermission godoc
// @Summary 获取当前用户权限信息
// @Description 获取当前用户权限信息
// @Tags permission
// @Produce  json
// @Param name path string true "权限名"
// @Success 200 {object} model.ApiJson{data=rbac.PermissionJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/permission/{name} [get]
func getPermission(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := GetPermissionService(name, auth)
	ctx.Values().Set("response", response)
}

// getAllPermissions godoc
// @Summary 获取所有权限信息
// @Description 获取所有权限信息 不分页
// @Tags permission
// @Produce  json
// @Success 200 {object} model.ApiJson{data=[]rbac.PermissionJson}
// @Failure 400 {object} model.ApiJson{data=[]string}
// @Failure 401 {object} model.ApiJson{data=[]string}
// @Failure 403 {object} model.ApiJson{data=[]string}
// @Failure 404 {object} model.ApiJson{data=[]string}
// @Failure 422 {object} model.ApiJson{data=[]string}
// @Failure 500 {object} model.ApiJson{data=[]string}
// @Router /v1/permission/all [get]
func getAllPermissions(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := GetAllPermissionsService(auth)
	ctx.Values().Set("response", response)
}
