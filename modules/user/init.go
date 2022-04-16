package user

import (
	"fmt"

	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/logger"
)

func initDefaultData() {
	createSystemAdmin()
}

func createSystemAdmin() {
	aul := &CreateUserRequest{}
	aul.Name = config.AppConfig.GetString("admin.name")
	aul.DisplayName = config.AppConfig.GetString("admin.display_name")
	aul.Password = config.AppConfig.GetString("admin.password")
	aul.RoleName = config.AppConfig.GetString("admin.role_name")

	count, err := dbGetUserCount()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		logger.Logger.Debug("Create default administrator account")
		if _, err := dbCreateUser(aul, 0); err != nil {
			panic(fmt.Errorf("failed to create default administrator: %v", err))
		}
	}
}
