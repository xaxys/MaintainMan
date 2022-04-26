package user

import (
	"fmt"
	"time"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/rbac"
	"github.com/xaxys/maintainman/core/util"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func dbGetUserCount() (uint, error) {
	return txGetUserCount(mctx.Database)
}

func txGetUserCount(tx *gorm.DB) (uint, error) {
	count := int64(0)
	if err := tx.Model(&User{}).Count(&count).Error; err != nil {
		mctx.Logger.Warnf("GetUserCountErr: %v\n", err)
		return 0, err
	}
	return uint(count), nil
}

func dbGetUserByID(id uint) (user *User, err error) {
	if user, err = cacheGetUserByID(id); err == nil {
		return
	}
	mctx.Logger.Debugf("CacheGetUserByIDErr: %v", err)
	user, err = txGetUserByID(mctx.Database, id)
	if err != nil {
		return
	}
	if err := cacheSaveUser(user); err != nil {
		mctx.Logger.Warnf("CacheSaveUserErr: %v", err)
	}
	return
}

func txGetUserByID(tx *gorm.DB, id uint) (*User, error) {
	user := &User{}
	if err := tx.First(user, id).Error; err != nil {
		mctx.Logger.Warnf("GetUserByIDErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbGetUserByName(name string) (*User, error) {
	return txGetUserByName(mctx.Database, name)
}

func txGetUserByName(tx *gorm.DB, name string) (*User, error) {
	user := &User{Name: name}
	if err := tx.Where(user).First(user).Error; err != nil {
		mctx.Logger.Warnf("GetUserByNameErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbGetUserByEmail(email string) (*User, error) {
	return txGetUserByEmail(mctx.Database, email)
}

func txGetUserByEmail(tx *gorm.DB, email string) (*User, error) {
	user := &User{Email: email}
	if err := tx.Where(user).First(user).Error; err != nil {
		mctx.Logger.Warnf("GetUserByEmailErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbGetUserByPhone(phone string) (*User, error) {
	return txGetUserByPhone(mctx.Database, phone)
}

func txGetUserByPhone(tx *gorm.DB, phone string) (*User, error) {
	user := &User{Phone: phone}
	if err := tx.Where(user).First(user).Error; err != nil {
		mctx.Logger.Warnf("GetUserByPhoneErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbGetUserByOpenID(openid string) (*User, error) {
	return txGetUserByOpenID(mctx.Database, openid)
}

func txGetUserByOpenID(tx *gorm.DB, openid string) (*User, error) {
	user := &User{OpenID: openid}
	if err := tx.Where(user).First(user).Error; err != nil {
		mctx.Logger.Warnf("GetUserByOpenIDErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbGetUsersByDivision(id uint, param *model.PageParam) (users []*User, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if users, count, err = txGetUserByDivision(tx, id, param); err != nil {
			mctx.Logger.Warnf("GetUsersByDivisionErr: %v\n", err)
		}
		return err
	})
	return
}

func txGetUserByDivision(tx *gorm.DB, id uint, param *model.PageParam) (users []*User, count uint, err error) {
	user := &User{}
	user.DivisionID = &id
	tx = dao.TxPageFilter(tx, param).Where(user)
	if id == 0 {
		tx = tx.Where("division_id is null")
	}
	if err = tx.Find(&users).Error; err != nil {
		return
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}

func dbGetAllUsersWithParam(aul *AllUserRequest) (users []*User, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if users, count, err = txGetAllUsersWithParam(tx, aul); err != nil {
			mctx.Logger.Warnf("GetAllUsersWithParamErr: %v\n", err)
		}
		return err
	})
	return
}

func txGetAllUsersWithParam(tx *gorm.DB, aul *AllUserRequest) (users []*User, count uint, err error) {
	user := &User{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
	}
	tx = dao.TxPageFilter(tx, &aul.PageParam).Where(user)
	if err = tx.Find(&users).Error; err != nil {
		return
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}

func dbCreateUser(json *CreateUserRequest, operator uint) (*User, error) {
	return txCreateUser(mctx.Database, json, operator)
}

func txCreateUser(tx *gorm.DB, json *CreateUserRequest, operator uint) (*User, error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(json.Password, salt)
	json.Password = string(hash)
	json.DisplayName = util.NotEmpty(json.DisplayName, json.Name)

	if json.RoleName == "" {
		json.RoleName = rbac.GetDefaultRoleName()
	} else if rbac.GetRole(json.RoleName) == nil {
		return nil, fmt.Errorf("role %s not found", json.RoleName)
	}

	user := &User{}
	copier.Copy(user, json)
	user.DivisionID = util.Tenary(json.DivisionID != 0, &json.DivisionID, nil)
	user.CreatedBy = operator
	user.LoginTime = time.Now()

	if err := tx.Create(user).Error; err != nil {
		mctx.Logger.Warnf("CreateUserErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbUpdateUser(id uint, json *UpdateUserRequest, operator uint) (user *User, err error) {
	user, err = txUpdateUser(mctx.Database, id, json, operator)
	if err != nil {
		return
	}
	cacheDeleteUser(id)
	return
}

func txUpdateUser(tx *gorm.DB, id uint, json *UpdateUserRequest, operator uint) (*User, error) {
	if json.Password != "" {
		salt, _ := bcrypt.Salt(10)
		hash, _ := bcrypt.Hash(json.Password, salt)
		json.Password = string(hash)
	}
	if json.RoleName != "" && rbac.GetRole(json.RoleName) == nil {
		return nil, fmt.Errorf("role %s not found", json.RoleName)
	}

	user := &User{}
	copier.Copy(user, json)
	user.ID = id
	user.UpdatedBy = operator
	division := uint(json.DivisionID)
	user.DivisionID = util.Tenary(json.DivisionID > 0, &division, nil)
	tx = tx.Model(user).Updates(user)
	if json.DivisionID == -1 {
		tx = tx.Update("division_id", nil)
	}
	if err := tx.Error; err != nil {
		mctx.Logger.Warnf("UpdateUserErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func dbAttachOpenIDToUser(id uint, openid string) error {
	err := txAttachOpenIDToUser(mctx.Database, id, openid)
	if err != nil {
		return err
	}
	cacheDeleteUser(id)
	return nil
}

func txAttachOpenIDToUser(tx *gorm.DB, id uint, openid string) error {
	user := &User{}
	user.ID = id
	if err := tx.Model(user).Update("open_id", openid).Error; err != nil {
		mctx.Logger.Warnf("AttachOpenIDToUserErr: %v\n", err)
		return err
	}
	return nil
}

func dbDeleteUser(id uint) error {
	err := txDeleteUser(mctx.Database, id)
	if err != nil {
		return err
	}
	cacheDeleteUser(id)
	return nil
}

func txDeleteUser(tx *gorm.DB, id uint) (err error) {
	if err = tx.Delete(&User{}, id).Error; err != nil {
		mctx.Logger.Warnf("DeleteUserByIdErr: %v\n", err)
	}
	return
}

func dbCheckLogin(user *User, password string) error {
	if ok := bcrypt.Match(password, user.Password); !ok {
		return fmt.Errorf("Wrong password")
	}
	return dbForceLogin(user.ID, user.LoginIP)
}

func dbForceLogin(id uint, ip string) error {
	err := txForceLogin(mctx.Database, id, ip)
	if err != nil {
		return err
	}
	cacheDeleteUser(id)
	return nil
}

func txForceLogin(tx *gorm.DB, id uint, ip string) error {
	user := &User{
		LoginIP:   ip,
		LoginTime: time.Now(),
	}
	user.ID = id
	if err := tx.Model(user).Updates(user).Error; err != nil {
		mctx.Logger.Warnf("ForceLoginErr: %v\n", err)
		return err
	}
	return nil
}
