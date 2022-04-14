package imagehost

import (
	"github.com/xaxys/maintainman/core/middleware"
	"github.com/xaxys/maintainman/module"

	"github.com/kataras/iris/v12"
)

var Module = module.Module{
	ModuleName:    "image",
	ModuleVersion: "1.0.0",
	ModuleConfig:  imageConfig,
	ModuleEnv: map[string]any{
		"cache.evict": onEvict,
	},
	ModuleExport: map[string]any{},
	EntryPoint:   entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	initLimiter()
	ctx.Route.PartyFunc("/image", func(image iris.Party) {
		image.Use(middleware.HeaderExtractor, middleware.TokenValidator)
		image.Done(middleware.ResponseHandler)
		image.SetExecutionRules(iris.ExecutionRules{Done: iris.ExecutionOptions{Force: true}})
		image.Post("/", middleware.PermInterceptor("image.upload"), rateLimiter, uploadImage)
		image.Get("/{id:uuid}", middleware.PermInterceptor("image.view"), getImage)
	})

	transformationPO = newTransformationPersistence(imageConfig)
	imageStorage = ctx.Storage
	imageCacheStorage = ctx.Storage.Sub("cache", imageConfig.GetBool("storage.cache.clean"))
}
