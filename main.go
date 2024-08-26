package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aqyuki/tubu/internal/config"
	"github.com/aqyuki/tubu/packages/bot/command"
	"github.com/aqyuki/tubu/packages/bot/handler"
	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v11"
	"github.com/redis/go-redis/v9"
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

	// load application config
	logger.Infof("try to load application config")
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		logger.Errorf("tried to load application config but failed with error: %v", err)
		return ExitFailure
	}
	logger.Infow("loaded application config", zap.Any("config", cfg))

	// initialize cache store
	var channelCache cache.CacheStore[discordgo.Channel]
	var guildCache cache.CacheStore[discordgo.Guild]

	// initialize cache store
	if cfg.IsRedisEnabled() {
		redisClient := redis.NewClient(buildRedisOption(&cfg))
		channelCache = cache.NewRedisCacheStore[discordgo.Channel](redisClient, 30*time.Minute)
		guildCache = cache.NewRedisCacheStore[discordgo.Guild](redisClient, 30*time.Minute)
	} else {
		channelCache = cache.NewInMemoryCacheStore[discordgo.Channel](5*time.Minute, 30*time.Minute)
		guildCache = cache.NewInMemoryCacheStore[discordgo.Guild](5*time.Minute, 30*time.Minute)
	}

	// initialize discord bot
	md := metadata.GetMetadata()
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

func buildRedisOption(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: cfg.RedisPoolSize,
	}
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}
