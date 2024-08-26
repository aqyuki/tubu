package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEnabled(t *testing.T) {
	t.Parallel()

	t.Run("AddressとPasswordが指定されている場合", func(t *testing.T) {
		t.Parallel()
		cfg := &RedisConfig{
			Address:  "localhost:6379",
			Password: "password",
		}
		assert.True(t, cfg.IsEnabled())
	})

	t.Run("Addressが指定されていない場合", func(t *testing.T) {
		t.Parallel()
		cfg := &RedisConfig{
			Password: "password",
		}
		assert.False(t, cfg.IsEnabled())
	})

	t.Run("Passwordが指定されていない場合", func(t *testing.T) {
		t.Parallel()
		cfg := &RedisConfig{
			Address: "localhost:6379",
		}
		assert.False(t, cfg.IsEnabled())
	})

	t.Run("AddressとPasswordが指定されていない場合", func(t *testing.T) {
		t.Parallel()
		cfg := &RedisConfig{}
		assert.False(t, cfg.IsEnabled())
	})

	t.Run("nilの場合", func(t *testing.T) {
		t.Parallel()
		var cfg *RedisConfig
		assert.False(t, cfg.IsEnabled())
	})
}

func TestNewRedisClient(t *testing.T) {
	t.Parallel()

	t.Run("RedisConfigがnilの場合", func(t *testing.T) {
		t.Parallel()
		client, err := NewRedisClient(nil)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		cfg := &RedisConfig{
			Address:  "localhost:6379",
			Password: "password",
			DB:       1,
			PoolSize: 20,
		}
		client, err := NewRedisClient(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
