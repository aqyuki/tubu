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
				Token:         viper.GetString("token"),
				Timeout:       viper.GetDuration("timeout"),
				RedisEnabled:  viper.GetBool("redis-enabled"),
				RedisAddr:     viper.GetString("redis-addr"),
				RedisPassword: viper.GetString("redis-password"),
				RedisDB:       viper.GetInt("redis-db"),
				RedisPoolSize: viper.GetInt("redis-pool-size"),
			}

			logger := logging.NewLoggerFromEnv()
			ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer done()
			ctx = logging.WithLogger(ctx, logger)

			if !conf.IsValid() {
				return config.ErrInvalidConfig
			}
			if conf.RedisEnabled && !conf.IsValidRedisConfig() {
				conf.RedisEnabled = false
				logger.Warnw("detected invalid Redis configuration. disabled Redis-cache and use in-memory cache", zap.Any("config", conf))
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
	viper.SetDefault("redis-enabled", false)
	viper.SetDefault("redis-db", 0)
	viper.SetDefault("redis-pool-size", 10)

	rootCmd.PersistentFlags().String("token", "", "token is a Discord bot token. It or TUBU_DISCORD_TOKEN is required.")
	rootCmd.PersistentFlags().Duration("timeout", 5*time.Second, "timeout is a duration for API requests.")
	rootCmd.PersistentFlags().Bool("redis-enabled", false, "redis-enabled is a flag to enable Redis to use cache.")
	rootCmd.PersistentFlags().String("redis-addr", "localhost:6379", "redis-addr is an address to connect to Redis.")
	rootCmd.PersistentFlags().String("redis-password", "", "redis-password is a password to connect to Redis.")
	rootCmd.PersistentFlags().Int("redis-db", 0, "redis-db is a database number to connect to Redis.")
	rootCmd.PersistentFlags().Int("redis-pool-size", 10, "redis-pool-size is a pool size of Redis.")

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
	if err := viper.BindEnv("redis-enabled", "REDIS_ENABLED"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("redis-addr", "REDIS_ADDR"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("redis-password", "REDIS_PASSWORD"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("redis-db", "REDIS_DB"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("redis-pool-size", "REDIS_POOL_SIZE"); err != nil {
		panic(err)
	}
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
