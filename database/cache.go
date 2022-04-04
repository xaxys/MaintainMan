package database

import (
	"context"
	"fmt"
	"maintainman/config"
	"maintainman/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
)

type CacheClient interface {
	Get(key string) (any, bool)
	Set(key string, value any, expire time.Duration)
}

type GoCacheClient struct {
	cache *cache.Cache
}

type RedisClient struct {
	rdb *redis.Client
}

var (
	Cache CacheClient
)

func init() {
	cacheType := config.AppConfig.GetString("cache.driver")
	switch cacheType {
	case "go-cache":
		Cache = initGoCache()
	case "redis":
		Cache = initRedis()
	default:
		panic(fmt.Errorf("support go-cache and redis only"))
	}
}

func initGoCache() CacheClient {
	cacheExpire, err := time.ParseDuration(config.AppConfig.GetString("cache.expire"))
	if err != nil {
		panic(fmt.Errorf("Can not parse cache.expire in app config: %v", err))
	}
	cachePurge, err := time.ParseDuration(config.AppConfig.GetString("cache.purge"))
	if err != nil {
		panic(fmt.Errorf("Can not parse cache.purge in app config: %v", err))
	}

	cache := &GoCacheClient{
		cache: cache.New(cacheExpire, cachePurge),
	}
	return cache
}

func initRedis() CacheClient {
	rdbHost := config.AppConfig.GetString("cache.redis.host")
	rdbPort := config.AppConfig.GetInt("cache.redis.port")
	rdbPasswd := config.AppConfig.GetString("cache.redis.password")
	rdbAddr := fmt.Sprintf("%s:%d", rdbHost, rdbPort)

	cache := &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr:     rdbAddr,
			Password: rdbPasswd,
			DB:       0,
		}),
	}
	return cache
}

func (client *GoCacheClient) Get(key string) (any, bool) {
	return client.cache.Get(key)
}

func (client *GoCacheClient) Set(key string, value any, expire time.Duration) {
	client.cache.Set(key, value, expire)
}

func (client *RedisClient) Get(key string) (any, bool) {
	ctx := context.Background()
	value, err := client.rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		logger.Logger.Warnf("Redis error: %+v", err)
	}
	return value, err == nil
}

func (client *RedisClient) Set(key string, value any, expire time.Duration) {
	ctx := context.Background()
	if _, err := client.rdb.Set(ctx, key, value, expire).Result(); err != nil {
		logger.Logger.Warnf("Redis error: %+v", err)
	}
}
