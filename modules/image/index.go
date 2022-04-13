package imagehost

import (
	"maintainman/middleware"
	"maintainman/module"

	"github.com/kataras/iris/v12"
)

var Module = &module.Module{
	ModuleName:    "imagehost",
	ModuleVersion: "1.0.0",
	ModuleConfig:  ImageConfig,
	ModuleFuncs: map[string]any{
		"onEvict": onEvict,
	},
	ModuleExport: map[string]any{},
	EntryPoint:   entry,
}

func entry(server *module.Server) {
	server.Route.PartyFunc("/image", func(image iris.Party) {
		image.Use(middleware.HeaderExtractor, middleware.TokenValidator)
		image.Done(middleware.ResponseHandler)
		image.SetExecutionRules(iris.ExecutionRules{Done: iris.ExecutionOptions{Force: true}})
		image.Post("/", middleware.PermInterceptor("image.upload"), middleware.RateLimiter, UploadImage)
		image.Get("/{id:uuid}", middleware.PermInterceptor("image.view"), GetImage)
	})
}
