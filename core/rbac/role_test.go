package rbac

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

func TestRole(t *testing.T) {
	config := viper.New()
	config.SetDefault("role", []any{
		map[string]any{
			"name":         "user",
			"display_name": "普通用户",
			"default":      true,
			"permissions":  []string{},
			"inheritance":  []string{},
		},
	})
	LoadRole(config)
	// test get all roles
	roles := GetAllRoles()
	fmt.Printf("all %d roles: %v\n", len(roles), roles)

	// test creat role
	aul := CreateRoleRequest{
		Name:        "test role 1",
		DisplayName: "test role 1",
		Permissions: []string{"perm.test"},
	}
	err := CreateRole(&aul)
	if err != nil {
		t.Error(err)
	}

	aul2 := CreateRoleRequest{
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
	aul3 := UpdateRoleRequest{
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

	aul4 := UpdateRoleRequest{
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

func TestRoleConcurrency(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config := viper.New()
	config.SetDefault("role", []any{
		map[string]any{
			"name":         "user",
			"display_name": "普通用户",
			"default":      true,
			"permissions":  []string{},
			"inheritance":  []string{},
		},
	})
	LoadRole(config)

	roles := GetAllRoles()
	fmt.Printf("all %d roles: %v\n", len(roles), roles)
	defRole := GetDefaultRole()
	guestRole := GetGuestRole()

	random2 := func(n int) (a int, b int) {
		b = rand.Intn(n-1) + 1
		a = rand.Intn(b)
		return
	}

	createFunc := func(prefix string, i, n int) error {
		Role1 := fmt.Sprintf("%s role %d", prefix, i)
		Perm1 := fmt.Sprintf("%s perm %d", prefix, rand.Intn(n))
		aul := CreateRoleRequest{
			Name:        Role1,
			DisplayName: Role1,
			Permissions: []string{Perm1},
		}
		fmt.Printf("[INFO] [CREATE] create role %s\n", Role1)
		err := CreateRole(&aul)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] [CREATE] role %s created\n", Role1)
		return nil
	}

	testFunc := func(prefix string, index, n int) error {
		roles := GetAllRoles()
		fmt.Printf("[%d] all %d roles: %v\n", index, len(roles), roles)

		a, b := random2(n)
		Role1 := fmt.Sprintf("%s role %d", prefix, a)
		Role2 := fmt.Sprintf("%s role %d", prefix, b)
		Perm1 := fmt.Sprintf("%s perm %d", prefix, rand.Intn(n))
		Perm2 := fmt.Sprintf("%s perm %d", prefix, rand.Intn(n))

		// test get defRole
		role := GetRole(Role1)
		if role == nil {
			return fmt.Errorf("%s not found", Role1)
		}
		fmt.Printf("[INFO] [TEST] [%d] GetRole(%s): %v\n", index, Role1, *role)

		// test update role
		aul3 := UpdateRoleRequest{
			DisplayName:    fmt.Sprintf("%s (update by %d)", Role1, index),
			AddPermissions: []string{Perm2},
			DelPermissions: []string{Perm1},
		}
		err := UpdateRole(Role1, &aul3)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] [TEST] [%d] UpdateRole(%s): Add:%v Del:%v\n", index, Role1, aul3.AddPermissions, aul3.DelPermissions)

		if rand.Intn(n)%2 == 0 {
			aul4 := UpdateRoleRequest{
				DelInheritance: []string{Role1},
			}
			err = UpdateRole(Role2, &aul4)
			if err != nil {
				return fmt.Errorf("update role %s, del %v failed: %v", Role2, aul3.DelPermissions, err)
			}
			fmt.Printf("[INFO] [TEST] [%d] UpdateRole(%s): Del:%v\n", index, Role2, aul3.DelPermissions)
		} else {
			aul4 := UpdateRoleRequest{
				AddInheritance: []string{Role1},
			}
			err = UpdateRole(Role2, &aul4)
			if err != nil {
				return fmt.Errorf("update role %s, add %v failed: %v", Role2, aul3.AddPermissions, err)
			}
			fmt.Printf("[INFO] [TEST] [%d] UpdateRole(%s): Add:%v\n", index, Role2, aul3.AddPermissions)
		}

		// test get default role
		defRole := GetDefaultRole()
		if defRole == nil {
			return fmt.Errorf("default role not found")
		}
		fmt.Printf("[INFO] [TEST] [%d] GetDefaultRole: %v\n", index, *defRole)

		// test set default role
		err = SetDefaultRole(Role1)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] [TEST] [%d] SetDefaultRole to: %s\n", index, Role1)

		role = GetDefaultRole()
		if role == nil {
			return fmt.Errorf("default role not set")
		}
		fmt.Printf("[INFO] [TEST] [%d] after set, GetDefaultRole: %v\n", index, *role)

		err = SetDefaultRole(defRole.Name)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] [TEST] [%d] set back SetDefaultRole to: %s\n", index, defRole.Name)

		// test get all roles
		roles = GetAllRoles()
		fmt.Printf("[INFO] [TEST] [%d] all %d roles: %v\n", index, len(roles), roles)
		return nil
	}

	deleteFunc := func(prefix string, i int) error {
		Role1 := fmt.Sprintf("%s role %d", prefix, i)
		err := DeleteRole(Role1)
		fmt.Printf("[INFO] [DELETE] delete role %s\n", Role1)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] [DELETE] role %s deleted\n", Role1)
		return nil
	}

	prefix := "cotest"
	n := 20
	wg := sync.WaitGroup{}

	fmt.Println("[INFO] [CREATE] start create")
	for i := 0; i < n; i++ {
		index := i
		wg.Add(1)
		go func() {
			fmt.Printf("[INFO] [CREATE] [%d] start\n", index)
			err := createFunc(prefix, index, n)
			if err != nil {
				t.Error(err)
				fmt.Printf("[ERR] [CREATE] [%d]: %s\n", index, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("[INFO] [CREATE] all create done")

	fmt.Println("[INFO] [TEST] start test")
	for i := 0; i < n*5; i++ {
		index := i
		wg.Add(1)
		go func() {
			fmt.Printf("[INFO] [TEST] [%d] start\n", index)
			err := testFunc(prefix, index, n)
			if err != nil {
				t.Error(err)
				fmt.Printf("[ERR] [TEST] [%d]: %s\n", index, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("[INFO] [TEST] all test done")

	SetDefaultRole(defRole.Name)
	if guestRole != nil {
		SetGuestRole(guestRole.Name)
	}
	fmt.Println("[INFO] [DELETE] start delete")
	for i := n - 1; i >= 0; i-- {
		index := i
		wg.Add(1)
		go func() {
			fmt.Printf("[INFO] [DELETE] [%d] start\n", index)
			err := deleteFunc(prefix, index)
			if err != nil {
				fmt.Printf("[ERR] [DELETE] [%d]: %s\n", index, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("[INFO] [DELETE] all delete done")

	for i := n - 1; i >= 0; i-- {
		index := i
		fmt.Printf("[INFO] [SEQ_DELETE] [%d] start\n", index)
		err := deleteFunc(prefix, index)
		if err != nil {
			fmt.Printf("[ERR] [SEQ_DELETE] [%d]: %s\n", index, err)
		}
	}
	if len(GetAllRoles()) != len(roles) {
		t.Errorf("[ERR] [SEQ_DELETE] not all created roles deleted")
	}
}
