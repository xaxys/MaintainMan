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
	AppConfig.SetDefault("app.page.limit", 100)
	AppConfig.SetDefault("app.page.default", 50)
	AppConfig.SetDefault("app.hit_expire_time", "12h")
	AppConfig.SetDefault("wechat.appid", "微信小程序的appid")
	AppConfig.SetDefault("wechat.secret", "微信小程序的secret")
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

	AppConfig.SetDefault("cache.expire", "24h")
	AppConfig.SetDefault("cache.purge", "10m")

	AppConfig.SetDefault("admin.name", "admin")
	AppConfig.SetDefault("admin.display_name", "maintainman default admin")
	AppConfig.SetDefault("admin.password", "123456")
	AppConfig.SetDefault("admin.role_name", "admin")

	if err := AppConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("App configuration file not found: %v\n", err)
			if err := AppConfig.SafeWriteConfig(); err != nil {
				panic(fmt.Errorf("Failed to write default app config: %v", err))
			}
			fmt.Println("Default app config file created.")
		} else {
			panic(fmt.Errorf("Fatal error reading config file: %v", err))
		}
	}
}
