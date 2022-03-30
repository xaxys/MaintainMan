package model

import (
	"maintainman/util"
	"sync"
)

type RoleInfo struct {
	Name           string   `mapstructure:"name"`
	DisplayName    string   `mapstructure:"display_name"`
	Default        bool     `mapstructure:"default"`
	Guest          bool     `mapstructure:"guest"`
	RawPermissions []string `mapstructure:"permissions"`
	RawInheritance []string `mapstructure:"inheritance"`
}

type Role struct {
	*RoleInfo
	Permissions *util.PermSet
	Inheritance []*Role
	sync.RWMutex
}

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required,gte=2,lte=50"`
	DisplayName string   `json:"display_name" validate:"required,lte=191"`
	Permissions []string `json:"permissions"`
	Inheritance []string `json:"inheritance"`
}

type UpdateRoleRequest struct {
	DisplayName    string   `json:"display_name" validate:"required,lte=191"`
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
