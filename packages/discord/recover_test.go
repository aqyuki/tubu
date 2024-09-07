package discord

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_recoveryPanic(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		defer recoveryPanic(context.Background(), "test")
		panic("test")
	})
}
