package discord

import "time"

const (
	minAPITimeout = 5 * time.Second
)

// Config is the configuration for the Discord bot.
type Config struct {
	APITimeout time.Duration
}

// DefaultConfig returns the default configuration for the Config.
func DefaultConfig() *Config {
	return &Config{
		APITimeout: 10 * time.Second,
	}
}

// Option is a configuration option for the Config.
type Option func(*Config)

// NewConfig creates a new Config with the given options.
func NewConfig(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithAPITimeout sets the timeout for API requests.
// If the timeout is less than the minimum allowed timeout, the minimum timeout is used.
func WithAPITimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		if timeout < minAPITimeout {
			timeout = minAPITimeout
		}
		cfg.APITimeout = timeout
	}
}
