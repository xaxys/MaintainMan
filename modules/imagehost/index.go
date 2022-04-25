package imagehost

import (
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"

	"github.com/kataras/iris/v12"
)

var Module = module.Module{
	ModuleName:    "image",
	ModuleVersion: "1.1.0",
	ModuleConfig:  imageConfig,
	ModuleEnv: map[string]any{
		"cache.evict": onEvict,
	},
	ModuleExport: map[string]any{},
	ModulePerm: map[string]string{
		"image.upload": "上传图片",
		"image.view":   "查看图片",
		"image.custom": "处理图片",
	},
	EntryPoint: entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	initLimiter()
	ctx.Route.PartyFunc("/image", func(image iris.Party) {
		if rateLimiter != nil {
			image.Post("/", rbac.PermInterceptor("image.upload"), rateLimiter, uploadImage)
		} else {
			image.Post("/", rbac.PermInterceptor("image.upload"), uploadImage)
		}
		image.Get("/{id:uuid}", rbac.PermInterceptor("image.view"), getImage)
	})

	transformationPO = newTransformationPersistence(imageConfig)
	imageStorage = ctx.Storage
	imageCacheStorage = ctx.Storage.Sub("cache", imageConfig.GetBool("storage.cache.clean"))
}
