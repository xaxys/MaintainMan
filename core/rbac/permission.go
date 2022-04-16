package rbac

var (
	PermPO = &PermissionPersistence{
		perm: make(map[string]string),
	}
)

type PermissionPersistence struct {
	perm map[string]string
}

func RegisterPerm(name string, perm map[string]string) {
	for k, v := range perm {
		PermPO.perm[k] = v
	}
}

// GetPermissionName 获取权限名称
func GetPermissionName(name string) string {
	if v, ok := PermPO.perm[name]; ok {
		return v
	}
	return name
}

// GetPermission 获取权限Json
func GetPermission(name string) *PermissionJson {
	return &PermissionJson{
		Name:        name,
		DisplayName: GetPermissionName(name),
	}
}

func GetAllPermissions() (perms []*PermissionJson) {
	for k, v := range PermPO.perm {
		perms = append(perms, &PermissionJson{
			Name:        k,
			DisplayName: v,
		})
	}
	return
}
