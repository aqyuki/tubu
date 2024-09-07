package setup

import (
	"time"

	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/profile"
)

const (
	defaultExpiration = 30 * time.Minute
	defaultCleanup    = 30 * time.Minute
)

func NewCacheStore[T any](prof *profile.Profile) cache.CacheStore[T] {
	return cache.NewInMemoryCacheStore[T](defaultExpiration, defaultCleanup)
}
