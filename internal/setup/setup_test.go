package setup

import (
	"testing"
	"time"

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
