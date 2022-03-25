package role

import (
	"errors"

	"maintainman/config"
	"maintainman/dao"
	"maintainman/logger"
	. "maintainman/model"

	"gorm.io/gorm"
)

func CreateDefaultUsers() {
	CreateSystemAdmin()
}

func CreateSystemAdmin() {
	aul := &CreateUserRequest{}
	aul.Name = config.AppConfig.GetString("admin.name")
	aul.DisplayName = config.AppConfig.GetString("admin.display_name")
	aul.Password = config.AppConfig.GetString("admin.password")
	aul.RoleName = config.AppConfig.GetString("admin.role_name")

	if _, err := dao.GetUserByID(1); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Debug("Create default administrator account")
			if _, err := dao.CreateUser(aul); err != nil {
				panic("Failed to create default administrator")
			}
		} else {
			panic(err)
		}
	}
}
