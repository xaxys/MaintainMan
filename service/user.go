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

func GetUserByID(id uint, auth *model.AuthInfo) *model.ApiJson {
	user, err := dao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(UserToJson(user), "获取成功")
}

func GetUserInfoByID(id uint, auth *model.AuthInfo) *model.ApiJson {
	user, err := dao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	json := UserToJson(user)
	json.Role = dao.GetRole(user.RoleName)
	return model.Success(json, "获取成功")
}

func GetUserByName(name string, auth *model.AuthInfo) *model.ApiJson {
	user, err := dao.GetUserByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(UserToJson(user), "获取成功")
}

func GetUserInfoByName(name string, auth *model.AuthInfo) *model.ApiJson {
	user, err := dao.GetUserByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	json := UserToJson(user)
	json.Role = dao.GetRole(user.RoleName)
	return model.Success(json, "获取成功")
}

func RegisterUser(aul *model.RegisterUserRequest, auth *model.AuthInfo) *model.ApiJson {
	req := &model.CreateUserRequest{
		RegisterUserRequest: *aul,
	}
	return CreateUser(req, auth)
}

func CreateUser(aul *model.CreateUserRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if util.EmailRegex.MatchString(aul.Name) || util.PhoneRegex.MatchString(aul.Name) {
		return model.ErrorValidation(errors.New("用户名不能为邮箱或手机号"))
	}
	operator := util.NilOrBaseValue(auth, func(v *model.AuthInfo) uint { return v.User }, 0)
	u, err := dao.CreateUser(aul, operator)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(UserToJson(u), "创建成功")

}

func UpdateUser(id uint, aul *model.UpdateUserRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	_, err := dao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	u, err := dao.UpdateUser(id, aul, auth.User)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(UserToJson(u), "更新成功")
}

func DeleteUser(id uint, auth *model.AuthInfo) *model.ApiJson {
	if err := dao.DeleteUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func GetAllUsers(aul *model.AllUserRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	users, err := dao.GetAllUsersWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	us := util.TransSlice(users, UserToJson)
	return model.Success(us, "获取成功")
}

func WxUserLogin(aul *model.WxLoginRequest, ip string, auth *model.AuthInfo) *model.ApiJson {
	const wxURL = "https://api.weixin.qq.com/sns/jscode2session"
	params := map[string]string{
		"appid":      config.AppConfig.GetString("wechat.appid"),
		"secret":     config.AppConfig.GetString("wechat.secret"),
		"js_code":    aul.Code,
		"grant_type": "authorization_code",
	}
	wxres := &model.WxLoginResponse{}
	if err := util.HttpRequest(wxURL, "GET", params, wxres); err != nil {
		return model.ErrorVerification(err)
	}
	if wxres.ErrCode != 0 {
		return model.ErrorVerification(fmt.Errorf(wxres.ErrMsg))
	}

	user, err := dao.GetUserByOpenID(wxres.OpenID)
	id := user.ID
	if err != nil {
		// If user related to openid not found, attach openid to current user OR create a new one
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorQueryDatabase(err)
		}
		if auth.User != 0 {
			// If already login, attach openid to current user
			if err := dao.AttachOpenIDToUser(auth.User, wxres.OpenID); err != nil {
				return model.ErrorUpdateDatabase(err)
			}
			id = auth.User
		} else {
			// If not login, create a new user
			aul := &model.CreateUserRequest{
				RegisterUserRequest: model.RegisterUserRequest{
					Name:     wxres.OpenID,
					Password: util.RandomString(32),
				},
			}
			user, err = dao.CreateUser(aul, auth.User)
			if err != nil {
				return model.ErrorInsertDatabase(err)
			}
			if err := dao.AttachOpenIDToUser(user.ID, wxres.OpenID); err != nil {
				return model.ErrorUpdateDatabase(err)
			}
			id = user.ID
		}
	}

	if err := dao.ForceLogin(id, ip); err != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("登录失败"))
	}
	token, err := util.GetJwtString(id, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func UserLogin(aul *model.LoginRequest, ip string, auth *model.AuthInfo) *model.ApiJson {
	var user *model.User
	var err error
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
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

	user.LoginIP = ip
	if err := dao.CheckLogin(user, aul.Password); err != nil {
		return model.ErrorVerification(fmt.Errorf("密码错误"))
	}
	token, err := util.GetJwtString(user.ID, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func UserRenew(id uint, ip string, auth *model.AuthInfo) *model.ApiJson {
	user, err := dao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if err := dao.ForceLogin(id, ip); err != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("登录失败"))
	}
	token, err := util.GetJwtString(id, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func UserToJson(user *model.User) *model.UserJson {
	return util.NilOrValue(user, &model.UserJson{
		ID:          user.ID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		RoleName:    user.RoleName,
	})
}
