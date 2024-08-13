package discord

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/bwmarrin/discordgo"
)

var (
	ErrAlreadyStarted = errors.New("tried to create a new session to activate the bot, but the session had already been created")
	ErrNoRunningBot   = errors.New("tried to shut down the bot, but the bot was not running")
)

type Bot struct {
	session  *discordgo.Session
	metadata *metadata.Metadata
	config   *Config
}

func NewBot(md *metadata.Metadata, cfg *Config) *Bot {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &Bot{
		session:  nil, // Session is initialized at bot startup
		metadata: md,
		config:   cfg,
	}
}

// Start starts the bot.
func (b *Bot) Start(token string) error {
	if b.session != nil {
		return ErrAlreadyStarted
	}
	ses, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("tried to create a new session to activate the bot, but failed to create a new session with error: %w", err)
	}
	b.session = ses

	client := &http.Client{
		Timeout: b.config.APITimeout,
	}
	b.session.Client = client

	// TODO: Add event handlers
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("tried to open a session to activate the bot, but failed to open the session with error: %w", err)
	}
	// TODO: register commands
	return nil
}

// Shutdown shuts down the bot.
func (b *Bot) Shutdown() error {
	if b.session == nil {
		return ErrNoRunningBot
	}

	// TODO: Remove event handlers
	// TODO: Unregister commands

	if err := b.session.Close(); err != nil {
		return fmt.Errorf("tried to close the session to shut down the bot, but failed to close the session with error: %w", err)
	}
	return nil
}
