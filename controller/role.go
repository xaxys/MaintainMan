package controller

import (
	"maintainman/model"
	"maintainman/service"

	"github.com/kataras/iris/v12"
)

func GetRole(ctx iris.Context) {
	role := ctx.Values().GetString("user_role")
	response := service.GetRoleByName(role)
	ctx.Values().Set("response", response)
}

func GetRoleByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	response := service.GetRoleByName(name)
	ctx.Values().Set("response", response)
}

func CreateRole(ctx iris.Context) {
	aul := &model.CreateRoleJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.CreateRole(aul)
	ctx.Values().Set("response", response)
}

func UpdateRoleByName(ctx iris.Context) {
	aul := &model.UpdateRoleJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	name := ctx.Params().GetString("name")
	response := service.UpdateRole(name, aul)
	ctx.Values().Set("response", response)
}

func SetDefaultRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	response := service.SetDefaultRole(name)
	ctx.Values().Set("response", response)
}

func SetGuestRole(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	response := service.SetGuestRole(name)
	ctx.Values().Set("response", response)
}

func DeleteRoleByName(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	response := service.DeleteRole(name)
	ctx.Values().Set("response", response)
}

func GetAllRoles(ctx iris.Context) {
	response := service.GetAllRoles()
	ctx.Values().Set("response", response)
}
