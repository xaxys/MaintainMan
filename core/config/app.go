package config

import (
	"github.com/spf13/viper"
)

const AppConfigVersion = "1.3.1"

var (
	AppConfig *viper.Viper
)

func init() {
	AppConfig = viper.New()
	AppConfig.SetConfigName("app")
	AppConfig.SetConfigType("yaml")
	AppConfig.AddConfigPath(".")
	AppConfig.AddConfigPath("./config")
	AppConfig.AddConfigPath("/etc/maintainman/")
	AppConfig.AddConfigPath("$HOME/.maintainman/")

	AppConfig.SetDefault("app.name", "maintainman")
	AppConfig.SetDefault("app.listen", ":8787")
	AppConfig.SetDefault("app.loglevel", "info")
	AppConfig.SetDefault("app.page.limit", 100)
	AppConfig.SetDefault("app.page.default", 50)

	AppConfig.SetDefault("token.key", "xaxys_2022_all_rights_reserved")
	AppConfig.SetDefault("token.expire", "30m")

	AppConfig.SetDefault("database.driver", "sqlite")
	AppConfig.SetDefault("database.sqlite.path", "maintainman.db")
	AppConfig.SetDefault("database.mysql.host", "localhost")
	AppConfig.SetDefault("database.mysql.port", 3306)
	AppConfig.SetDefault("database.mysql.name", "maintainman")
	AppConfig.SetDefault("database.mysql.params", "parseTime=true&loc=Local&charset=utf8mb4")
	AppConfig.SetDefault("database.mysql.user", "root")
	AppConfig.SetDefault("database.mysql.password", "123456")

	AppConfig.SetDefault("cache.driver", "local")
	AppConfig.SetDefault("cache.limit", 268435456)
	AppConfig.SetDefault("cache.redis.host", "localhost")
	AppConfig.SetDefault("cache.redis.port", 6379)
	AppConfig.SetDefault("cache.redis.password", "")

	AppConfig.SetDefault("storage.driver", "local")
	AppConfig.SetDefault("storage.local.path", "./files")
	AppConfig.SetDefault("storage.s3.access_key", "ACCESS_KEY")
	AppConfig.SetDefault("storage.s3.secret_key", "SECRET_KEY")
	AppConfig.SetDefault("storage.s3.bucket", "BUCKET")
	AppConfig.SetDefault("storage.s3.region", "REGION")

	AppConfig.SetDefault("throttling.enable", false)
	AppConfig.SetDefault("throttling.burst", 100)
	AppConfig.SetDefault("throttling.rate", 10)
	AppConfig.SetDefault("throttling.purge", "1m")
	AppConfig.SetDefault("throttling.expire", "10m")

	AppConfig.SetDefault("bus_buffer", 1000)

	ReadAndUpdateConfig(AppConfig, "app", AppConfigVersion)
}
