package user

import "github.com/spf13/viper"

var userConfig = viper.New()

func init() {
	userConfig.SetDefault("wechat.appid", "微信小程序的appid")
	userConfig.SetDefault("wechat.secret", "微信小程序的secret")
	userConfig.SetDefault("wechat.fastlogin", true)
	userConfig.SetDefault("admin.name", "admin")
	userConfig.SetDefault("admin.display_name", "maintainman default admin")
	userConfig.SetDefault("admin.password", "12345678")
	userConfig.SetDefault("admin.role_name", "super_admin")
}
