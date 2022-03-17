package database

import (
	"maintainman/config"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	Cache *cache.Cache
)

func init() {
	Cache = initCache()
}

func initCache() *cache.Cache {
	cacheExpire := config.AppConfig.GetInt64("cache.expire")
	cachePurge := config.AppConfig.GetInt64("cache.purge")

	cache := cache.New(time.Duration(cacheExpire)*time.Second, time.Duration(cachePurge)*time.Second)
	return cache
}
