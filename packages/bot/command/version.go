package command

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/metadata"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
)

var _ discord.Command = (*VersionCommand)(nil)

type VersionCommand struct {
	md *metadata.Metadata
}

func NewVersionCommand(md *metadata.Metadata) *VersionCommand {
	return &VersionCommand{md: md}
}

func (c *VersionCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "version",
		Description: "Botのバージョンを表示します",
	}
}

func (c *VersionCommand) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("version command is called")

		embed := &discordgo.MessageEmbed{
			Title:       "現在のバージョン",
			Color:       EmbedColor,
			Description: fmt.Sprintf("現在のバージョンは `%s` です", c.md.Version),
		}
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond to the interaction", "error", err)
		}
	}
}
