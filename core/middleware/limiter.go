package middleware

import (
	"github.com/kataras/iris/v12"
)

var (
	RateLimiter iris.Handler
)

func init() {
	// tRate := config.ImageConfig.GetInt("upload.throttling.rate")
	// burst := config.ImageConfig.GetInt("upload.throttling.burst")
	// purge := config.ImageConfig.GetDuration("upload.throttling.purge")
	// expire := config.ImageConfig.GetDuration("upload.throttling.expire")
	// RateLimiter = rate.Limit(float64(tRate), burst, rate.PurgeEvery(purge, expire))
}
