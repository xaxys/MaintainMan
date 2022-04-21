package user

import (
	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"
)

var Module module.Module

func init() {
	Module = module.Module{
		ModuleName:    "user",
		ModuleVersion: "1.0.0",
		ModuleConfig:  userConfig,
		ModuleEnv: map[string]any{
			"orm.model": []any{
				&User{},
				&Division{},
			},
		},
		ModuleExport: map[string]any{
			"appid":     "",
			"appsecret": "",
		},
		ModulePerm: map[string]string{
			"user.view":        "查看当前用户",
			"user.create":      "创建用户",
			"user.update":      "更新用户",
			"user.updateall":   "更新所有用户",
			"user.delete":      "删除用户",
			"user.viewall":     "查看所有用户",
			"user.login":       "登录",
			"user.register":    "注册",
			"user.wxlogin":     "微信登录",
			"user.wxregister":  "微信注册",
			"user.renew":       "更新Token",
			"division.viewall": "查看所有分组",
			"division.create":  "创建分组",
			"division.update":  "更新分组",
			"division.delete":  "删除分组",
		},
		EntryPoint: entry,
	}
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	initDefaultData()
	Module.ModuleExport["appid"] = userConfig.GetString("wechat.appid")
	Module.ModuleExport["appsecret"] = userConfig.GetString("wechat.secret")

	mctx.Route.Post("/login", rbac.PermInterceptor("user.login"), userLogin)
	mctx.Route.Post("/wxlogin", rbac.PermInterceptor("user.wxlogin"), wxUserLogin)
	mctx.Route.Post("/register", rbac.PermInterceptor("user.register"), userRegister)
	mctx.Route.Post("/wxregister", rbac.PermInterceptor("user.wxregister"), wxUserRegister)
	mctx.Route.Get("/renew", rbac.PermInterceptor("user.renew"), userRenew)
	mctx.Route.Get("/wxappid", getAppID)

	mctx.Route.PartyFunc("/user", func(user iris.Party) {
		user.Get("/", rbac.PermInterceptor("user.view"), getUser)
		user.Put("/", rbac.PermInterceptor("user.update"), updateUser)
		user.Post("/", rbac.PermInterceptor("user.create"), createUser)
		user.Get("/all", rbac.PermInterceptor("user.viewall"), getAllUsers)
		user.Get("/{id:uint}", rbac.PermInterceptor("user.viewall"), getUserByID)
		user.Put("/{id:uint}", rbac.PermInterceptor("user.updateall"), forceUpdateUser)
		user.Delete("/{id:uint}", rbac.PermInterceptor("user.delete"), forceDeleteUser)
		user.Get("/division/{id:uint}", rbac.PermInterceptor("user.viewall"), getUsersByDivision)
	})

	mctx.Route.PartyFunc("/division", func(division iris.Party) {
		division.Get("/{id:uint}", rbac.PermInterceptor("division.viewall"), getDivision)
		division.Get("/{id:uint}/children", rbac.PermInterceptor("division.viewall"), getDivisionsByParentID)
		division.Post("/", rbac.PermInterceptor("division.create"), createDivision)
		division.Put("/{id:uint}", rbac.PermInterceptor("division.update"), updateDivision)
		division.Delete("/{id:uint}", rbac.PermInterceptor("division.delete"), deleteDivision)
	})
}

// getAppID godoc
// @Summary 获取微信AppID
// @Description 获取微信AppID
// @Tags user
// @Produce text/plain
// @Success 200 {string} string "微信AppID"
// @Router /v1/wxappid [get]
func getAppID(ctx iris.Context) {
	ctx.Write([]byte(userConfig.GetString("wechat.appid")))
}
