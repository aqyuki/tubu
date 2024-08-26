package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address  string `env:"TUBU_REDIS_ADDRESS"`
	Password string `env:"TUBU_REDIS_PASSWORD"`
	DB       int    `env:"TUBU_REDIS_DB" envDefault:"0"`
	PoolSize int    `env:"TUBU_REDIS_POOL_SIZE" envDefault:"10"`
}

func (c *RedisConfig) IsEnabled() bool {
	if c == nil {
		return false
	}
	return c.Address != "" && c.Password != ""
}

var ErrInvalidRedisOption = errors.New("redis option is invalid")

func NewRedisClient(cfg *RedisConfig) (*redis.Client, error) {
	if cfg == nil {
		return nil, ErrInvalidRedisOption
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	}), nil
}
