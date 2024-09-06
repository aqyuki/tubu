package logging

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLoggerFromEnv(t *testing.T) {
	t.Run("LOG_MODEがdevelopの場合", func(t *testing.T) {
		t.Setenv("LOG_MODE", "develop")
		actual := NewLoggerFromEnv()

		assert.NotNil(t, actual)
	})

	t.Run("LOG_MODEがdevelopでない場合", func(t *testing.T) {
		t.Setenv("LOG_MODE", "production")
		actual := NewLoggerFromEnv()

		assert.NotNil(t, actual)
	})
}

func TestNewLogger(t *testing.T) {
	t.Parallel()
	t.Run("developがtrueの場合", func(t *testing.T) {
		t.Parallel()
		actual := NewLogger(true, "debug")

		assert.NotNil(t, actual)
	})

	t.Run("developがfalseの場合", func(t *testing.T) {
		t.Parallel()
		actual := NewLogger(false, "info")

		assert.NotNil(t, actual)
	})
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	actual1 := DefaultLogger()
	actual2 := DefaultLogger()

	assert.NotNil(t, actual1)
	assert.NotNil(t, actual2)
	assert.Equal(t, actual1, actual2)
}

func TestWithLogger(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	ctx := WithLogger(context.Background(), logger)

	assert.NotNil(t, ctx)
}

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("loggerがcontextに設定されている場合", func(t *testing.T) {
		t.Parallel()
		logger := zap.NewNop()
		ctx := WithLogger(context.Background(), logger)
		actual := FromContext(ctx)

		assert.NotNil(t, actual)
	})

	t.Run("loggerがcontextに設定されていない場合", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		actual := FromContext(ctx)

		assert.NotNil(t, actual)
	})
}

func Test_levelToZapLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level string
		want  zapcore.Level
	}{
		{
			name:  "debug",
			level: "debug",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "info",
			level: "info",
			want:  zapcore.InfoLevel,
		},
		{
			name:  "warning",
			level: "warning",
			want:  zapcore.WarnLevel,
		},
		{
			name:  "error",
			level: "error",
			want:  zapcore.ErrorLevel,
		},
		{
			name:  "critical",
			level: "critical",
			want:  zapcore.DPanicLevel,
		},
		{
			name:  "alert",
			level: "alert",
			want:  zapcore.PanicLevel,
		},
		{
			name:  "emergency",
			level: "emergency",
			want:  zapcore.FatalLevel,
		},
		{
			name:  "unknown",
			level: "unknown",
			want:  zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := levelToZapLevel(tt.level)

			assert.Equal(t, tt.want, actual)
		})
	}
}

func Test_levelEncoder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level zapcore.Level
		want  string
	}{
		{
			name:  "debug",
			level: zapcore.DebugLevel,
			want:  "DEBUG",
		},
		{
			name:  "info",
			level: zapcore.InfoLevel,
			want:  "INFO",
		},
		{
			name:  "warning",
			level: zapcore.WarnLevel,
			want:  "WARNING",
		},
		{
			name:  "error",
			level: zapcore.ErrorLevel,
			want:  "ERROR",
		},
		{
			name:  "critical",
			level: zapcore.DPanicLevel,
			want:  "CRITICAL",
		},
		{
			name:  "alert",
			level: zapcore.PanicLevel,
			want:  "ALERT",
		},
		{
			name:  "emergency",
			level: zapcore.FatalLevel,
			want:  "EMERGENCY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mem := zapcore.NewMapObjectEncoder()
			err := mem.AddArray("k", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
				levelEncoder()(tt.level, arr)
				return nil
			}))
			assert.NoError(t, err)
			arr := mem.Fields["k"].([]any)
			assert.Len(t, arr, 1)
			assert.Equal(t, tt.want, arr[0])
		})
	}
}

func Test_timeEncoder(t *testing.T) {
	t.Parallel()
	moment := time.Unix(100, 50005000).UTC()
	mem := zapcore.NewMapObjectEncoder()
	err := mem.AddArray("k", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
		timeEncoder()(moment, arr)
		return nil
	}))
	assert.NoError(t, err)
	arr := mem.Fields["k"].([]any)
	assert.Len(t, arr, 1)
	assert.Equal(t, "1970-01-01T00:01:40.050005Z", arr[0])
}
