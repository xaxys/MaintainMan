package announce

import (
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"

	"github.com/kataras/iris/v12"
)

var Module = module.Module{
	ModuleName:    "announce",
	ModuleVersion: "1.0.0",
	ModuleConfig:  announceConfig,
	ModuleEnv: map[string]any{
		"orm.model": []any{
			&Announce{},
		},
	},
	ModuleExport: map[string]any{},
	ModulePerm: map[string]string{
		"announce.view":    "查看公告",
		"announce.hit":     "点击公告",
		"announce.create":  "创建公告",
		"announce.update":  "更新公告",
		"announce.delete":  "删除公告",
		"announce.viewall": "查看所有公告",
	},
	EntryPoint: entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	ctx.Route.PartyFunc("/announce", func(announce iris.Party) {
		announce.Get("/", rbac.PermInterceptor("announce.view"), getLatestAnnounces)
		announce.Get("/all", rbac.PermInterceptor("announce.viewall"), getAllAnnounces)
		announce.Get("/{id:uint}", rbac.PermInterceptor("announce.viewall"), getAnnounce)
		announce.Post("/", rbac.PermInterceptor("announce.create"), createAnnounce)
		announce.Put("/{id:uint}", rbac.PermInterceptor("announce.update"), updateAnnounce)
		announce.Delete("/{id:uint}", rbac.PermInterceptor("announce.delete"), deleteAnnounce)
		announce.Get("/{id:uint}/hit", rbac.PermInterceptor("announce.hit"), hitAnnounce)
	})
}
