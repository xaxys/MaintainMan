package config

import (
	"github.com/spf13/viper"
)

const AppConfigVersion = "1.0.3"

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
	AppConfig.SetDefault("app.hit_expire.announce", "12h")
	AppConfig.SetDefault("app.appraise.timeout", "72h")
	AppConfig.SetDefault("app.appraise.purge", "10m")
	AppConfig.SetDefault("app.appraise.default", 5)
	AppConfig.SetDefault("wechat.appid", "微信小程序的appid")
	AppConfig.SetDefault("wechat.secret", "微信小程序的secret")
	AppConfig.SetDefault("wechat.fastlogin", true)
	AppConfig.SetDefault("token.key", "xaxys_2022_all_rights_reserved")
	AppConfig.SetDefault("token.exp", "30m")

	AppConfig.SetDefault("database.driver", "sqlite")
	AppConfig.SetDefault("database.sqlite.path", "maintainman.db")
	AppConfig.SetDefault("database.mysql.host", "localhost")
	AppConfig.SetDefault("database.mysql.port", 3306)
	AppConfig.SetDefault("database.mysql.name", "maintainman")
	AppConfig.SetDefault("database.mysql.params", "parseTime=true&loc=Local&charset=utf8mb4")
	AppConfig.SetDefault("database.mysql.user", "root")
	AppConfig.SetDefault("database.mysql.password", "123456")

	AppConfig.SetDefault("cache.driver", "go-cache")
	AppConfig.SetDefault("cache.redis.host", "localhost")
	AppConfig.SetDefault("cache.redis.port", 6379)
	AppConfig.SetDefault("cache.redis.password", "")
	AppConfig.SetDefault("cache.expire", "24h")
	AppConfig.SetDefault("cache.purge", "10m")

	AppConfig.SetDefault("admin.name", "admin")
	AppConfig.SetDefault("admin.display_name", "maintainman default admin")
	AppConfig.SetDefault("admin.password", "12345678")
	AppConfig.SetDefault("admin.role_name", "super_admin")

	ReadAndUpdateConfig(AppConfig, "app", AppConfigVersion)
}
