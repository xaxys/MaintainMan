package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/xaxys/maintainman/core/config"
)

var (
	RateLimiter iris.Handler
)

func init() {
	if config.AppConfig.GetBool("throttling.enable") {
		tRate := config.AppConfig.GetInt("throttling.rate")
		burst := config.AppConfig.GetInt("throttling.burst")
		purge := config.AppConfig.GetDuration("throttling.purge")
		expire := config.AppConfig.GetDuration("throttling.expire")
		RateLimiter = rate.Limit(float64(tRate), burst, rate.PurgeEvery(purge, expire))
	}
}
