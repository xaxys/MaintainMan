package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	PermConfig *viper.Viper
)

func init() {
	PermConfig = viper.New()
	PermConfig.SetConfigName("permission")
	PermConfig.SetConfigType("yaml")
	PermConfig.AddConfigPath(".")
	PermConfig.AddConfigPath("./config")
	PermConfig.AddConfigPath("/etc/maintainman/")
	PermConfig.AddConfigPath("$HOME/.maintainman/")

	if err := PermConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Permission configuration file not found: %v\n", err))
		} else {
			panic(fmt.Errorf("Fatal error reading config file: %v\n", err))
		}
	}
}
