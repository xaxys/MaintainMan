package model

import (
	"encoding/json"
	"maintainman/util"
	"sync"
)

type RoleInfo struct {
	Name        string   `mapstructure:"name"         yaml:"name"`
	DisplayName string   `mapstructure:"display_name" yaml:"display_name"`
	Default     bool     `mapstructure:"default"      yaml:"default,omitempty"`
	Guest       bool     `mapstructure:"guest"        yaml:"guest,omitempty"`
	Permissions []string `mapstructure:"permissions"  yaml:"permissions"`
	Inheritance []string `mapstructure:"inheritance"  yaml:"inheritance"`
}

type Role struct {
	*RoleInfo
	PermSet  *util.PermSet
	InheRole []*Role
	sync.RWMutex
}

func (r *Role) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

type CreateRoleRequest struct {
	Name        string   `json:"name"         validate:"required,gte=2,lte=50"`
	DisplayName string   `json:"display_name" validate:"required,lte=191"`
	Position    uint     `json:"position"`
	Permissions []string `json:"permissions"`
	Inheritance []string `json:"inheritance"`
}

type UpdateRoleRequest struct {
	DisplayName    string   `json:"display_name" validate:"required,lte=191"`
	Position       uint     `json:"position"`
	AddPermissions []string `json:"add_permissions"`
	DelPermissions []string `json:"del_permissions"`
	AddInheritance []string `json:"add_inheritance"`
	DelInheritance []string `json:"del_inheritance"`
}

type RoleJson struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Default     bool              `json:"default"`
	Guest       bool              `json:"guest"`
	Permissions []*PermissionJson `json:"permissions,omitempty"`
	Inheritance []string          `json:"inheritance,omitempty"`
}
