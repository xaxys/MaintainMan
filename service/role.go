package service

import (
	"fmt"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"
)

func GetRoleByName(name string) *model.ApiJson {
	role := dao.GetRole(name)
	return model.Success(role, "获取成功")
}

func CreateRole(aul *model.CreateRoleJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if dao.GetRole(aul.Name) != nil {
		return model.ErrorInsertDatabase(fmt.Errorf("Role %s already exists", aul.Name))
	}
	if aul.DisplayName != "" {
		aul.DisplayName = aul.Name
	}
	err := dao.CreateRole(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	role := dao.GetRole(aul.Name)
	return model.SuccessCreate(role, "创建成功")

}

func UpdateRole(name string, aul *model.UpdateRoleJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if dao.GetRole(name) != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("Role %s already exists", name))
	}

	err := dao.UpdateRole(name, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	role := dao.GetRole(name)
	return model.SuccessUpdate(role, "更新成功")
}

func DeleteRole(name string) *model.ApiJson {
	if dao.GetRole(name) == nil {
		return model.ErrorNotFound(fmt.Errorf("Role %s not found", name))
	}
	err := dao.DeleteRole(name)
	if err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func SetDefaultRole(name string) *model.ApiJson {
	if dao.GetRole(name) == nil {
		return model.ErrorNotFound(fmt.Errorf("Role %s not found", name))
	}
	err := dao.SetDefaultRole(name)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "操作成功")
}

func SetGuestRole(name string) *model.ApiJson {
	if dao.GetRole(name) == nil {
		return model.ErrorNotFound(fmt.Errorf("Role %s not found", name))
	}
	err := dao.SetGuestRole(name)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "操作成功")
}

func GetAllRoles() *model.ApiJson {
	roles := dao.GetAllRoles()
	return model.Success(roles, "操作成功")
}
