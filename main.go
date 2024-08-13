package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aqyuki/tubu/packages/bot"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/samber/lo"
)

type exitCode int

const (
	ExitSuccess exitCode = iota
	ExitFailure
)

var (
	timeout = lo.FromPtr(flag.Duration("timeout", 10*time.Second, "API request timeout"))
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

	logger.Infof("try to load DISCORD_TOKEN")
	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		logger.Errorf("DISCORD_TOKEN is not set")
		return ExitFailure
	}
	logger.Infof("DISCORD_TOKEN was successfully loaded")

	md := metadata.GetMetadata()

	config := discord.NewConfig(
		discord.WithAPITimeout(timeout),
	)

	contextFunc := func() context.Context {
		return ctx
	}

	handler := discord.NewHandler(
		discord.WithContextFunc(contextFunc),
		discord.WithReadyHandler(bot.ReadyHandler(md)),
	)

	router := discord.NewCommandRouter(
		discord.WithCommandContextFunc(contextFunc),
	)

	discordBot := discord.NewBot(md, config, handler, router)
	if err := discordBot.Start(token); err != nil {
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
