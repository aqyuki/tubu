package service

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type HealthService struct {
	md *metadata.Metadata
}

func NewHealthService(md *metadata.Metadata) *HealthService {
	return &HealthService{md: md}
}

func (s *HealthService) HealthCheck(ctx context.Context, session *discordgo.Session, ready *discordgo.Ready) {
	logger := logging.FromContext(ctx)
	logger.Info(fmt.Sprintf("Bot is ready (username : %s)", ready.User.Username), zap.Any("metadata", s.md))
	// Set the playing status.
	if err := session.UpdateCustomStatus(fmt.Sprintf("%sをプレイ中", s.md.Version)); err != nil {
		logger.Error("failed to set the playing status", zap.Error(err))
	}
}
