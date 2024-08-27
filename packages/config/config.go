package config

import (
	"errors"
	"time"
)

var (
	ErrInvalidConfig = errors.New("config: invalid config")
)

type Config struct {
	// Token is a Discord bot token. It is required.
	Token string

	// Timeout is a duration for Discord API requests.
	Timeout time.Duration

	// RedisEnabled is a flag to enable Redis to use cache.
	RedisEnabled bool

	// RedisAddr is an address to connect to Redis.
	RedisAddr string

	// RedisPassword is a password to connect to Redis.
	RedisPassword string

	// RedisDB is a database number to connect to Redis.
	RedisDB int

	// RedisPoolSize is a pool size of Redis.
	RedisPoolSize int
}

func (c Config) IsValid() bool {
	return c.Token != ""
}

func (c Config) IsValidRedisConfig() bool {
	return c.RedisAddr != "" && c.RedisDB >= 0 && c.RedisPoolSize > 0
}
