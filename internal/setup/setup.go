package setup

import (
	"fmt"
	"time"

	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/infra/redis"
	"github.com/caarlos0/env/v11"
)

const (
	environmentPrefix = "TUBU_"
	defaultExpiration = 30 * time.Minute
	defaultCleanup    = 30 * time.Minute
)

type Config struct {
	Token      string        `env:"DISCORD_TOKEN,required"`
	APITimeout time.Duration `env:"API_TIMEOUT" envDefault:"10s"`
}

func ParseBotConfig() (*Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{Prefix: environmentPrefix})
	if err != nil {
		return nil, fmt.Errorf("failed to parse bot config: %w", err)
	}
	return &cfg, nil
}

func ParseRedisConfig() (*redis.RedisConfig, error) {
	cfg, err := env.ParseAsWithOptions[redis.RedisConfig](env.Options{Prefix: environmentPrefix})
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis config: %w", err)
	}
	return &cfg, nil
}

func SetupCacheStore[T any](cfg *redis.RedisConfig) cache.CacheStore[T] {
	if cfg.IsEnabled() {
		client, err := redis.NewRedisClient(cfg)
		if err == nil {
			return cache.NewRedisCacheStore[T](client, defaultExpiration)
		}
	}
	return cache.NewInMemoryCacheStore[T](defaultExpiration, defaultCleanup)
}
