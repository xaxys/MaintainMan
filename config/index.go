package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	AppConfig *viper.Viper
)

func init() {
	AppConfig = viper.New()
	AppConfig.SetConfigName("app")
	AppConfig.SetConfigType("yaml")
	AppConfig.AddConfigPath(".")
	AppConfig.AddConfigPath("./config")
	AppConfig.AddConfigPath("/etc/srs_wrappper/")
	AppConfig.AddConfigPath("$HOME/.srs_wrappper/")

	AppConfig.SetDefault("app.name", "maintainman")
	AppConfig.SetDefault("app.listen", ":8787")
	AppConfig.SetDefault("app.loglevel", "info")
	AppConfig.SetDefault("token.key", "xaxys_2022_all_rights_reserved")
	AppConfig.SetDefault("token.exp", "30m")

	AppConfig.SetDefault("database.driver", "sqlite")
	AppConfig.SetDefault("database.sqlite.path", "maintainman.db")
	AppConfig.SetDefault("database.mysql.host", "localhost")
	AppConfig.SetDefault("database.mysql.port", "3306")
	AppConfig.SetDefault("database.mysql.name", "maintainman")
	AppConfig.SetDefault("database.mysql.params", "parseTime=true&loc=Local&charset=utf8mb4")
	AppConfig.SetDefault("database.mysql.user", "root")
	AppConfig.SetDefault("database.mysql.password", "123456")

	AppConfig.SetDefault("cache.expire", 86400)
	AppConfig.SetDefault("cache.purge", 600)

	AppConfig.SetDefault("admin.name", "admin")
	AppConfig.SetDefault("admin.display_name", "maintainman default admin")
	AppConfig.SetDefault("admin.password", "123456")
	AppConfig.SetDefault("admin.role_name", "admin")

	if err := AppConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			AppConfig.SafeWriteConfig()
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
}
