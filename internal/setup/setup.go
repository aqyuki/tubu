package setup

import (
	"time"

	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/config"
)

const (
	defaultExpiration = 30 * time.Minute
	defaultCleanup    = 30 * time.Minute
)

func NewCacheStore[T any](cnf *config.Config) cache.CacheStore[T] {
	return cache.NewInMemoryCacheStore[T](defaultExpiration, defaultCleanup)
}
