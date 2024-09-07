package service

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ discord.Command = (*VersionService)(nil)

type VersionService struct {
	md *metadata.Metadata
}

func NewVersionService(md *metadata.Metadata) *VersionService {
	return &VersionService{md: md}
}

func (s *VersionService) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "version",
		Description: "Botのバージョンを確認します",
	}
}

func (s *VersionService) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, ses *discordgo.Session, i *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("version command is called")

		embed := &discordgo.MessageEmbed{
			Title:       "現在のバージョン",
			Color:       EmbedColor,
			Description: fmt.Sprintf("現在のバージョンは `%s` です", s.md.Version),
		}
		if err := ses.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond to the interaction", zap.Error(err))
		}
	}
}
