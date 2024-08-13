package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
)

type InMemoryCacheStore[T any] struct {
	cache             *cache.Cache
	defaultExpiration time.Duration
	cleanupExpiration time.Duration
}

func NewInMemoryCacheStore[T any](exp, cleanup time.Duration) *InMemoryCacheStore[T] {
	return &InMemoryCacheStore[T]{
		cache:             cache.New(exp, cleanup),
		defaultExpiration: exp,
		cleanupExpiration: cleanup,
	}
}

func (s *InMemoryCacheStore[T]) Set(key string, value T) {
	s.cache.Set(key, value, s.defaultExpiration)
}

func (s *InMemoryCacheStore[T]) Get(key string) (*T, bool) {
	v, ok := s.cache.Get(key)
	if !ok {
		return nil, false
	}
	return lo.ToPtr(v.(T)), true
}
