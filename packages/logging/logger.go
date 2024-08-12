package logging

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// defaultLogger holds a logger used to default logger.
	// If you want to get a default logger. You must get default logger from DefaultLogger function.
	// because, defaultLogger initialized when first call DefaultLogger function.
	defaultLogger *zap.SugaredLogger

	// defaultLoggerOnce is a sync.Once variable to create default logger once.
	defaultLoggerOnce sync.Once
)

// NewLoggerFromEnv creates a logger with configuration from environment variables.
// If not set environment variables, it will return a logger with production mode and info level.
func NewLoggerFromEnv() *zap.SugaredLogger {
	// develop is a flag variable to switch logger mode between develop mode or not develop mode.
	// default is not develop mode.
	// TODO: change to get environment variable from config provider.
	develop := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "develop"

	// level is a log level variable to set log level.
	// TODO: change to get environment variable from config provider.
	level := os.Getenv("LOG_LEVEL")

	return NewLogger(develop, level)
}

// NewLogger creates a logger with given configuration.
// If not parse level argument, it will return a logger with info level.
func NewLogger(develop bool, level string) *zap.SugaredLogger {
	// config is a configuration to use base to create logger.
	var config zap.Config

	// TODO: customize logger configuration. but now, it is just a simple configuration.
	if develop {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Level = zap.NewAtomicLevelAt(stringToZapLevel(level))

	logger, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	return logger.Sugar()
}

// DefaultLogger returns a logger from configuration based on environment variables.
// If not created default logger, it will creates  a new logger and set it to default logger.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

// stringToZapLevel convert given string to zap level.
// If not match, it will return info level.
func stringToZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// contextKey is a private type used to define context key.
type contextKey string

// loggerKey is a context key to store logger in context.
const loggerKey = contextKey("logger")

// WithLogger stores a given logger to given context.
// If context is nil, it will panic. because, this function is wrapper for context.WithValue.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns a logger from given context.
// If not contained logger from given context, it will return a default logger.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}
