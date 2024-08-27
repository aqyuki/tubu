package setup

import (
	"testing"

	"github.com/aqyuki/tubu/packages/config"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheStore(t *testing.T) {
	t.Parallel()

	t.Run("Redisが有効", func(t *testing.T) {
		t.Parallel()

		cnf := &config.Config{
			RedisEnabled:  true,
			RedisAddr:     "localhost:6379",
			RedisPassword: "password",
			RedisDB:       0,
			RedisPoolSize: 10,
		}

		cacheStore := NewCacheStore[interface{}](cnf)
		assert.NotNil(t, cacheStore)
	})

	t.Run("Redisが無効", func(t *testing.T) {
		t.Parallel()

		cnf := &config.Config{
			RedisEnabled: false,
		}

		cacheStore := NewCacheStore[interface{}](cnf)
		assert.NotNil(t, cacheStore)
	})
}
