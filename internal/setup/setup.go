package setup

import (
	"time"

	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/config"
	"github.com/aqyuki/tubu/packages/infra/redis"
)

const (
	defaultExpiration = 30 * time.Minute
	defaultCleanup    = 30 * time.Minute
)

func NewCacheStore[T any](cnf *config.Config) cache.CacheStore[T] {
	if cnf.RedisEnabled {
		client, err := redis.NewRedisClient(&redis.RedisConfig{
			Address:  cnf.RedisAddr,
			Password: cnf.RedisPassword,
			DB:       cnf.RedisDB,
			PoolSize: cnf.RedisPoolSize,
		})
		if err == nil {
			return cache.NewRedisCacheStore[T](client, defaultExpiration)
		}
	}
	return cache.NewInMemoryCacheStore[T](defaultExpiration, defaultCleanup)
}
