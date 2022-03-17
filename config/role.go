package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	RoleConfig *viper.Viper
)

func init() {
	RoleConfig = viper.New()
	RoleConfig.SetConfigName("role")
	RoleConfig.SetConfigType("yaml")
	RoleConfig.AddConfigPath(".")
	RoleConfig.AddConfigPath("./config")
	RoleConfig.AddConfigPath("/etc/maintainman/")
	RoleConfig.AddConfigPath("$HOME/.maintainman/")

	if err := RoleConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Role configuration file not found: %v\n", err))
		} else {
			panic(fmt.Errorf("Fatal error reading config file: %v\n", err))
		}
	}
}
