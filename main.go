package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aqyuki/tubu/internal/config"
	"github.com/aqyuki/tubu/packages/bot/command"
	"github.com/aqyuki/tubu/packages/bot/handler"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type exitCode int

const (
	ExitSuccess exitCode = iota
	ExitFailure
)

func main() {
	ctx := context.Background()
	ctx = logging.WithLogger(ctx, logging.NewLoggerFromEnv())
	exit(run(ctx))
}

func run(ctx context.Context) exitCode {
	ctx, done := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer done()
	logger := logging.FromContext(ctx)

	logger.Infof("try to load application config")
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		logger.Errorf("tried to load application config but failed with error: %v", err)
		return ExitFailure
	}
	logger.Infow("loaded application config", zap.Any("config", cfg))

	md := metadata.GetMetadata()

	config := discord.NewConfig(
		discord.WithAPITimeout(cfg.APITimeout),
	)

	contextFunc := func() context.Context {
		return ctx
	}

	handler := discord.NewHandler(
		discord.WithContextFunc(contextFunc),
		discord.WithReadyHandler(handler.ReadyHandler(md)),
	)

	router := discord.NewCommandRouter(
		discord.WithCommandContextFunc(contextFunc),
		discord.WithCommand(command.NewVersionCommand(md)),
	)

	discordBot := discord.NewBot(md, config, handler, router)
	if err := discordBot.Start(cfg.Token); err != nil {
		logger.Errorf("tried to start discord bot but failed with error: %v", err)
		return ExitFailure
	}

	<-ctx.Done()
	logger.Infof("received signal to shutdown")
	if err := discordBot.Shutdown(); err != nil {
		logger.Errorf("tried to stop discord bot but failed with error: %v", err)
		return ExitFailure
	}
	return ExitSuccess
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}
