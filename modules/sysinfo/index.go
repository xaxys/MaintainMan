package sysinfo

import (
	_ "embed"

	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/middleware"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/module"
)

var Module = module.Module{
	ModuleName:    "sysinfo",
	ModuleVersion: "1.0.0",
	ModuleEnv:     map[string]any{},
	ModuleExport:  map[string]any{},
	EntryPoint:    entry,
}

func entry(mctx *module.ModuleContext) {
	mctx.Route.PartyFunc("/sysinfo", func(sysinfo iris.Party) {
		sysinfo.Get("/", middleware.PermInterceptor("sysinfo.view"), getSysInfo)
	})
}

// getSysInfo godoc
// @Summary 获取系统信息
// @Description 获取Go Runtime信息
// @Tags sysinfo
// @Produce json
// @Success 200 {object} RuntimeStatus
// @Router /v1/sysinfo [get]
func getSysInfo(ctx iris.Context) {
	response := model.Success(newStatus(), "")
	ctx.Values().Set("response", response)
}
