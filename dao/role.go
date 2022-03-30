package dao

import (
	"fmt"
	"maintainman/config"
	"maintainman/model"
	"maintainman/util"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	RolePO = NewRolePersistence(config.RoleConfig)
)

type RolePersistence struct {
	sync.RWMutex
	data  *viper.Viper
	roles []model.RoleInfo
	index util.CoPtrMap[string, model.Role]
	def   util.AtomPtr[model.Role] // Default role
	guest util.AtomPtr[model.Role] // Guest role
}

func NewRolePersistence(config *viper.Viper) (s *RolePersistence) {
	s = &RolePersistence{
		data: config,
	}

	config.UnmarshalKey("role", &s.roles)
	for i := range s.roles {
		role := &model.Role{
			RoleInfo: &s.roles[i],
		}
		if role.Default {
			if s.def.Get() != nil {
				panic("Default role can only be set once")
			}
			s.def.Set(role)
		}
		if role.Guest {
			if s.guest.Get() != nil {
				panic("Guest role can only be set once")
			}
			s.guest.Set(role)
		}
		role.Permissions = util.NewPermSet().Add(role.RawPermissions...)
		s.index.Set(role.Name, role)
		for _, inhe := range role.RawInheritance {
			if !s.index.Has(inhe) {
				panic(fmt.Sprintf("Role %s can not be inherited by %s. Only buttom-up inheritance is valid (latter roles are superior)", inhe, role.Name))
			}
			role.Inheritance = append(role.Inheritance, s.index.Get(inhe))
		}
	}
	if s.def.Get() == nil {
		panic("Default role is not set")
	}
	return
}

func saveRole() {
	RolePO.Lock()
	RolePO.data.Set("role", RolePO.roles)
	RolePO.data.WriteConfig()
	RolePO.Unlock()
}

func getRole(role string) (*model.Role, error) {
	r := RolePO.index.Get(role)
	if r == nil {
		return nil, fmt.Errorf("Role %s does not exist", role)
	}
	return r, nil
}

func AddPermission(role string, perms ...string) error {
	r, err := getRole(role)
	if err != nil {
		return err
	}
	addPermission(r, perms...)
	saveRole()
	return nil
}

func addPermission(role *model.Role, perms ...string) {
	role.Lock()
	role.RawPermissions = append(role.RawPermissions, perms...)
	role.Permissions.Add(perms...)
	role.Unlock()
}

func DeletePermission(role string, permission ...string) error {
	r, err := getRole(role)
	if err != nil {
		return err
	}
	deletePermission(r, permission...)
	saveRole()
	return nil
}

func deletePermission(role *model.Role, permission ...string) {
	role.Lock()
	role.RawPermissions = util.Remove(role.RawPermissions, permission...)
	role.Permissions.Delete(permission...)
	role.Unlock()
}

func HasPermission(role, permission string) bool {
	r, err := getRole(role)
	if err != nil {
		return false
	}
	return hasPermission(r, permission)
}

func hasPermission(role *model.Role, permission string) bool {
	role.RLock()
	defer role.RUnlock()
	if has, ok := role.Permissions.Find(permission); ok {
		return has
	}
	for _, v := range role.Inheritance {
		if has, ok := v.Permissions.Find(permission); ok {
			return has
		}
	}
	return false
}

func GuestHasPermission(permission string) bool {
	r := getGuestRole()
	if r == nil {
		return false
	}
	return hasPermission(r, permission)
}

func CheckPermission(role, perm string) error {
	if role == "" {
		if !GuestHasPermission(perm) {
			return fmt.Errorf("权限不足：%s", GetPermissionName(perm))
		}
	} else {
		if !HasPermission(role, perm) {
			return fmt.Errorf("权限不足：%s", GetPermissionName(perm))
		}
	}
	return nil
}

func AddInheritance(role string, inherit ...string) error {
	r, err := getRole(role)
	if err != nil {
		return err
	}
	addInheritance(r, inherit...)
	saveRole()
	return nil
}

func addInheritance(role *model.Role, inherit ...string) error {
	role.Lock()
	nonexist := []string{}
	role.RawInheritance = append(role.RawInheritance, inherit...)
	for _, inhe := range inherit {
		inheRole := RolePO.index.Get(inhe)
		if inheRole == nil {
			nonexist = append(nonexist, inhe)
		} else {
			role.Inheritance = append(role.Inheritance, inheRole)
		}
	}
	role.Unlock()

	if len(nonexist) != 0 {
		return fmt.Errorf("Role %s does not exist", strings.Join(nonexist, " "))
	}
	return nil
}

func DeleteInheritance(role string, inherit ...string) error {
	r, err := getRole(role)
	if err != nil {
		return err
	}
	deleteInheritance(r, inherit...)
	saveRole()
	return nil
}

func deleteInheritance(role *model.Role, inherit ...string) error {
	role.Lock()
	nonexist := []string{}
	role.RawInheritance = util.Remove(role.RawInheritance, inherit...)
	for _, inhe := range inherit {
		inheRole := RolePO.index.Get(inhe)
		if inheRole == nil {
			nonexist = append(nonexist, inhe)
		} else {
			role.Inheritance = util.Remove(role.Inheritance, inheRole)
		}
	}
	role.Unlock()

	if len(nonexist) != 0 {
		return fmt.Errorf("Role %s does not exist", strings.Join(nonexist, " "))
	}
	return nil
}

