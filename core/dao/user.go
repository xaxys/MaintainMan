package dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xaxys/maintainman/core/database"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetUserCount() (int, error) {
	return TxGetUserCount(database.DB)
}

func TxGetUserCount(tx *gorm.DB) (int, error) {
	count := int64(0)
	if err := tx.Model(&model.User{}).Count(&count).Error; err != nil {
		logger.Logger.Debugf("GetUserCountErr: %v\n", err)
		return 0, err
	}
	return int(count), nil
}

func GetUserByID(id uint) (*model.User, error) {
	return TxGetUserByID(database.DB, id)
}

func TxGetUserByID(tx *gorm.DB, id uint) (*model.User, error) {
	user := &model.User{}
	if err := tx.First(user, id).Error; err != nil {
		logger.Logger.Debugf("GetUserByIDErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetUserByName(name string) (*model.User, error) {
	return TxGetUserByName(database.DB, name)
}

func TxGetUserByName(tx *gorm.DB, name string) (*model.User, error) {
	user := &model.User{Name: name}
	if err := tx.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByNameErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	return TxGetUserByEmail(database.DB, email)
}

func TxGetUserByEmail(tx *gorm.DB, email string) (*model.User, error) {
	user := &model.User{Email: email}
	if err := tx.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByEmailErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetUserByPhone(phone string) (*model.User, error) {
	return TxGetUserByPhone(database.DB, phone)
}

func TxGetUserByPhone(tx *gorm.DB, phone string) (*model.User, error) {
	user := &model.User{Phone: phone}
	if err := tx.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByPhoneErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetUserByOpenID(openid string) (*model.User, error) {
	return TxGetUserByOpenID(database.DB, openid)
}

func TxGetUserByOpenID(tx *gorm.DB, openid string) (*model.User, error) {
	user := &model.User{OpenID: openid}
	if err := tx.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByOpenIDErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetUserByDivision(id uint, param *model.PageParam) ([]*model.User, error) {
	return TxGetUserByDivision(database.DB, id, param)
}

func TxGetUserByDivision(tx *gorm.DB, id uint, param *model.PageParam) (users []*model.User, err error) {
	user := &model.User{}
	user.DivisionID = sql.NullInt64{Int64: int64(id), Valid: id != 0}
	tx = TxPageFilter(tx, param).Where(user)
	if id == 0 {
		tx = tx.Where("division_id is null")
	}
	if err = tx.Find(&users).Error; err != nil {
		logger.Logger.Debugf("GetUserByDivisionErr: %v\n", err)
	}
	return
}

func GetAllUsersWithParam(aul *model.AllUserRequest) ([]*model.User, error) {
	return TxGetAllUsersWithParam(database.DB, aul)
}

func TxGetAllUsersWithParam(tx *gorm.DB, aul *model.AllUserRequest) (users []*model.User, err error) {
	user := &model.User{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
	}
	if err = TxPageFilter(tx, &aul.PageParam).Where(user).Find(&users).Error; err != nil {
		logger.Logger.Debugf("GetAllUserErr: %v\n", err)
	}
	return
}

func CreateUser(json *model.CreateUserRequest, operator uint) (*model.User, error) {
	return TxCreateUser(database.DB, json, operator)
}

func TxCreateUser(tx *gorm.DB, json *model.CreateUserRequest, operator uint) (*model.User, error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(json.Password, salt)
	json.Password = string(hash)
	json.DisplayName = util.NotEmpty(json.DisplayName, json.Name)
	json.RoleName = util.NotEmpty(json.RoleName, GetDefaultRoleName())

	user := &model.User{}
	copier.Copy(user, json)
	user.DivisionID = sql.NullInt64{Int64: int64(json.DivisionID), Valid: json.DivisionID != 0}
	user.CreatedBy = operator

	if err := tx.Create(user).Error; err != nil {
		logger.Logger.Debugf("CreateUserErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func UpdateUser(id uint, json *model.UpdateUserRequest, operator uint) (*model.User, error) {
	return TxUpdateUser(database.DB, id, json, operator)
}

func TxUpdateUser(tx *gorm.DB, id uint, json *model.UpdateUserRequest, operator uint) (*model.User, error) {
	if json.Password != "" {
		salt, _ := bcrypt.Salt(10)
		hash, _ := bcrypt.Hash(json.Password, salt)
		json.Password = string(hash)
	}

	user := &model.User{}
	copier.Copy(user, json)
	user.ID = id
	user.UpdatedBy = operator
	user.DivisionID = sql.NullInt64{Int64: 0, Valid: false}
	tx = tx.Model(user).Updates(user)
	if json.DivisionID != 0 {
		tx = tx.Update("division_id", sql.NullInt64{Int64: int64(json.DivisionID), Valid: json.DivisionID != -1})
	}
	if err := tx.Error; err != nil {
		logger.Logger.Debugf("UpdateUserErr: %v\n", err)
		return nil, err
	}
	return user, nil
}

func AttachOpenIDToUser(id uint, openid string) error {
	return TxAttachOpenIDToUser(database.DB, id, openid)
}

func TxAttachOpenIDToUser(tx *gorm.DB, id uint, openid string) error {
	user := &model.User{}
	user.ID = id
	if err := tx.Model(user).Update("open_id", openid).Error; err != nil {
		logger.Logger.Debugf("AttachOpenIDToUserErr: %v\n", err)
		return err
	}
	return nil
}

func DeleteUser(id uint) error {
	return TxDeleteUser(database.DB, id)
}

func TxDeleteUser(tx *gorm.DB, id uint) (err error) {
	if err = tx.Delete(&model.User{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteUserByIdErr: %v\n", err)
	}
	return
}

func CheckLogin(user *model.User, password string) error {
	if ok := bcrypt.Match(password, user.Password); !ok {
		return fmt.Errorf("Wrong password")
	}
	return ForceLogin(user.ID, user.LoginIP)
}

func ForceLogin(id uint, ip string) error {
	return TxForceLogin(database.DB, id, ip)
}

func TxForceLogin(tx *gorm.DB, id uint, ip string) error {
	user := &model.User{
		LoginIP:   ip,
		LoginTime: time.Now(),
	}
	user.ID = id
	if err := tx.Model(user).Updates(user).Error; err != nil {
		logger.Logger.Debugf("ForceLoginErr: %v\n", err)
		return err
	}
	return nil
}
