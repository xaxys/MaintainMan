package module

import (
	"github.com/xaxys/maintainman/core/cache"
	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/rbac"
	"github.com/xaxys/maintainman/core/router"
	"github.com/xaxys/maintainman/core/storage"
)

type Registry struct {
	modules map[string]*Module
	server  *Server
}

func NewRegistry(server *Server) *Registry {
	registry := &Registry{
		modules: make(map[string]*Module),
		server:  server,
	}
	server.Registry = registry
	return registry
}

func (r *Registry) Register(module ...*Module) {
	model := []any{}
	for _, m := range module {
		r.modules[m.ModuleName] = m
		model = append(model, m.getModel()...)
		rbac.RegisterPerm(m.ModuleName, m.ModulePerm)
		if m.ModuleConfig != nil {
			m.ModuleConfig.SetConfigName(m.ModuleName)
			m.ModuleConfig.SetConfigType("yaml")
			m.ModuleConfig.AddConfigPath(".")
			m.ModuleConfig.AddConfigPath("./config")
			m.ModuleConfig.AddConfigPath("/etc/maintainman/")
			m.ModuleConfig.AddConfigPath("$HOME/.maintainman/")
			config.ReadAndUpdateConfig(m.ModuleConfig, m.ModuleName, m.ModuleVersion)
		}
	}
	database.SyncModel(model...)
	for _, m := range module {
		mctx := &ModuleContext{
			Server:  r.server,
			Route:   router.APIRoute.Party(m.ModuleRoute),
			Storage: storage.InitStorage(m.ModuleConfig),
			Cache:   cache.InitCache(m.ModuleName, m.ModuleConfig, m.getOnEvict()),
		}
		m.EntryPoint(mctx)
	}
}

func (r *Registry) Get(moduleName string) IModule {
	return r.modules[moduleName]
}
