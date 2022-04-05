package initialize

import (
	"fmt"
	"maintainman/config"
	"maintainman/dao"
	"maintainman/logger"
	"maintainman/model"
)

func InitDefaultData() {
	CreateSystemAdmin()
}

func CreateSystemAdmin() {
	aul := &model.CreateUserRequest{}
	aul.Name = config.AppConfig.GetString("admin.name")
	aul.DisplayName = config.AppConfig.GetString("admin.display_name")
	aul.Password = config.AppConfig.GetString("admin.password")
	aul.RoleName = config.AppConfig.GetString("admin.role_name")

	count, err := dao.GetUserCount()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		logger.Logger.Debug("Create default administrator account")
		if _, err := dao.CreateUser(aul, 0); err != nil {
			panic(fmt.Errorf("failed to create default administrator: %v", err))
		}
	}
}
