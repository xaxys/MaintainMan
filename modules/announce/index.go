package announce

import (
	"github.com/xaxys/maintainman/core/middleware"
	"github.com/xaxys/maintainman/module"

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
	EntryPoint:   entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	ctx.Route.PartyFunc("/announce", func(announce iris.Party) {
		announce.Get("/", middleware.PermInterceptor("announce.view"), getLatestAnnounces)
		announce.Get("/all", middleware.PermInterceptor("announce.viewall"), getAllAnnounces)
		announce.Get("/{id:uint}", middleware.PermInterceptor("announce.viewall"), getAnnounce)
		announce.Post("/", middleware.PermInterceptor("announce.create"), createAnnounce)
		announce.Put("/{id:uint}", middleware.PermInterceptor("announce.update"), updateAnnounce)
		announce.Delete("/{id:uint}", middleware.PermInterceptor("announce.delete"), deleteAnnounce)
		announce.Get("/{id:uint}/hit", middleware.PermInterceptor("announce.hit"), hitAnnounce)
	})
}
