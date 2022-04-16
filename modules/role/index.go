package role

import (
	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"
)

var Module = module.Module{
	ModuleName:    "permission",
	ModuleVersion: "1.2.0",
	ModuleConfig:  roleConfig,
	ModuleEnv:     map[string]any{},
	ModuleExport:  map[string]any{},
	ModulePerm: map[string]string{
		"role.view":          "查看当前角色",
		"role.create":        "创建角色",
		"role.update":        "更新角色",
		"role.delete":        "删除角色",
		"role.viewall":       "查看所有角色",
		"permission.viewall": "查看所有权限",
	},
	EntryPoint: entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	rbac.LoadRole(roleConfig)
	mctx.Route.PartyFunc("/role", func(role iris.Party) {
		role.Get("/", rbac.PermInterceptor("role.view"), getRole)
		role.Post("/", rbac.PermInterceptor("role.create"), createRole)
		role.Get("/all", rbac.PermInterceptor("role.viewall"), getAllRoles)
		role.Get("/{name:string}", rbac.PermInterceptor("role.viewall"), getRoleByName)
		role.Post("/{name:string}/default", rbac.PermInterceptor("role.update"), setDefaultRole)
		role.Post("/{name:string}/guest", rbac.PermInterceptor("role.update"), setGuestRole)
		role.Put("/{name:string}", rbac.PermInterceptor("role.update"), updateRole)
		role.Delete("/{name:string}", rbac.PermInterceptor("role.delete"), deleteRole)
	})
	mctx.Route.PartyFunc("/permission", func(perm iris.Party) {
		perm.Get("/all", rbac.PermInterceptor("permission.viewall"), getAllPermissions)
		perm.Get("/{name:string}", rbac.PermInterceptor("permission.viewall"), getPermission)
	})
}
