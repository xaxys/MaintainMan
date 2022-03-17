package initialize

import (
	_ "maintainman/database"
	. "maintainman/initialize/user"
)

func InitDefaultData() {
	CreateDefaultUsers()
}
