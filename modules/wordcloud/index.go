package wordcloud

import (
	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/module"
)

var Module = module.Module{
	ModuleName:    "word",
	ModuleVersion: "1.0.0",
	ModuleConfig:  userConfig,
	ModuleEnv: map[string]any{
		"orm.model": []any{
			&Word{},
		},
	},
	ModuleExport: map[string]any{},
	ModulePerm: map[string]string{
		"word.upload":  "上传词库",
		"word.getall":  "获取全部词库",
		"word.getword": "获取订单词库",
	},
	EntryPoint: entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
    //TODO:增加权限管理
	mctx.Route.PartyFunc("/word", func(user iris.Party) {
		user.Get("/all", getAllWords)
		user.Get("/{id:unit}", getWordByOrderId)
		user.Post("/{id:uint}", uploadWords)
	})
}
