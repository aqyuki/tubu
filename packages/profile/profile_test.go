package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValid(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		cnf := &Profile{Token: "token"}
		assert.True(t, cnf.IsValid())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		cnf := &Profile{}
		assert.False(t, cnf.IsValid())
	})
}
