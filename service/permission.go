package service

import (
	"maintainman/dao"
	"maintainman/model"
)

func GetPermission(name string) *model.ApiJson {
	perm := dao.GetPermission(name)
	return model.Success(perm, "获取成功")
}

func GetAllPermissions() *model.ApiJson {
	perm := dao.GetAllPermissions()
	return model.Success(perm, "获取成功")
}
