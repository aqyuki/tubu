package config

import "time"

// Config holds the configuration for the application.
type Config struct {
	Token      string        `env:"TUBU_DISCORD_TOKEN,required"`
	APITimeout time.Duration `env:"TUBU_API_TIMEOUT" envDefault:"10s"`
}
