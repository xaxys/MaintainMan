package module

import (
	"fmt"

	"github.com/xaxys/maintainman/core/cache"
	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/logger"
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
		configKey := fmt.Sprintf("module.%s", m.ModuleName)
		if config.AppConfig.Get(configKey) != nil && !config.AppConfig.GetBool(configKey) {
			logger.Logger.Warnf("module disabled: %s", m.ModuleName)
			continue
		}
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
		for _, dep := range m.ModuleDepends {
			if r.modules[dep] == nil {
				logger.Logger.Fatalf("module %s dependencies missing: %s", m.ModuleName, dep)
			}
		}
		mctx := &ModuleContext{
			Server:  r.server,
			Route:   router.APIRoute.Party(m.ModuleRoute),
			Storage: storage.InitStorage(m.ModuleConfig),
			Cache:   cache.InitCache(m.ModuleName, m.ModuleConfig, m.getOnEvict()),
		}
		m.EntryPoint(mctx)
		logger.Logger.Infof("module loaded: %s", m.ModuleName)
	}
}

func (r *Registry) Get(moduleName string) IModule {
	return r.modules[moduleName]
}
