package database

import (
	"fmt"
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
	cacheExpire, err := time.ParseDuration(config.AppConfig.GetString("cache.expire"))
	if err != nil {
		panic(fmt.Errorf("Can not parse cache.expire in app config: %v", err))
	}
	cachePurge, err := time.ParseDuration(config.AppConfig.GetString("cache.purge"))
	if err != nil {
		panic(fmt.Errorf("Can not parse cache.purge in app config: %v", err))
	}

	cache := cache.New(cacheExpire, cachePurge)
	return cache
}
