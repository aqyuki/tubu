package handler

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func ReadyHandler(md *metadata.Metadata) discord.ReadyHandler {
	return func(ctx context.Context, s *discordgo.Session, r *discordgo.Ready) {
		logger := logging.FromContext(ctx)
		logger.Infow(fmt.Sprintf("Bot is ready (username : %s)", r.User.Username), zap.Any("metadata", md))
		// Set the playing status.
		if err := s.UpdateCustomStatus(fmt.Sprintf("%sをプレイ中", md.Version)); err != nil {
			logger.Errorf("failed to set the playing status because of %v", err)
		}
	}
}
