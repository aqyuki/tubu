package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValid(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		cnf := &Config{Token: "token"}
		assert.True(t, cnf.IsValid())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		cnf := &Config{}
		assert.False(t, cnf.IsValid())
	})
}

func Test_IsValidRedisConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		cnf := &Config{
			RedisAddr:     "localhost:6379",
			RedisDB:       0,
			RedisPoolSize: 10,
		}
		assert.True(t, cnf.IsValidRedisConfig())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		cnf := &Config{
			RedisAddr:     "localhost:6379",
			RedisDB:       -1,
			RedisPoolSize: 0,
		}
		assert.False(t, cnf.IsValidRedisConfig())
	})
}
