package config

import "time"

// Config holds the configuration for the application.
type Config struct {
	Token         string        `env:"TUBU_DISCORD_TOKEN,required"`
	APITimeout    time.Duration `env:"TUBU_API_TIMEOUT" envDefault:"10s"`
	RedisAddr     string        `env:"TUBU_REDIS_ADDR"`
	RedisPassword string        `env:"TUBU_REDIS_PASSWORD"`
	RedisDB       int           `env:"TUBU_REDIS_DB" envDefault:"0"`
	RedisPoolSize int           `env:"TUBU_REDIS_POOL_SIZE" envDefault:"10"`
}
