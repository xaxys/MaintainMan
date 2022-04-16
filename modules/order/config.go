package order

import "github.com/spf13/viper"

var orderConfig = viper.New()

func init() {
	orderConfig.SetDefault("hit_expire.announce", "12h")
	orderConfig.SetDefault("appraise.timeout", "72h")
	orderConfig.SetDefault("appraise.purge", "1m")
	orderConfig.SetDefault("appraise.default", 5)
	orderConfig.SetDefault("item_can_negative", true)
}
