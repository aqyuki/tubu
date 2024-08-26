package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
)

var ErrNotFound = errors.New("cache: key not found")

type CacheStore[T any] interface {
	Set(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) (*T, error)
}

var _ CacheStore[string] = (*InMemoryCacheStore[string])(nil)

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

func (s *InMemoryCacheStore[T]) Set(_ context.Context, key string, value T) error {
	s.cache.Set(key, value, s.defaultExpiration)
	return nil
}

func (s *InMemoryCacheStore[T]) Get(_ context.Context, key string) (*T, error) {
	v, ok := s.cache.Get(key)
	if !ok {
		return nil, ErrNotFound
	}
	return lo.ToPtr(v.(T)), nil
}
