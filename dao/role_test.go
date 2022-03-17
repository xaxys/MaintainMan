package dao

import (
	"fmt"
	"maintainman/model"
	"testing"
)

func TestRole(t *testing.T) {
	// test get all roles
	roles := GetAllRoles()
	fmt.Printf("all %d roles: %v\n", len(roles), roles)

	// test creat role
	aul := *&model.CreateRoleJson{
		Name:        "test role 1",
		DisplayName: "test role 1",
		Permissions: []string{"perm.test"},
	}
	err := CreateRole(&aul)
	if err != nil {
		t.Error(err)
	}

	aul2 := *&model.CreateRoleJson{
		Name:        "test role 2",
		Inheritance: []string{"test role 1"},
	}
	err = CreateRole(&aul2)
	if err != nil {
		t.Error(err)
	}

	// test has permission
	if !HasPermission("test role 1", "perm.test") {
		t.Error("test role 1 does not has test permission")
	}
	if !HasPermission("test role 2", "perm.test") {
		t.Error("test role 2 does not has test permission")
	}

	// test get defRole
	role := GetRole("test role 1")
	if role == nil {
		t.Error("test role 1 not found")
	}
	fmt.Printf("role test role 1: %v\n", *role)

	// test update role
	aul3 := *&model.UpdateRoleJson{
		DisplayName:    "test role 1 (new)",
		AddPermissions: []string{"perm.test.2"},
		DelPermissions: []string{"perm.test"},
	}
	err = UpdateRole("test role 1", &aul3)
	if err != nil {
		t.Error(err)
	}
	if HasPermission("test role 1", "perm.test") {
		t.Error("test role 1 has test permission")
	}
	if !HasPermission("test role 1", "perm.test.2") {
		t.Error("test role 1 does not has test role 2 permission")
	}

	aul4 := *&model.UpdateRoleJson{
		DelInheritance: []string{"test role 1"},
	}
	err = UpdateRole("test role 2", &aul4)
	if err != nil {
		t.Error(err)
	}
	if HasPermission("test role 2", "perm.test.2") {
		t.Error("test role 2 has test permission")
	}

	// test get default role
	defRole := GetDefaultRole()
	if defRole == nil {
		t.Error("default role not found")
	}
	fmt.Printf("role default role: %v\n", *defRole)

	// test set default role
	err = SetDefaultRole("test role 1")
	if err != nil {
		t.Error(err)
	}
	role = GetDefaultRole()
	if role == nil || role.Name != "test role 1" {
		t.Error("default role not set")
	}
	err = SetDefaultRole(defRole.Name)
	if err != nil {
		t.Error(err)
	}

	// test delete role
	err = DeleteRole("test role 1")
	if err != nil {
		t.Error(err)
	}

	err = DeleteRole("test role 2")
	if err != nil {
		t.Error(err)
	}

	// test get all roles
	roles = GetAllRoles()
	fmt.Printf("all %d roles: %v\n", len(roles), roles)
}
