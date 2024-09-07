package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aqyuki/tubu/internal/setup"
	"github.com/aqyuki/tubu/packages/discord"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/profile"
	"github.com/aqyuki/tubu/packages/service"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrInvalidConfig = errors.New("config: invalid config")

	rootCmd = &cobra.Command{
		Use:   "tubu",
		Short: "tubu is a discord bot",
		RunE: func(_ *cobra.Command, _ []string) error {
			prof := profile.Profile{
				Token:   viper.GetString("token"),
				Timeout: viper.GetDuration("timeout"),
			}

			logger := logging.NewLoggerFromEnv()
			ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer done()
			ctx = logging.WithLogger(ctx, logger)

			if !prof.IsValid() {
				return ErrInvalidConfig
			}
			logger.Info("bot configuration was loaded successfully")

			md := metadata.GetMetadata()

			config := discord.NewConfig(
				discord.WithAPITimeout(prof.Timeout),
			)
			handler := discord.NewHandler(
				discord.WithHandlerContextFunc(buildContextFunc(ctx)),
				discord.WithReadyHandler(service.NewHealthService(md).HealthCheck),
				discord.WithMessageCreateHandler(service.NewCitationService(setup.NewCacheStore[discordgo.Channel](&prof)).Citation),
			)
			router := discord.NewCommandRouter(
				discord.WithCommandContextFunc(buildContextFunc(ctx)),
				discord.WithCommand(service.NewVersionService(md)),
				discord.WithCommand(service.NewDiceService()),
				discord.WithCommand(service.NewChannelInformationService()),
				discord.WithCommand(service.NewGuildInformationService(setup.NewCacheStore[discordgo.Guild](&prof))),
				discord.WithCommand(service.NewSendDMService()),
			)

			bot := discord.NewBot(config, handler, router)
			if err := bot.Start(prof.Token); err != nil {
				logger.Error("failed to start bot", zap.Error(err))
				return err
			}

			<-ctx.Done()
			logger.Info("received signal to stop bot")

			if err := bot.Shutdown(); err != nil {
				logger.Error("failed to stop bot", zap.Error(err))
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
