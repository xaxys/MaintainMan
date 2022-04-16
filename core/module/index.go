package module

import (
	"github.com/xaxys/maintainman/core/cache"
	"github.com/xaxys/maintainman/core/storage"

	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

type ModuleContext struct {
	*Server
	Route   iris.Party
	Storage storage.IStorage
	Cache   cache.ICache
}

type IModule interface {
	Name() string
	Version() string
	Export(string) (any, bool)
}

type Module struct {
	ModuleName    string
	ModuleVersion string
	ModuleConfig  *viper.Viper
	ModuleEnv     map[string]any // unexported functions or variables, only accessible to system
	ModuleExport  map[string]any // exported functions or variables, accessible to all modules
	ModuleRoute   string         // route prefix
	ModulePerm    map[string]string
	EntryPoint    func(mctx *ModuleContext)
}

func (m *Module) Name() string {
	return m.ModuleName
}

func (m *Module) Version() string {
	return m.ModuleVersion
}

func (m *Module) Export(name string) (any, bool) {
	v, ok := m.ModuleExport[name]
	return v, ok
}

func (m *Module) getOnEvict() func(any) error {
	evictValue, ok := m.ModuleEnv["cache.evict"]
	if !ok {
		return nil
	}
	evictFunc, ok := evictValue.(func(any) error)
	if !ok {
		return nil
	}
	return evictFunc
}

func (m *Module) getModel() []any {
	modelValue, ok := m.ModuleEnv["orm.model"]
	if !ok {
		return nil
	}
	models, ok := modelValue.([]any)
	if !ok {
		return []any{modelValue}
	}
	return models
}
