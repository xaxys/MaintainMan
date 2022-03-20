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
	RolePO   = NewRolePersistence(config.RoleConfig)
	roleLock = sync.RWMutex{}
)

type RoleWithLock struct {
	*model.Role
	sync.RWMutex
}

type RolePersistence struct {
	data  *viper.Viper
	roles []model.RoleInfo
	index map[string]*RoleWithLock
	def   *RoleWithLock // Default role
	guest *RoleWithLock // Guest role
}

func NewRolePersistence(config *viper.Viper) (s *RolePersistence) {
	s = &RolePersistence{
		data:  config,
		index: make(map[string]*RoleWithLock),
	}

	config.UnmarshalKey("role", &s.roles)
	for i := range s.roles {
		role := &RoleWithLock{
			Role: &model.Role{
				RoleInfo: &s.roles[i],
			},
		}
		if role.Default {
			if s.def != nil {
				panic("Default role can only be set once")
			}
			s.def = role
		}
		if role.Guest {
			if s.guest != nil {
				panic("Guest role can only be set once")
			}
			s.guest = role
		}
		role.Permissions = model.NewPermissionSet().Add(role.RawPermissions...)
		s.index[role.Name] = role
		for _, inhe := range role.RawInheritance {
			inheRole := s.index[inhe]
			if inheRole == nil {
				panic(fmt.Sprintf("Role %s can not be inherited by %s. Only buttom-up inheritance is valid (latter roles are superior)", inhe, role.Name))
			}
			role.Inheritance = append(role.Inheritance, inheRole.Role)
		}
	}
	if s.def == nil {
		panic("Default role is not set")
	}
	return
}

func saveRole() {
	roleLock.Lock()
	RolePO.data.Set("role", RolePO.roles)
	go RolePO.data.WriteConfig()
	roleLock.Unlock()
}

func getRole(role string) (*RoleWithLock, error) {
	roleLock.RLock()
	defer roleLock.RUnlock()
	r := RolePO.index[role]
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

func addPermission(role *RoleWithLock, perms ...string) {
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

func deletePermission(role *RoleWithLock, permission ...string) {
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

func hasPermission(role *RoleWithLock, permission string) bool {
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

func addInheritance(role *RoleWithLock, inherit ...string) error {
	role.Lock()
	nonexist := []string{}
	role.RawInheritance = append(role.RawInheritance, inherit...)
	for _, inhe := range inherit {
		inheRole := RolePO.index[inhe]
		if inheRole == nil {
			nonexist = append(nonexist, inhe)
		}
		role.Inheritance = append(role.Inheritance, inheRole.Role)
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

func deleteInheritance(role *RoleWithLock, inherit ...string) error {
	role.Lock()
	nonexist := []string{}
	role.RawInheritance = util.Remove(role.RawInheritance, inherit...)
	for _, inhe := range inherit {
		inheRole := RolePO.index[inhe]
		if inheRole == nil {
			nonexist = append(nonexist, inhe)
		}
		role.Inheritance = util.Remove(role.Inheritance, inheRole.Role)
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

	roleLock.Lock()
	RolePO.def = r
	roleLock.Unlock()

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

	var r *RoleWithLock
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

	roleLock.Lock()
	RolePO.guest = r
	roleLock.Unlock()

	if r == nil {
		r.Lock()
		r.Guest = true
		r.Unlock()
	}

	saveRole()
	return nil
}

func CreateRole(aul *model.CreateRoleJson) (err error) {
	roleLock.RLock()
	if _, ok := RolePO.index[aul.Name]; ok {
		err = fmt.Errorf("Role %s already exists", aul.Name)
	}
	roleLock.RUnlock()
	if err != nil {
		return
	}

	info := model.RoleInfo{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
		Default:     false,
	}

	role := &RoleWithLock{
		Role: &model.Role{
			RoleInfo:    &info,
			Permissions: model.NewPermissionSet(),
		},
	}
	addPermission(role, aul.Permissions...)
	err = addInheritance(role, aul.Inheritance...)

	roleLock.Lock()
	RolePO.roles = append(RolePO.roles, info)
	RolePO.index[aul.Name] = role
	roleLock.Unlock()

	saveRole()
	return
}

func UpdateRole(name string, aul *model.UpdateRoleJson) error {
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
	roleLock.RLock()
	if RolePO.def == RolePO.index[name] {
		roleLock.RUnlock()
		return fmt.Errorf("Cannot delete default role")
	}
	for _, role := range RolePO.index {
		for _, inhe := range role.Inheritance {
			if inhe.Name == name {
				roleLock.RUnlock()
				return fmt.Errorf("Cannot delete role %s, it is inherited by %s", name, role.Name)
			}
		}
	}
	roleLock.RUnlock()

	roleLock.Lock()
	r := RolePO.index[name]
	RolePO.roles = util.RemoveByRef(RolePO.roles, r.RoleInfo)
	delete(RolePO.index, name)
	roleLock.Unlock()

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

func getDefaultRole() *RoleWithLock {
	roleLock.RLock()
	defer roleLock.RUnlock()
	return RolePO.def
}

func GetDefaultRole() *model.RoleJson {
	r := getDefaultRole()
	return RoleToJson(r)
}

func GetDefaultRoleName() string {
	roleLock.RLock()
	defer roleLock.RUnlock()
	return RolePO.def.Name
}

func getGuestRole() *RoleWithLock {
	roleLock.RLock()
	defer roleLock.RUnlock()
	return RolePO.guest
}

func GetGuestRole() *model.RoleJson {
	r := getGuestRole()
	if r == nil {
		return nil
	}
	return RoleToJson(r)
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
	roleLock.RLock()
	for _, r := range RolePO.index {
		role := RoleToJson(r)
		roles = append(roles, role)
	}
	roleLock.RUnlock()
	return
}

func RoleToJson(role *RoleWithLock) *model.RoleJson {
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
