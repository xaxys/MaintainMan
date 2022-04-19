package user

import (
	"database/sql"
	"time"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/rbac"
)

type User struct {
	model.BaseModel
	Name        string        `gorm:"not null; size:50; unique; comment:用户名"`
	Password    string        `gorm:"not null; size:191; comment:密码"`
	DisplayName string        `gorm:"not null; size:191; comment:昵称"`
	RoleName    string        `gorm:"not null; size:50; index; comment:所属角色"`
	DivisionID  sql.NullInt64 `gorm:"comment:所属分组id"`
	Division    *Division     `gorm:"foreignkey:DivisionID"`
	Phone       string        `gorm:"not null; size:191; index; comment:手机号"`
	Email       string        `gorm:"not null; size:191; index; comment:邮箱"`
	LoginIP     string        `gorm:"not null; size:40; default:0.0.0.0; comment:最后登录IP"`
	LoginTime   time.Time     `gorm:"not null; comment:最后登录时间"`
	RealName    string        `gorm:"not null; size:191; comment:真实姓名"`
	OpenID      string        `gorm:"not null; size:191; index; comment:微信openid"` //TODO: 更改open_id的orm格式
}

type LoginRequest struct {
	Account  string `json:"account" validate:"required,lte=191"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

type WxLoginRequest struct {
	Code string `json:"code"`
}

type WxLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type RegisterUserRequest struct {
	Name        string `json:"name" validate:"required,gte=2,lte=50"`
	Password    string `json:"password" validate:"required,gte=8,lte=32"`
	DisplayName string `json:"display_name" validate:"required,lte=191"`
	Phone       string `json:"phone" validate:"omitempty,alphanum,lte=191"`
	Email       string `json:"email" validate:"omitempty,email,lte=191"`
	RealName    string `json:"real_name" validate:"omitempty,lte=191"`
}

type WxRegisterRequest struct {
	WxLoginRequest
	RegisterUserRequest
}

type CreateUserRequest struct {
	RegisterUserRequest
	RoleName   string `json:"role_name" validate:"omitempty,lte=50"`
	DivisionID uint   `json:"division_id"`
	OpenID     string `json:"-"`
}

type UpdateUserRequest struct {
	Name        string `json:"name" validate:"omitempty,gte=2,lte=50"`
	Password    string `json:"password" validate:"omitempty,gte=8,lte=32"`
	DisplayName string `json:"display_name" validate:"omitempty,lte=191"`
	Phone       string `json:"phone" validate:"omitempty,alphanum,lte=191"`
	Email       string `json:"email" validate:"omitempty,email,lte=191"`
	RealName    string `json:"real_name" validate:"omitempty,lte=191"`
	RoleName    string `json:"role_name" validate:"omitempty,lte=50"`
	DivisionID  int64  `json:"division_id" validate:"omitempty,gte=-1"` // -1: 修改为null 0: 不修改 n: 修改为指定的分组
}

type AllUserRequest struct {
	Name        string `json:"name" url:"name" validate:"omitempty,gte=2,lte=50"`
	DisplayName string `json:"display_name" url:"display_name" validate:"omitempty,lte=191"`
	model.PageParam
}

type UserJson struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"` // 昵称
	Division    *DivisionJson  `json:"division"`
	RoleName    string         `json:"user_role"`
	Role        *rbac.RoleJson `json:"role,omitempty"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	RealName    string         `json:"real_name"`
	LoginTime   int64          `json:"login_time"` // unix timestamp in seconds (UTC)
}
