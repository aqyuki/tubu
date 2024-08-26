package setup

import (
	"testing"
	"time"

	"github.com/aqyuki/tubu/packages/infra/redis"
	"github.com/stretchr/testify/assert"
)

func TestParseBotConfig(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		t.Run("全ての環境変数が指定されている場合", func(t *testing.T) {
			t.Setenv("TUBU_DISCORD_TOKEN", "token")
			t.Setenv("TUBU_API_TIMEOUT", "15s")

			except := &Config{
				Token:      "token",
				APITimeout: 15 * time.Second,
			}

			config, err := ParseBotConfig()
			assert.NoError(t, err, "expected no error but received %v", err)
			assert.Equal(t, except, config, "expected %v but received %v", except, config)
		})

		t.Run("TUBU_API_TIMEOUTが指定されていない場合，API_TIMEOUTにはデフォルト値が使用される", func(t *testing.T) {
			t.Setenv("TUBU_DISCORD_TOKEN", "token")

			except := &Config{
				Token:      "token",
				APITimeout: 10 * time.Second,
			}

			config, err := ParseBotConfig()
			assert.NoError(t, err, "expected no error but received %v", err)
			assert.Equal(t, except, config, "expected %v but received %v", except, config)
		})
	})

	t.Run("異常系", func(t *testing.T) {
		t.Run("必須の環境変数が指定されて居ない場合", func(t *testing.T) {
			config, err := ParseBotConfig()
			assert.Error(t, err, "expected error but received nil")
			assert.Nil(t, config, "expected nil but received %v", config)
		})
	})
}

func TestParseRedisConfig(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		t.Run("全ての環境変数がｓ指定されている場合", func(t *testing.T) {
			t.Setenv("TUBU_REDIS_ADDRESS", "localhost:6379")
			t.Setenv("TUBU_REDIS_PASSWORD", "password")
			t.Setenv("TUBU_REDIS_DB", "1")
			t.Setenv("TUBU_REDIS_POOL_SIZE", "20")

			except := &redis.RedisConfig{
				Address:  "localhost:6379",
				Password: "password",
				DB:       1,
				PoolSize: 20,
			}

			config, err := ParseRedisConfig()
			assert.NoError(t, err, "expected no error but received %v", err)
			assert.Equal(t, except, config, "expected %v but received %v", except, config)
		})

		t.Run("TUBU_REDIS_DBが指定されていない場合，REDIS_DBにはデフォルト値が使用される", func(t *testing.T) {
			t.Setenv("TUBU_REDIS_ADDRESS", "localhost:6379")
			t.Setenv("TUBU_REDIS_PASSWORD", "password")
			t.Setenv("TUBU_REDIS_POOL_SIZE", "20")

			except := &redis.RedisConfig{
				Address:  "localhost:6379",
				Password: "password",
				DB:       0,
				PoolSize: 20,
			}

			config, err := ParseRedisConfig()
			assert.NoError(t, err, "expected no error but received %v", err)
			assert.Equal(t, except, config, "expected %v but received %v", except, config)
		})

		t.Run("TUBU_REDIS_POOL_SIZEが指定されていない場合，REDIS_POOL_SIZEにはデフォルト値が使用される", func(t *testing.T) {
			t.Setenv("TUBU_REDIS_ADDRESS", "localhost:6379")
			t.Setenv("TUBU_REDIS_PASSWORD", "password")
			t.Setenv("TUBU_REDIS_DB", "1")

			except := &redis.RedisConfig{
				Address:  "localhost:6379",
				Password: "password",
				DB:       1,
				PoolSize: 10,
			}

			config, err := ParseRedisConfig()
			assert.NoError(t, err, "expected no error but received %v", err)
			assert.Equal(t, except, config, "expected %v but received %v", except, config)
		})
	})

	t.Run("異常系", func(t *testing.T) {
		t.Run("型が不一致の場合", func(t *testing.T) {
			t.Setenv("TUBU_REDIS_DB", "invalid")

			config, err := ParseRedisConfig()
			assert.Error(t, err, "expected error but received nil")
			assert.Nil(t, config, "expected nil but received %v", config)
		})
	})
}

func TestSetupCacheStore(t *testing.T) {
	t.Parallel()
	t.Run("Redisを使用する場合", func(t *testing.T) {
		t.Parallel()

		config := &redis.RedisConfig{
			Address:  "localhost:6379",
			Password: "password",
			DB:       1,
			PoolSize: 20,
		}
		store := SetupCacheStore[string](config)
		assert.NotNil(t, store, "expected not nil but received nil")
	})

	t.Run("InMemoryを使用する場合", func(t *testing.T) {
		t.Parallel()

		config := &redis.RedisConfig{}
		store := SetupCacheStore[string](config)
		assert.NotNil(t, store, "expected not nil but received nil")
	})

	t.Run("Redisの設定がnilの場合", func(t *testing.T) {
		t.Parallel()

		store := SetupCacheStore[string](nil)
		assert.NotNil(t, store, "expected not nil but received nil")
	})
}
