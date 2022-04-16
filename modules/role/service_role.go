package role

import (
	"fmt"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/rbac"
	"github.com/xaxys/maintainman/core/util"
)

func getRoleByNameService(name string, auth *model.AuthInfo) *model.ApiJson {
	role := rbac.GetRole(name)
	return model.Success(role, "获取成功")
}

func createRoleService(aul *rbac.CreateRoleRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if rbac.GetRole(aul.Name) != nil {
		return model.ErrorInsertDatabase(fmt.Errorf("Role %s already exists", aul.Name))
	}
	if aul.DisplayName != "" {
		aul.DisplayName = aul.Name
	}
	err := rbac.CreateRole(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	role := rbac.GetRole(aul.Name)
	return model.SuccessCreate(role, "创建成功")

}

func updateRoleService(name string, aul *rbac.UpdateRoleRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if rbac.GetRole(name) != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("Role %s already exists", name))
	}

	err := rbac.UpdateRole(name, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	role := rbac.GetRole(name)
	return model.SuccessUpdate(role, "更新成功")
}

func deleteRoleService(name string, auth *model.AuthInfo) *model.ApiJson {
	if rbac.GetRole(name) == nil {
		return model.ErrorNotFound(fmt.Errorf("Role %s not found", name))
	}
	err := rbac.DeleteRole(name)
	if err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func setDefaultRoleService(name string, auth *model.AuthInfo) *model.ApiJson {
	if rbac.GetRole(name) == nil {
		return model.ErrorNotFound(fmt.Errorf("Role %s not found", name))
	}
	err := rbac.SetDefaultRole(name)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "操作成功")
}

func setGuestRoleService(name string, auth *model.AuthInfo) *model.ApiJson {
	if rbac.GetRole(name) == nil {
		return model.ErrorNotFound(fmt.Errorf("Role %s not found", name))
	}
	err := rbac.SetGuestRole(name)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "操作成功")
}

func getAllRolesService(auth *model.AuthInfo) *model.ApiJson {
	roles := rbac.GetAllRoles()
	return model.Success(roles, "操作成功")
}
