package announce

import "github.com/spf13/viper"

var announceConfig = viper.New()

func init() {
	announceConfig.SetDefault("hit_expire", "12h")
}
