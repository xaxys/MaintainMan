package user

import (
	"errors"
	"fmt"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/rbac"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func getUserByIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	user, err := dbGetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(userToJson(user), "获取成功")
}

func getUserInfoByIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	user, err := dbGetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	json := userToJson(user)
	json.Role = rbac.GetRole(user.RoleName)
	return model.Success(json, "获取成功")
}

func getUserByNameService(name string, auth *model.AuthInfo) *model.ApiJson {
	user, err := dbGetUserByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(userToJson(user), "获取成功")
}

func getUserInfoByNameService(name string, auth *model.AuthInfo) *model.ApiJson {
	user, err := dbGetUserByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	json := userToJson(user)
	json.Role = rbac.GetRole(user.RoleName)
	return model.Success(json, "获取成功")
}

func getUsersByDivisionService(id uint, param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(param); err != nil {
		return model.ErrorValidation(err)
	}
	users, count, err := dbGetUsersByDivision(id, param)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	us := util.TransSlice(users, userToJson)
	return model.SuccessPaged(us, count, "获取成功")
}

func registerUserService(aul *RegisterUserRequest, auth *model.AuthInfo) *model.ApiJson {
	req := &CreateUserRequest{
		RegisterUserRequest: *aul,
	}
	return createUserService(req, auth)
}

func createUserService(aul *CreateUserRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if util.EmailRegex.MatchString(aul.Name) || util.PhoneRegex.MatchString(aul.Name) {
		return model.ErrorValidation(fmt.Errorf("用户名不能为邮箱或手机号"))
	}
	operator := util.NilOrBaseValue(auth, func(v *model.AuthInfo) uint { return v.User }, 0)
	u, err := dbCreateUser(aul, operator)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(userToJson(u), "创建成功")

}

func updateUserService(id uint, aul *UpdateUserRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	_, err := dbGetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	u, err := dbUpdateUser(id, aul, auth.User)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(userToJson(u), "更新成功")
}

