package dao

import (
	"fmt"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"
	"time"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/copier"
)

func GetUserByID(id uint) (*model.User, error) {
	user := &model.User{}

	if err := database.DB.First(user, id).Error; err != nil {
		logger.Logger.Debugf("GetUserByIDErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByName(name string) (*model.User, error) {
	user := &model.User{Name: name}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByNameErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{Email: email}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByEmailErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByPhone(phone string) (*model.User, error) {
	user := &model.User{Phone: phone}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByPhoneErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByOpenID(openid string) (*model.User, error) {
	user := &model.User{OpenID: openid}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		logger.Logger.Debugf("GetUserByOpenIDErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetAllUsersWithParam(aul *model.AllUserRequest) (users []*model.User, err error) {
	user := &model.User{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
	}
	if err = Filter(aul.OrderBy, aul.Offset, aul.Limit).Where(user).Find(&users).Error; err != nil {
		logger.Logger.Debugf("GetAllUserErr: %v\n", err)
	}
	return
}

func CreateUser(json *model.CreateUserRequest, operator uint) (*model.User, error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(json.Password, salt)
	json.Password = string(hash)
	if json.DisplayName == "" {
		json.DisplayName = json.Name
	}
	if json.RoleName == "" {
		json.RoleName = GetDefaultRoleName()
	}

	user := &model.User{}
	copier.Copy(user, json)
	user.CreatedBy = operator

	if err := database.DB.Create(user).Error; err != nil {
		logger.Logger.Debugf("CreateUserErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func UpdateUser(id uint, json *model.UpdateUserRequest, operator uint) (*model.User, error) {
	if json.Password != "" {
		salt, _ := bcrypt.Salt(10)
		hash, _ := bcrypt.Hash(json.Password, salt)
		json.Password = string(hash)
	}

	user := &model.User{}
	copier.Copy(user, json)
	user.ID = id
	user.UpdatedBy = operator
	if err := database.DB.Model(user).Updates(user).Error; err != nil {
		logger.Logger.Debugf("UpdateUserErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func AttachOpenIDToUser(id uint, openid string) error {
	user := &model.User{}
	user.ID = id
	if err := database.DB.Where(user).Update("openid", openid).Error; err != nil {
		logger.Logger.Debugf("AttachOpenIDToUserErr: %v\n", err)
		return err
	}
	return nil
}

func DeleteUser(id uint) error {
	if err := database.DB.Delete(&model.User{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteUserByIdErr: %v\n", err)
		return err
	}
	return nil
}

func CheckLogin(user *model.User, password string) error {
	if ok := bcrypt.Match(password, user.Password); !ok {
		return fmt.Errorf("Wrong password")
	}
	return ForceLogin(user.ID, user.LoginIP)
}

func ForceLogin(id uint, ip string) error {
	u := &model.User{
		LoginIP:   ip,
		LoginTime: time.Now(),
	}
	u.ID = id
	if err := database.DB.Model(u).Updates(u).Error; err != nil {
		logger.Logger.Debugf("ForceLoginErr: %v\n", err)
		return err
	}
	return nil
}
