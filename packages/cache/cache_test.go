package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryCacheStore(t *testing.T) {
	t.Parallel()

	type sample struct{}

	expiration := 5 * time.Minute
	cleanup := 10 * time.Minute

	actual := NewInMemoryCacheStore[sample](expiration, cleanup)
	assert.Equal(t, expiration, actual.defaultExpiration)
	assert.Equal(t, cleanup, actual.cleanupExpiration)
}

func TestInMemoryCacheStore(t *testing.T) {
	t.Parallel()

	type sample struct{}
	expiration := 5 * time.Minute
	cleanup := 10 * time.Minute

	t.Run("success to get the value", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryCacheStore[sample](expiration, cleanup)
		store.Set("key", sample{})
		v, ok := store.Get("key")
		assert.True(t, ok)
		assert.NotNil(t, v)
	})

	t.Run("failed to get the value", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryCacheStore[sample](expiration, cleanup)
		store.Set("key", sample{})
		v, ok := store.Get("key2")
		assert.False(t, ok)
		assert.Nil(t, v)
	})
}
