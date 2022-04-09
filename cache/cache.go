package cache

import (
	"context"
	"fmt"
	"maintainman/config"
	"maintainman/logger"
	"maintainman/util"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	Cache      ICache
	ImageCache ICache
	redisConn  *redis.Client
)

type ICache interface {
	Get(key string) (any, bool)
	Set(key string, value any, expire time.Duration) bool
	SetWithCost(key string, value any, cost int64, expire time.Duration) bool
	Del(key string)
}

type Ristretto struct {
	limit int64
	cache *ristretto.Cache
}

// Default implemented cache strategy: LFU
type Redis struct {
	prefix  string
	limit   int64
	onEvict func(any) error
	rdb     *redis.Client
}

func init() {
	initRedisConn(config.AppConfig)
	Cache = initCache("app", config.AppConfig, nil)

}

func CreateImageCache(fn func(any) error) {
	ImageCache = initCache("image", config.ImageConfig, fn)
}

func initCache(name string, config *viper.Viper, fn func(any) error) ICache {
	cacheType := config.GetString("cache.driver")
	limit := config.GetInt64("cache.limit")
	switch cacheType {
	case "local":
		return newRistretto(limit, fn)
	case "redis":
		return newRedis(name, limit, fn)
	default:
		panic(fmt.Errorf("support local and redis only"))
	}
}

func initRedisConn(config *viper.Viper) {
	rdbHost := config.GetString("cache.redis.host")
	rdbPort := config.GetInt("cache.redis.port")
	rdbPasswd := config.GetString("cache.redis.password")
	rdbAddr := fmt.Sprintf("%s:%d", rdbHost, rdbPort)

	redisConn = redis.NewClient(&redis.Options{
		Addr:     rdbAddr,
		Password: rdbPasswd,
		DB:       0,
	})
}

func newRistretto(limit int64, onEvict func(any) error) ICache {
	max_cost := util.Tenary(limit > 0, limit, 1024)
	ristretto, err := ristretto.NewCache(&ristretto.Config{
		IgnoreInternalCost: true,
		NumCounters:        max_cost << 3,
		MaxCost:            max_cost,
		BufferItems:        64,
		OnEvict: func(item *ristretto.Item) {
			if onEvict != nil {
				onEvict(item.Value)
			}
		},
	})
	if err != nil {
		panic(fmt.Errorf("Can not create ristretto cache: %v", err))
	}
	cache := &Ristretto{
		limit: limit,
		cache: ristretto,
	}
	return cache
}

func newRedis(prefix string, limit int64, onEvict func(any) error) ICache {
	cache := &Redis{
		prefix:  prefix,
		limit:   limit,
		onEvict: onEvict,
		rdb:     redisConn,
	}
	return cache
}

func (client *Ristretto) Get(key string) (any, bool) {
	return client.cache.Get(key)
}

func (client *Ristretto) Set(key string, value any, expire time.Duration) bool {
	size := util.Tenary(client.limit > 0, int64(unsafe.Sizeof(value)), 0)
	return client.cache.SetWithTTL(key, value, size, expire)
}

func (client *Ristretto) SetWithCost(key string, value any, cost int64, expire time.Duration) bool {
	size := util.Tenary(client.limit > 0, cost, 0)
	return client.cache.SetWithTTL(key, value, size, expire)
}

func (client *Ristretto) Del(key string) {
	client.cache.Del(key)
}

func (client *Redis) Get(key string) (any, bool) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("%s:%s", client.prefix, key)
	value, err := client.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil && err != redis.Nil {
		logger.Logger.Warnf("Redis error: %+v", err)
	}
	if client.limit > 0 {
		if _, err := client.rdb.ZAdd(ctx, client.prefix+"timestamp", &redis.Z{Score: float64(time.Now().Unix()), Member: redisKey}).Result(); err != nil {
			logger.Logger.Warnf("Redis error: %+v", err)
		}
	}
	return value, true
}

func (client *Redis) Set(key string, value any, expire time.Duration) bool {
	size := util.Tenary(client.limit > 0, int64(unsafe.Sizeof(value)), 0)
	return client.SetWithCost(key, value, size, expire)
}

func (client *Redis) SetWithCost(key string, value any, cost int64, expire time.Duration) bool {
	ctx := context.Background()
	redisKey := fmt.Sprintf("%s:%s", client.prefix, key)
	if _, err := client.rdb.Set(ctx, redisKey, value, expire).Result(); err != nil {
		logger.Logger.Warnf("Redis error: %+v", err)
		return false
	}
	if cost != 0 && client.limit > 0 {
		if _, err := client.rdb.SetNX(ctx, client.prefix+"size", 0, 0).Result(); err != nil {
			logger.Logger.Warnf("Redis error: %+v", err)
		}
		totalSize, err := client.rdb.IncrBy(ctx, client.prefix+"size", cost).Result()
		if err != nil {
			logger.Logger.Warnf("Redis error: %+v", err)
		}
		if _, err := client.rdb.ZAdd(ctx, client.prefix+"timestamp", &redis.Z{Score: float64(time.Now().Unix()), Member: redisKey}).Result(); err != nil {
			logger.Logger.Warnf("Redis error: %+v", err)
		}

		if totalSize > client.limit {
			go func() {
				candidates, err := client.rdb.ZRange(ctx, client.prefix+"timestamp", 0, 5).Result()
				if err != nil {
					logger.Logger.Warnf("Redis error: %+v", err)
				}
				if client.onEvict != nil {
					for _, candidate := range candidates {
						value, err := client.rdb.Get(ctx, redisKey).Result()
						if err != nil && err != redis.Nil {
							logger.Logger.Warnf("Redis error: %+v", err)
						}
						if client.onEvict(value) != nil {
							logger.Logger.Debugf("Failed to run evict function on %s: %+v", candidate, err)
						}
					}
				}
				for _, candidate := range candidates {
					client.Del(candidate)
				}
			}()
		}
	}
	return true
}

func (client *Redis) Del(key string) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("%s:%s", client.prefix, key)
	value, err := client.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return
	}
	if err != nil && err != redis.Nil {
		logger.Logger.Warnf("Redis error: %+v", err)
	}
	if _, err := client.rdb.Del(ctx, redisKey).Result(); err != nil {
		logger.Logger.Warnf("Redis error: %+v", err)
	}
	if client.limit > 0 {
		if _, err := client.rdb.ZRem(ctx, client.prefix+"timestamp", redisKey).Result(); err != nil {
			logger.Logger.Warnf("Redis error: %+v", err)
		}
		if _, err := client.rdb.DecrBy(ctx, client.prefix+"size", int64(unsafe.Sizeof(value))).Result(); err != nil {
			logger.Logger.Warnf("Redis error: %+v", err)
		}
	}
}
