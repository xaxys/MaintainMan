package dao

import (
	"maintainman/config"
	"maintainman/model"

	"github.com/spf13/viper"
)

var (
	PermPO = NewPermissionPersistence(config.PermConfig)
)

type PermissionPersistence struct {
	data map[string]string
}

func NewPermissionPersistence(config *viper.Viper) (s *PermissionPersistence) {
	s = &PermissionPersistence{
		data: make(map[string]string),
	}

	var getPermission func(string)
	getPermission = func(prefix string) {
		for k, v := range config.GetStringMap(prefix) {
			if name, ok := v.(string); ok {
				s.data[prefix+"."+k] = name
			} else if _, ok := v.(map[string]interface{}); ok {
				getPermission(prefix + "." + k)
			}
		}
	}
	getPermission("permission")
	return
}

// GetPermissionName 获取权限名称
func GetPermissionName(name string) string {
	return PermPO.data[name]
}

// GetPermission 获取权限Json
func GetPermission(name string) *model.PermissionJson {
	return &model.PermissionJson{
		Name:        name,
		DisplayName: GetPermissionName(name),
	}
}

func GetAllPermissions() (perms []*model.PermissionJson) {
	for k, v := range PermPO.data {
		perms = append(perms, &model.PermissionJson{
			Name:        k,
			DisplayName: v,
		})
	}
	return
}
