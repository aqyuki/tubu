package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aqyuki/tubu/internal/setup"
	"github.com/aqyuki/tubu/packages/bot/command"
	"github.com/aqyuki/tubu/packages/bot/handler"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
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

	// load configurations
	cfg, err := setup.ParseBotConfig()
	if err != nil {
		logger.Errorf("failed to parse bot config: %v", err)
		return ExitFailure
	}

	redisConfig, err := setup.ParseRedisConfig()
	if err != nil {
		logger.Warnf("failed to parse redis config: %v", err)
	}
	logger.Info("loaded bot configurations successfully")

	// initialize discord bot
	md := metadata.GetMetadata()
	channelCache := setup.SetupCacheStore[discordgo.Channel](redisConfig)
	guildCache := setup.SetupCacheStore[discordgo.Guild](redisConfig)

	config := discord.NewConfig(
		discord.WithAPITimeout(cfg.APITimeout),
	)
	handler := discord.NewHandler(
		discord.WithHandlerContextFunc(BuildContextFunc(ctx)),
		discord.WithReadyHandler(handler.ReadyHandler(md)),
		discord.WithMessageCreateHandler(handler.NewExpandHandler(channelCache).Expand),
	)
	router := discord.NewCommandRouter(
		discord.WithCommandContextFunc(BuildContextFunc(ctx)),
		discord.WithCommand(command.NewVersionCommand(md)),
		discord.WithCommand(command.NewDiceCommand()),
		discord.WithCommand(command.NewChannelCommand()),
		discord.WithCommand(command.NewGuildCommand(guildCache)),
		discord.WithCommand(command.NewSendDMCommand()),
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

func BuildContextFunc(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}
