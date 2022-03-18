package dao

import (
	"fmt"
	"maintainman/database"
	"maintainman/model"
	. "maintainman/model"

	"github.com/jameskeane/bcrypt"
)

func GetUserByID(id uint) (*User, error) {
	user := &User{}

	if err := database.DB.First(user, id).Error; err != nil {
		fmt.Printf("GetUserByIDErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByName(name string) (*User, error) {
	user := &User{Name: name}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		fmt.Printf("GetUserByNameErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{Email: email}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		fmt.Printf("GetUserByEmailErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func GetUserByPhone(phone string) (*User, error) {
	user := &User{Phone: phone}

	if err := database.DB.Where(user).First(user).Error; err != nil {
		fmt.Printf("GetUserByPhoneErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func DeleteUserByID(id uint) error {
	if err := database.DB.Delete(&User{}, id).Error; err != nil {
		fmt.Printf("DeleteUserByIdErr: %v\n", err)
		return err
	}
	return nil
}

func GetAllUsers() ([]*User, error) {
	return GetAllUsersWithParam("", "", "", 0, 0)
}

func GetAllUsersWithParam(name, displayName, orderBy string, offset, limit int) (users []*User, err error) {
	user := &User{
		Name:        name,
		DisplayName: displayName,
	}
	if err = database.DB.Where(user).Find(&users).Error; err != nil {
		fmt.Printf("GetAllUserErr: %v\n", err)
	}
	return
}

func CreateUser(ujson *ModifyUserJson) (*User, error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(ujson.Password, salt)
	ujson.Password = string(hash)
	if ujson.DisplayName == "" {
		ujson.DisplayName = ujson.Name
	}
	if ujson.RoleName == "" {
		ujson.RoleName = GetDefaultRoleName()
	}

	user := JsonToUser(ujson)

	if err := database.DB.Create(user).Error; err != nil {
		fmt.Printf("CreateUserErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func UpdateUser(id uint, ujson *ModifyUserJson) (*User, error) {
	user := JsonToUser(ujson)
	user.ID = id
	if ujson.Password != "" {
		salt, _ := bcrypt.Salt(10)
		hash, _ := bcrypt.Hash(ujson.Password, salt)
		user.Password = string(hash)
	}

	if err := database.DB.Model(&user).Updates(user).Error; err != nil {
		fmt.Printf("UpdateUserErr: %v\n", err)
		return nil, err
	}

	return user, nil
}

func CheckLogin(user *User, password string) error {
	if ok := bcrypt.Match(password, user.Password); !ok {
		return fmt.Errorf("Wrong password")
	}
	return nil
}

func JsonToUser(json *model.ModifyUserJson) (user *model.User) {
	return &model.User{
		Name:        json.Name,
		Password:    json.Password,
		DisplayName: json.DisplayName,
		RoleName:    json.RoleName,
		DivisionID:  json.DivisionID,
		Phone:       json.Phone,
		Email:       json.Email,
		RealName:    json.RealName,
	}
}
