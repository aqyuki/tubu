package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
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

var _ CacheStore[string] = (*RedisCacheStore[string])(nil)

type RedisCacheStore[T any] struct {
	expiration time.Duration
	client     *redis.Client
}

func NewRedisCacheStore[T any](client *redis.Client, exp time.Duration) *RedisCacheStore[T] {
	return &RedisCacheStore[T]{
		expiration: exp,
		client:     client,
	}
}

func (s *RedisCacheStore[T]) Set(ctx context.Context, key string, value T) error {
	return s.client.Set(ctx, key, value, s.expiration).Err()
}

func (s *RedisCacheStore[T]) Get(ctx context.Context, key string) (*T, error) {
	v, err := s.client.Do(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the value: %w", err)
	}
	return lo.ToPtr(v.(T)), nil
}
