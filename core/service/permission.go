package service

import (
	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"
)

func GetPermission(name string, auth *model.AuthInfo) *model.ApiJson {
	perm := dao.GetPermission(name)
	return model.Success(perm, "获取成功")
}

func GetAllPermissions(auth *model.AuthInfo) *model.ApiJson {
	perm := dao.GetAllPermissions()
	return model.Success(perm, "获取成功")
}
