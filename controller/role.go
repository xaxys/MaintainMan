package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetRole(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetRoleByName(auth.Role, auth)
	ctx.Values().Set("response", response)
}

func GetRoleByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetRoleByName(name, auth)
	ctx.Values().Set("response", response)
}

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

func UpdateRoleByName(ctx iris.Context) {
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

func SetDefaultRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.SetDefaultRole(name, auth)
	ctx.Values().Set("response", response)
}

func SetGuestRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.SetGuestRole(name, auth)
	ctx.Values().Set("response", response)
}

func DeleteRoleByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.DeleteRole(name, auth)
	ctx.Values().Set("response", response)
}

func GetAllRoles(ctx iris.Context) {
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllRoles(auth)
	ctx.Values().Set("response", response)
}
