package model

import "time"

type User struct {
	BaseModel
	Name        string    `gorm:"unique; VARCHAR(50) comment '用户名'"`
	Password    string    `gorm:"not null VARCHAR(191) comment '密码'"`
	DisplayName string    `gorm:"not null VARCHAR(191) comment '昵称'"`
	RoleName    string    `gorm:"not null VARCHAR(191) comment '所属角色'"`
	DivisionID  uint      `gorm:"not null default '0' comment '所属分组id'"`
	Phone       string    `gorm:"VARCHAR(191) default NULL comment '手机号'"`
	Email       string    `gorm:"VARCHAR(191) default NULL comment '邮箱'"`
	LoginIP     string    `gorm:"not null VARCHAR(40) default '0.0.0.0' column:login_ip comment '最后登录IP'"`
	LoginTime   time.Time `gorm:"not null default '0000-00-00 00:00:00' comment '最后登录时间'"`
	RealName    string    `gorm:"not null VARCHAR(191) default '' comment '真实姓名'"`
}

type LoginJson struct {
	Name     string `json:"name" validate:"required,gte=2,lte=50"`
	Password string `json:"password" validate:"gte=8,lte=32"`
}

type ModifyUserJson struct {
	Name        string `json:"name" validate:"required,gte=4,lte=50"`
	Password    string `json:"password" validate:"required,gte=8,lte=32"`
	DisplayName string `json:"display_name" validate:"required,lte=191"`
	RoleName    string `json:"role_name" validate:"lte=191"`
	DivisionID  uint   `json:"division_id"`
	Phone       string `json:"phone" validate:"lte=191"`
	Email       string `json:"email" validate:"lte=191"`
	RealName    string `json:"real_name" validate:"lte=191"`
}

type AllUserJson struct {
	Name        string `json:"name" validate:"gte=2,lte=50"`
	DisplayName string `json:"display_name" validate:"gte=2,lte=50"`
	OrderBy     string `json:"order_by"`
	Limit       int    `json:"limit" validate:"number"`
	Offset      int    `json:"offset" validate:"number"`
}

type UserJson struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" validate:"required,gte=2,lte=50"`
	DisplayName string    `json:"display_name"`
	RoleName    string    `json:"user_role"`
	Role        *RoleJson `json:"role,omitempty"`
}
