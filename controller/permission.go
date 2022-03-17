package controller

import (
	"maintainman/service"

	"github.com/kataras/iris/v12"
)

func GetPermissionByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	response := service.GetPermission(name)
	ctx.Values().Set("response", response)
}

func GetAllPermissions(ctx iris.Context) {
	response := service.GetAllPermissions()
	ctx.Values().Set("response", response)
}
