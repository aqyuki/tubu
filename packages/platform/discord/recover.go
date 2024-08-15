package discord

import (
	"context"

	"github.com/aqyuki/tubu/packages/logging"
	"go.uber.org/zap"
)

func recoveryPanic(ctx context.Context, a any) {
	if r := recover(); r != nil {
		logger := logging.FromContext(ctx)
		logger.Errorw("recovered from panic", zap.Any("panic", r), zap.Any("detail", a))
	}
}
