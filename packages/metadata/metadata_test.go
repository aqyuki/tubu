package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetadata(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, GetMetadata())
}
