package service

import (
	"maintainman/dao"
	"maintainman/model"
)

func GetPermission(name string, auth *model.AuthInfo) *model.ApiJson {
	perm := dao.GetPermission(name)
	return model.Success(perm, "获取成功")
}

func GetAllPermissions(auth *model.AuthInfo) *model.ApiJson {
	perm := dao.GetAllPermissions()
	return model.Success(perm, "获取成功")
}
