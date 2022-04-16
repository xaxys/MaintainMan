package role

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/rbac"
)

func GetPermissionService(name string, auth *model.AuthInfo) *model.ApiJson {
	perm := rbac.GetPermission(name)
	return model.Success(perm, "获取成功")
}

func GetAllPermissionsService(auth *model.AuthInfo) *model.ApiJson {
	perm := rbac.GetAllPermissions()
	return model.Success(perm, "获取成功")
}
