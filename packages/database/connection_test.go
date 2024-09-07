package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	if testing.Short() {
		t.Skip("🚧 -shortフラグが指定されたためテストをスキップ")
	}

	dns := LoadConnectionString(t)
	db, err := NewDB(context.Background(), dns)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db.Pool.Close()
}
