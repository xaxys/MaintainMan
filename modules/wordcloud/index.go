package wordcloud

import (
	"github.com/kataras/iris/v12"
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"
)

var Module = module.Module{
	ModuleName:    "word",
	ModuleVersion: "1.0.0",
	ModuleEnv: map[string]any{
		"orm.model": []any{
			&OrderWord{},
			&GlobalWord{},
		},
	},
	ModuleExport: map[string]any{},
	ModulePerm: map[string]string{
		"word.view": "查看词云",
	},
	EntryPoint: entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	mctx.Route.PartyFunc("/word", func(word iris.Party) {
		word.Get("/", rbac.PermInterceptor("word.view"), getAllWords)
		word.Get("/{id:uint}", rbac.PermInterceptor("word.view"), getWordsByOrder)
	})
	go listener()
}
