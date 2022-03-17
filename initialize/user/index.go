package role

import (
	"errors"
	"fmt"

	"maintainman/config"
	"maintainman/dao"
	. "maintainman/model"

	"gorm.io/gorm"
)

func CreateDefaultUsers() {
	CreateSystemAdmin()
}

func CreateSystemAdmin() {
	aul := &ModifyUserJson{
		Name:        config.AppConfig.GetString("admin.name"),
		DisplayName: config.AppConfig.GetString("admin.display_name"),
		Password:    config.AppConfig.GetString("admin.password"),
		RoleName:    config.AppConfig.GetString("admin.role_name"),
	}

	if _, err := dao.GetUserByID(1); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Create default administrator account")
			if _, err := dao.CreateUser(aul); err != nil {
				panic("Failed to create default administrator")
			}
		} else {
			panic(err)
		}
	}
}
