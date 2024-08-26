package discord

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampFromSnowflake(t *testing.T) {
	t.Parallel()
	expected := time.Unix(1420070400, 0)
	actual := TimestampFromSnowflake("0")
	assert.Equal(t, expected, actual)
}