func SetDefaultRole(name string) error {
	def := getDefaultRole()
	def.RLock()
	defName := def.Name
	def.RUnlock()
	if defName == name {
		return nil
	}

	r, err := getRole(name)
	if err != nil {
		return err
	}

	def.Lock()
	def.Default = false
	def.Unlock()

	RolePO.def.Set(r)

	r.Lock()
	r.Default = true
	r.Unlock()

	saveRole()
	return nil
}

func SetGuestRole(name string) (err error) {
	guest := getGuestRole()
	if guest != nil {
		guest.RLock()
		guestName := guest.Name
		guest.RUnlock()
		if guestName == name {
			return nil
		}
	}

	var r *model.Role
	if name != "" {
		r, err = getRole(name)
		if err != nil {
			return err
		}
	}

	if r == nil && guest == nil {
		return nil
	}

	if guest != nil {
		guest.Lock()
		guest.Guest = false
		guest.Unlock()
	}

	RolePO.guest.Set(r)

	if r == nil {
		r.Lock()
		r.Guest = true
		r.Unlock()
	}

	saveRole()
	return nil
}

func CreateRole(aul *model.CreateRoleRequest) (err error) {
	if RolePO.index.Has(aul.Name) {
		return fmt.Errorf("Role %s already exists", aul.Name)
	}

	info := model.RoleInfo{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
		Default:     false,
	}

	role := &model.Role{
		RoleInfo:    &info,
		Permissions: util.NewPermSet(),
	}
	addPermission(role, aul.Permissions...)
	err = addInheritance(role, aul.Inheritance...)

	RolePO.Lock()
	RolePO.roles = append(RolePO.roles, info)
	RolePO.Unlock()

	RolePO.index.Set(aul.Name, role)

	saveRole()
	return
}

func UpdateRole(name string, aul *model.UpdateRoleRequest) error {
	r, err := getRole(name)
	if err != nil {
		return fmt.Errorf("Role %s does not exist", name)
	}

	if aul.DisplayName != "" {
		r.Lock()
		r.DisplayName = aul.DisplayName
		r.Unlock()
	}
	if len(aul.AddPermissions) != 0 {
		addPermission(r, aul.AddPermissions...)
	}
	if len(aul.DelPermissions) != 0 {
		deletePermission(r, aul.DelPermissions...)
	}
	if len(aul.AddInheritance) != 0 {
		addInheritance(r, aul.AddInheritance...)
	}
	if len(aul.DelInheritance) != 0 {
		deleteInheritance(r, aul.DelInheritance...)
	}

	saveRole()
	return nil
}

func DeleteRole(name string) error {
	if !RolePO.index.Has(name) {
		return fmt.Errorf("Role %s does not exist", name)
	}
	if RolePO.def.Get() == RolePO.index.Get(name) {
		RolePO.RUnlock()
		return fmt.Errorf("Cannot delete default role")
	}
	err := RolePO.index.Range(func(k string, role *model.Role) error {
		role.RLock()
		for _, inhe := range role.Inheritance {
			if inhe.Name == name {
				return fmt.Errorf("Cannot delete role %s, it is inherited by %s", name, role.Name)
			}
		}
		role.RUnlock()
		return nil
	})
	if err != nil {
		return err
	}
	if r := RolePO.index.LoadAndDelete(name); r != nil {
		RolePO.Lock()
		RolePO.roles = util.RemoveByRef(RolePO.roles, r.RoleInfo)
		RolePO.Unlock()
	}

	saveRole()
	return nil
}

func GetRole(name string) *model.RoleJson {
	r, err := getRole(name)
	if err != nil {
		return nil
	}
	return RoleToJson(r)
}

func getDefaultRole() *model.Role {
	return RolePO.def.Get()
}

func GetDefaultRole() *model.RoleJson {
	r := getDefaultRole()
	return RoleToJson(r)
}

func GetDefaultRoleName() string {
	r := getDefaultRole()
	if r == nil {
		return ""
	}
	r.RLock()
	defer r.RUnlock()
	return r.Name
}

func getGuestRole() *model.Role {
	return RolePO.guest.Get()
}

func GetGuestRole() *model.RoleJson {
	if r := getGuestRole(); r != nil {
		return RoleToJson(r)
	}
	return nil
}

func GetGuestRoleName() string {
	r := getGuestRole()
	if r == nil {
		return ""
	}
	r.RLock()
	defer r.RUnlock()
	return r.Name
}

func GetAllRoles() (roles []*model.RoleJson) {
	RolePO.index.Range(func(k string, r *model.Role) error {
		role := RoleToJson(r)
		roles = append(roles, role)
		return nil
	})
	return
}

func RoleToJson(role *model.Role) *model.RoleJson {
	if role == nil {
		return nil
	}
	role.RLock()
	defer role.RUnlock()
	return &model.RoleJson{
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Default:     role.Default,
		Guest:       role.Guest,
		Inheritance: role.RawInheritance,
		Permissions: util.TransSlice(role.RawPermissions, GetPermission),
	}
}
