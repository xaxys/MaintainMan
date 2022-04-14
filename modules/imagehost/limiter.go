package imagehost

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
)

var (
	rateLimiter iris.Handler
)

func initLimiter() {
	tRate := imageConfig.GetInt("upload.throttling.rate")
	burst := imageConfig.GetInt("upload.throttling.burst")
	purge := imageConfig.GetDuration("upload.throttling.purge")
	expire := imageConfig.GetDuration("upload.throttling.expire")
	rateLimiter = rate.Limit(float64(tRate), burst, rate.PurgeEvery(purge, expire))
}
