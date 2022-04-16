package rbac

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xaxys/maintainman/core/util"

	"github.com/spf13/viper"
)

var (
	RolePO *RolePersistence
)

type RolePersistence struct {
	sync.RWMutex
	data  *viper.Viper
	roles []RoleInfo
	index util.CoPtrMap[string, Role]
	def   util.AtomPtr[Role] // Default role
	guest util.AtomPtr[Role] // Guest role
}

func LoadRole(config *viper.Viper) {
	s := &RolePersistence{
		data: config,
	}

	config.UnmarshalKey("role", &s.roles)
	for i := range s.roles {
		role := &Role{
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
		role.PermSet = newPermSet().Add(role.Permissions...)
		s.index.Set(role.Name, role)
		for _, inhe := range role.Inheritance {
			if !s.index.Has(inhe) {
				panic(fmt.Sprintf("Role %s can not be inherited by %s. Only buttom-up inheritance is valid (latter roles are superior)", inhe, role.Name))
			}
			role.InheRole = append(role.InheRole, s.index.Get(inhe))
		}
	}
	if s.def.Get() == nil {
		panic("Default role is not set")
	}

	RolePO = s
}

func (s *RolePersistence) saveRole() {
	s.Lock()
	s.data.Set("role", s.roles)
	s.data.WriteConfig()
	s.Unlock()
}

func (s *RolePersistence) getRole(role string) (*Role, error) {
	r := s.index.Get(role)
	if r == nil {
		return nil, fmt.Errorf("Role %s does not exist", role)
	}
	return r, nil
}

func AddPermission(role string, permission ...string) error {
	return RolePO.AddPermission(role, permission...)
}

func (s *RolePersistence) AddPermission(role string, perms ...string) error {
	r, err := s.getRole(role)
	if err != nil {
		return err
	}
	addPermission(r, perms...)
	s.saveRole()
	return nil
}

func addPermission(role *Role, perms ...string) {
	role.Lock()
	role.Permissions = append(role.Permissions, perms...)
	role.PermSet.Add(perms...)
	role.Unlock()
}

func DeletePermission(role string, permission ...string) error {
	return RolePO.DeletePermission(role, permission...)
}

func (s *RolePersistence) DeletePermission(role string, permission ...string) error {
	r, err := s.getRole(role)
	if err != nil {
		return err
	}
	deletePermission(r, permission...)
	s.saveRole()
	return nil
}

func deletePermission(role *Role, permission ...string) {
	role.Lock()
	role.Permissions = util.Remove(role.Permissions, permission...)
	role.PermSet.Delete(permission...)
	role.Unlock()
}

func HasPermission(role, permission string) bool {
	return RolePO.HasPermission(role, permission)
}

func (s *RolePersistence) HasPermission(role, permission string) bool {
	r, err := s.getRole(role)
	if err != nil {
		return false
	}
	return hasPermission(r, permission)
}

func hasPermission(role *Role, permission string) bool {
	role.RLock()
	defer role.RUnlock()
	if has, ok := role.PermSet.Find(permission); ok {
		return has
	}
	for _, v := range role.InheRole {
		if has, ok := v.PermSet.Find(permission); ok {
			return has
		}
	}
	return false
}

func GuestHasPermission(permission string) bool {
	return RolePO.GuestHasPermission(permission)
}

func (s *RolePersistence) GuestHasPermission(permission string) bool {
	r := s.getGuestRole()
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
	return RolePO.AddInheritance(role, inherit...)
}

func (s *RolePersistence) AddInheritance(role string, inherit ...string) error {
	r, err := s.getRole(role)
	if err != nil {
		return err
	}
	s.addInheritance(r, inherit...)
	s.saveRole()
	return nil
}

func (s *RolePersistence) addInheritance(role *Role, inherit ...string) error {
	role.Lock()
	nonexist := []string{}
	role.Inheritance = append(role.Inheritance, inherit...)
	for _, inhe := range inherit {
		inheRole := s.index.Get(inhe)
		if inheRole == nil {
			nonexist = append(nonexist, inhe)
		} else {
			role.InheRole = append(role.InheRole, inheRole)
		}
	}
	role.Unlock()

	if len(nonexist) != 0 {
		return fmt.Errorf("Role %s does not exist", strings.Join(nonexist, " "))
	}
	return nil
}

func DeleteInheritance(role string, inherit ...string) error {
	return RolePO.DeleteInheritance(role, inherit...)
}

func (s *RolePersistence) DeleteInheritance(role string, inherit ...string) error {
	r, err := s.getRole(role)
	if err != nil {
		return err
	}
	s.deleteInheritance(r, inherit...)
	s.saveRole()
	return nil
}

func (s *RolePersistence) deleteInheritance(role *Role, inherit ...string) error {
	role.Lock()
	nonexist := []string{}
	role.Inheritance = util.Remove(role.Inheritance, inherit...)
	for _, inhe := range inherit {
		inheRole := s.index.Get(inhe)
		if inheRole == nil {
			nonexist = append(nonexist, inhe)
		} else {
			role.InheRole = util.Remove(role.InheRole, inheRole)
		}
	}
	role.Unlock()

	if len(nonexist) != 0 {
		return fmt.Errorf("Role %s does not exist", strings.Join(nonexist, " "))
	}
	return nil
}

func SetDefaultRole(name string) error {
	return RolePO.SetDefaultRole(name)
}

func (s *RolePersistence) SetDefaultRole(name string) error {
	def := s.getDefaultRole()
	def.RLock()
	defName := def.Name
	def.RUnlock()
	if defName == name {
		return nil
	}

	r, err := s.getRole(name)
	if err != nil {
		return err
	}

	def.Lock()
	def.Default = false
	def.Unlock()

	s.def.Set(r)

	r.Lock()
	r.Default = true
	r.Unlock()

	s.saveRole()
	return nil
}

func SetGuestRole(name string) error {
	return RolePO.SetGuestRole(name)
}

func (s *RolePersistence) SetGuestRole(name string) (err error) {
	guest := s.getGuestRole()
	if guest != nil {
		guest.RLock()
		guestName := guest.Name
		guest.RUnlock()
		if guestName == name {
			return nil
		}
	}

	var r *Role
	if name != "" {
		r, err = s.getRole(name)
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

	s.guest.Set(r)

	if r == nil {
		r.Lock()
		r.Guest = true
		r.Unlock()
	}

	s.saveRole()
	return nil
}

func CreateRole(aul *CreateRoleRequest) error {
	return RolePO.CreateRole(aul)
}

func (s *RolePersistence) CreateRole(aul *CreateRoleRequest) (err error) {
	if s.index.Has(aul.Name) {
		return fmt.Errorf("Role %s already exists", aul.Name)
	}

	info := RoleInfo{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
		Default:     false,
	}

	role := &Role{
		RoleInfo: &info,
		PermSet:  newPermSet(),
	}
	addPermission(role, aul.Permissions...)
	err = s.addInheritance(role, aul.Inheritance...)

	s.Lock()
	// TODO: Allow insert position be specified
	s.roles = append(s.roles, info)
	s.Unlock()

	s.index.Set(aul.Name, role)

	s.saveRole()
	return
}

func UpdateRole(name string, aul *UpdateRoleRequest) error {
	return RolePO.UpdateRole(name, aul)
}

func (s *RolePersistence) UpdateRole(name string, aul *UpdateRoleRequest) error {
	r, err := s.getRole(name)
	if err != nil {
		return fmt.Errorf("Role %s does not exist", name)
	}

	// TODO: Allow role position be adjusted
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
		s.addInheritance(r, aul.AddInheritance...)
	}
	if len(aul.DelInheritance) != 0 {
		s.deleteInheritance(r, aul.DelInheritance...)
	}

	s.saveRole()
	return nil
}

func DeleteRole(name string) error {
	return RolePO.DeleteRole(name)
}

func (s *RolePersistence) DeleteRole(name string) error {
	if !s.index.Has(name) {
		return fmt.Errorf("Role %s does not exist", name)
	}
	if s.def.Get() == s.index.Get(name) {
		return fmt.Errorf("Cannot delete default role")
	}
	err := s.index.Range(func(k string, role *Role) error {
		role.RLock()
		for _, inhe := range role.InheRole {
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
	if r := s.index.LoadAndDelete(name); r != nil {
		s.Lock()
		s.roles = util.RemoveByRef(s.roles, r.RoleInfo)
		s.Unlock()
	}

	s.saveRole()
	return nil
}

func GetRole(name string) *RoleJson {
	return RolePO.GetRole(name)
}

func (s *RolePersistence) GetRole(name string) *RoleJson {
	r, err := s.getRole(name)
	if err != nil {
		return nil
	}
	return roleToJson(r)
}

func (s *RolePersistence) getDefaultRole() *Role {
	return s.def.Get()
}

func GetDefaultRole() *RoleJson {
	return RolePO.GetDefaultRole()
}

func (s *RolePersistence) GetDefaultRole() *RoleJson {
	r := s.getDefaultRole()
	return roleToJson(r)
}

func GetDefaultRoleName() string {
	r := RolePO.getDefaultRole()
	if r == nil {
		return ""
	}
	r.RLock()
	defer r.RUnlock()
	return r.Name
}

func (s *RolePersistence) getGuestRole() *Role {
	return s.guest.Get()
}

func GetGuestRole() *RoleJson {
	return RolePO.GetGuestRole()
}

func (s *RolePersistence) GetGuestRole() *RoleJson {
	if r := s.getGuestRole(); r != nil {
		return roleToJson(r)
	}
	return nil
}

func GetGuestRoleName() string {
	return RolePO.GetGuestRoleName()
}

func (s *RolePersistence) GetGuestRoleName() string {
	r := s.getGuestRole()
	if r == nil {
		return ""
	}
	r.RLock()
	defer r.RUnlock()
	return r.Name
}

func GetAllRoles() []*RoleJson {
	return RolePO.GetAllRoles()
}

func (s *RolePersistence) GetAllRoles() (roles []*RoleJson) {
	s.index.Range(func(k string, r *Role) error {
		role := roleToJson(r)
		roles = append(roles, role)
		return nil
	})
	return
}

func roleToJson(role *Role) *RoleJson {
	if role == nil {
		return nil
	}
	role.RLock()
	defer role.RUnlock()
	return &RoleJson{
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Default:     role.Default,
		Guest:       role.Guest,
		Inheritance: role.Inheritance,
		Permissions: util.TransSlice(role.Permissions, GetPermission),
	}
}
