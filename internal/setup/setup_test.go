package setup

import (
	"testing"

	"github.com/aqyuki/tubu/packages/profile"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheStore(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		cnf := &profile.Profile{}

		cacheStore := NewCacheStore[interface{}](cnf)
		assert.NotNil(t, cacheStore)
	})
}
