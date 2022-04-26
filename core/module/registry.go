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
	"github.com/xaxys/maintainman/core/util"
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
	disabledModule := []*Module{}
	enabledModule := []*Module{}

	// register module
	for _, m := range module {
		// register model
		model = append(model, m.getModel()...)

		// check and register module
		configKey := fmt.Sprintf("module.%s", m.ModuleName)
		if config.AppConfig.Get(configKey) != nil && !config.AppConfig.GetBool(configKey) {
			disabledModule = append(disabledModule, m)
			continue
		}
		enabledModule = append(enabledModule, m)
		r.modules[m.ModuleName] = m

		// register permission
		rbac.RegisterPerm(m.ModuleName, m.ModulePerm)

		// read and update config
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

	// sync database model
	database.SyncModel(model...)

	// load module
	for _, m := range enabledModule {
		// check dependencies
		for _, dep := range m.ModuleDepends {
			if r.modules[dep] == nil {
				logger.Logger.Fatalf("Module %s dependencies missing: %s", m.ModuleName, dep)
			}
		}

		// init module context
		mctx := &ModuleContext{
			Server:  r.server,
			Route:   router.APIRoute.Party(m.ModuleRoute),
			Storage: storage.InitStorage(m.ModuleConfig),
			Cache:   cache.InitCache(m.ModuleName, m.ModuleConfig, m.getOnEvict()),
		}

		// start loading
		logger.Logger.Debugf("Module Loading: %s", m.ModuleName)

		// load module
		m.EntryPoint(mctx)

		// finish loading
		logger.Logger.Debugf("Module Loaded: %s", m.ModuleName)
	}

	// module log
	em := util.TransSlice(enabledModule, func(m *Module) string { return m.ModuleName })
	dm := util.TransSlice(disabledModule, func(m *Module) string { return m.ModuleName })
	logger.Logger.Infof("%d modules loaded: %v", len(em), em)
	logger.Logger.Infof("%d modules disabled: %v", len(dm), dm)
}

func (r *Registry) Get(moduleName string) IModule {
	return r.modules[moduleName]
}
