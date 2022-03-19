package service

import (
	"errors"
	"fmt"
	"maintainman/config"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"

	"gorm.io/gorm"
)

func GetUserByID(id uint) *model.ApiJson {
	user, err := dao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(UserToJson(user), "获取成功")
}

func GetUserInfoByID(id uint) *model.ApiJson {
	user, err := dao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	json := UserToJson(user)
	json.Role = dao.GetRole(user.RoleName)
	return model.Success(json, "获取成功")
}

func GetUserByName(name string) *model.ApiJson {
	user, err := dao.GetUserByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(UserToJson(user), "获取成功")
}

func GetUserInfoByName(name string) *model.ApiJson {
	user, err := dao.GetUserByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	json := UserToJson(user)
	json.Role = dao.GetRole(user.RoleName)
	return model.Success(json, "获取成功")
}

func CreateUser(aul *model.ModifyUserJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	if util.EmailRegex.MatchString(aul.Name) || util.PhoneRegex.MatchString(aul.Name) {
		return model.ErrorVerification(errors.New("用户名不能为邮箱或手机号"))
	}
	u, err := dao.CreateUser(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)

	}
	return model.SuccessCreate(UserToJson(u), "创建成功")

}

func UpdateUser(id uint, aul *model.ModifyUserJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	u, err := dao.UpdateUser(id, aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorUpdateDatabase(err)
		}
	}
	u.Password = ""
	return model.SuccessUpdate(UserToJson(u), "更新成功")
}

func DeleteUser(id uint) *model.ApiJson {
	if err := dao.DeleteUserByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorDeleteDatabase(err)
		}
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func GetAllUsers(aul *model.AllUserJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	if aul.Limit == 0 {
		aul.Limit = config.AppConfig.GetInt("page_limit_default")
	}
	users, err := dao.GetAllUsersWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	us := util.TransSlice(users, UserToJson)
	return model.Success(us, "获取成功")
}

func UserLogin(aul *model.LoginJson) *model.ApiJson {
	var user *model.User
	var err error
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	if util.EmailRegex.MatchString(aul.Account) {
		user, err = dao.GetUserByEmail(aul.Account)
		if err != nil {
			return model.ErrorNotFound(fmt.Errorf("邮箱不存在"))
		}
	} else if util.PhoneRegex.MatchString(aul.Account) {
		user, err = dao.GetUserByPhone(aul.Account)
		if err != nil {
			return model.ErrorNotFound(fmt.Errorf("手机号不存在"))
		}
	} else {
		user, err = dao.GetUserByName(aul.Account)
		if err != nil {
			return model.ErrorNotFound(fmt.Errorf("用户名不存在"))
		}
	}

	if err := dao.CheckLogin(user, aul.Password); err != nil {
		return model.ErrorUnauthorized(fmt.Errorf("密码错误"))
	}
	token, err := util.GetJwtString(user.ID, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func UserRenew(uid uint) *model.ApiJson {
	user, err := dao.GetUserByID(uid)
	token, err := util.GetJwtString(user.ID, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func UserToJson(user *model.User) *model.UserJson {
	return &model.UserJson{
		ID:          user.ID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		RoleName:    user.RoleName,
	}
}
