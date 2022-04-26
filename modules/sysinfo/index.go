package sysinfo

import (
	_ "embed"

	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"
)

var Module = module.Module{
	ModuleName:    "sysinfo",
	ModuleVersion: "1.0.0",
	ModuleDepends: []string{},
	ModuleEnv:     map[string]any{},
	ModuleExport:  map[string]any{},
	ModulePerm: map[string]string{
		"sysinfo.view": "查看系统信息",
	},
	EntryPoint: entry,
}

func entry(mctx *module.ModuleContext) {
	mctx.Route.PartyFunc("/sysinfo", func(sysinfo iris.Party) {
		sysinfo.Get("/", rbac.PermInterceptor("sysinfo.view"), getSysInfo)
	})
}

// getSysInfo godoc
// @Summary      获取系统信息
// @Description  获取Go Runtime信息
// @Tags         sysinfo
// @Produce      json
// @Success      200  {object}  RuntimeStatus
// @Router       /v1/sysinfo [get]
func getSysInfo(ctx iris.Context) {
	response := model.Success(newStatus(), "")
	ctx.Values().Set("response", response)
}