func deleteUserService(id uint, auth *model.AuthInfo) *model.ApiJson {
	if err := dbDeleteUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func getAllUsersService(aul *AllUserRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	users, count, err := dbGetAllUsersWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	us := util.TransSlice(users, userToJson)
	return model.SuccessPaged(us, count, "获取成功")
}

const wxURL = "https://api.weixin.qq.com/sns/jscode2session"

func wxUserLoginService(aul *WxLoginRequest, ip string, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}

	openID := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string {
		if id := v.Other["openid"]; id != nil {
			if openID, ok := id.(string); ok {
				return openID
			}
		}
		return ""
	}, "")
	if aul.Code != "" {
		id, err := getWxUserOpenID(aul.Code)
		if err != nil {
			return model.ErrorValidation(err)
		}
		openID = id
	}
	if openID == "" {
		return model.ErrorInvalidData(fmt.Errorf("未获取到openid"))
	}

	id := uint(0)
	user, err := dbGetUserByOpenID(openID)
	if err != nil {
		// If user related to openid not found, attach openid to current user OR create a new one
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorQueryDatabase(err)
		}
		if auth != nil && auth.User != 0 {
			// If already login, attach openid to current user
			if err := dbAttachOpenIDToUser(auth.User, openID); err != nil {
				return model.ErrorUpdateDatabase(err)
			}
			id = auth.User
		} else if userConfig.GetBool("wechat.fastlogin") {
			// If not login, create a new user
			aul := &CreateUserRequest{
				RegisterUserRequest: RegisterUserRequest{
					Name:     "微信用户" + openID,
					Password: util.RandomString(32),
				},
			}
			operator := util.NilOrBaseValue(auth, func(v *model.AuthInfo) uint { return v.User }, 0)
			u, response := createUserWithOpenID(aul, openID, operator)
			if response != nil {
				return response
			}
			id = u.ID
		} else {
			jwt, err := util.GetJwtStringWithClaims(0, "未登录用户", "", map[string]any{"openid": openID})
			if err != nil {
				return model.ErrorBuildJWT(err)
			}
			return model.Fail(jwt, "未绑定微信账号，请先绑定微信账号")
		}
	} else {
		id = user.ID
	}

	if err := dbForceLogin(id, ip); err != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("登录失败"))
	}
	token, err := util.GetJwtString(id, user.Name, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func wxUserRegisterService(aul *WxRegisterRequest, ip string, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if util.EmailRegex.MatchString(aul.Name) || util.PhoneRegex.MatchString(aul.Name) {
		return model.ErrorValidation(errors.New("用户名不能为邮箱或手机号"))
	}

	openID := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string {
		if id := v.Other["openid"]; id != nil {
			if openID, ok := id.(string); ok {
				return openID
			}
		}
		return ""
	}, "")
	if aul.Code != "" {
		id, err := getWxUserOpenID(aul.Code)
		if err != nil {
			return model.ErrorValidation(err)
		}
		openID = id
	}
	if openID == "" {
		return model.ErrorInvalidData(fmt.Errorf("未获取到openid"))
	}

	req := &CreateUserRequest{RegisterUserRequest: aul.RegisterUserRequest}
	operator := util.NilOrBaseValue(auth, func(v *model.AuthInfo) uint { return v.User }, 0)
	user, response := createUserWithOpenID(req, openID, operator)
	if response != nil {
		return response
	}

	if err := dbForceLogin(user.ID, ip); err != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("登录失败"))
	}
	token, err := util.GetJwtString(user.ID, user.Name, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func userLoginService(aul *LoginRequest, ip string, auth *model.AuthInfo) *model.ApiJson {
	var user *User
	var err error
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	if util.EmailRegex.MatchString(aul.Account) {
		user, err = dbGetUserByEmail(aul.Account)
		if err != nil {
			return model.ErrorNotFound(fmt.Errorf("邮箱不存在"))
		}
	} else if util.PhoneRegex.MatchString(aul.Account) {
		user, err = dbGetUserByPhone(aul.Account)
		if err != nil {
			return model.ErrorNotFound(fmt.Errorf("手机号不存在"))
		}
	} else {
		user, err = dbGetUserByName(aul.Account)
		if err != nil {
			return model.ErrorNotFound(fmt.Errorf("用户名不存在"))
		}
	}

	user.LoginIP = ip
	if err := dbCheckLogin(user, aul.Password); err != nil {
		return model.ErrorVerification(fmt.Errorf("密码错误"))
	}
	token, err := util.GetJwtString(user.ID, user.Name, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	openID := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string {
		if id := v.Other["openid"]; id != nil {
			if openID, ok := id.(string); ok {
				return openID
			}
		}
		return ""
	}, "")
	if openID != "" && user.OpenID == "" {
		dbAttachOpenIDToUser(user.ID, openID)
	}
	return model.Success(token, "登陆成功")
}

func userRenewService(id uint, ip string, auth *model.AuthInfo) *model.ApiJson {
	user, err := dbGetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if err := dbForceLogin(id, ip); err != nil {
		return model.ErrorUpdateDatabase(fmt.Errorf("登录失败"))
	}
	token, err := util.GetJwtString(id, user.Name, user.RoleName)
	if err != nil {
		return model.ErrorBuildJWT(err)
	}
	return model.Success(token, "登陆成功")
}

func getWxUserOpenID(code string) (string, error) {
	params := map[string]string{
		"appid":      userConfig.GetString("wechat.appid"),
		"secret":     userConfig.GetString("wechat.secret"),
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	wxres, err := util.HTTPRequest[WxLoginResponse](wxURL, "GET", params)
	if err != nil {
		mctx.Logger.Warnf("WeChatLoginErr: %+v", err)
		return "", err
	}
	if wxres.ErrCode != 0 {
		return "", fmt.Errorf(wxres.ErrMsg)
	}
	return wxres.OpenID, nil
}

func createUserWithOpenID(aul *CreateUserRequest, openID string, operator uint) (user *User, response *model.ApiJson) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		user, err := dbCreateUser(aul, operator)
		if err != nil {
			response = model.ErrorInsertDatabase(err)
			return err
		}
		if err := dbAttachOpenIDToUser(user.ID, openID); err != nil {
			response = model.ErrorUpdateDatabase(err)
			return err
		}
		return nil
	})
	return
}

func userToJson(user *User) *UserJson {
	if user == nil {
		return nil
	} else {
		return &UserJson{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			RoleName:    user.RoleName,
			Division:    divisionToJson(user.Division),
			Phone:       user.Phone,
			Email:       user.Email,
			RealName:    user.RealName,
			LoginTime:   user.LoginTime.Unix(),
		}
	}
}
