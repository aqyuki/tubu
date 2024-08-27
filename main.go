package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aqyuki/tubu/internal/setup"
	"github.com/aqyuki/tubu/packages/bot/command"
	"github.com/aqyuki/tubu/packages/bot/handler"
	"github.com/aqyuki/tubu/packages/config"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tubu",
		Short: "tubu is a discord bot",
		RunE: func(_ *cobra.Command, _ []string) error {
			conf := config.Config{
				Token:   viper.GetString("token"),
				Timeout: viper.GetDuration("timeout"),
			}

			logger := logging.NewLoggerFromEnv()
			ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer done()
			ctx = logging.WithLogger(ctx, logger)

			if !conf.IsValid() {
				return config.ErrInvalidConfig
			}
			logger.Info("bot configuration was loaded successfully")

			md := metadata.GetMetadata()

			config := discord.NewConfig(
				discord.WithAPITimeout(conf.Timeout),
			)
			handler := discord.NewHandler(
				discord.WithHandlerContextFunc(buildContextFunc(ctx)),
				discord.WithReadyHandler(handler.ReadyHandler(md)),
				discord.WithMessageCreateHandler(handler.NewExpandHandler(setup.NewCacheStore[discordgo.Channel](&conf)).Expand),
			)
			router := discord.NewCommandRouter(
				discord.WithCommandContextFunc(buildContextFunc(ctx)),
				discord.WithCommand(command.NewVersionCommand(md)),
				discord.WithCommand(command.NewDiceCommand()),
				discord.WithCommand(command.NewChannelCommand()),
				discord.WithCommand(command.NewGuildCommand(setup.NewCacheStore[discordgo.Guild](&conf))),
				discord.WithCommand(command.NewSendDMCommand()),
			)

			bot := discord.NewBot(md, config, handler, router)
			if err := bot.Start(conf.Token); err != nil {
				logger.Errorw("failed to start bot", zap.Error(err))
				return err
			}

			<-ctx.Done()
			logger.Info("received signal to stop bot")

			if err := bot.Shutdown(); err != nil {
				logger.Errorw("failed to stop bot", zap.Error(err))
				return err
			}
			return nil
		},
	}
)

func init() {
	viper.SetDefault("timeout", "10s")

	rootCmd.PersistentFlags().String("token", "", "token is a Discord bot token. It or TUBU_DISCORD_TOKEN is required.")
	rootCmd.PersistentFlags().Duration("timeout", 5*time.Second, "timeout is a duration for API requests.")

	if err := viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("redis-enabled", rootCmd.PersistentFlags().Lookup("redis-enabled")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("redis-addr", rootCmd.PersistentFlags().Lookup("redis-addr")); err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("tubu")
	viper.AutomaticEnv()
}

func buildContextFunc(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}

func main() {
	// 発生したエラーによって終了コードを変更したいため，正常終了した場合を先に処理する．
	if err := rootCmd.Execute(); err == nil {
		os.Exit(0)
	}
	os.Exit(1)
}
