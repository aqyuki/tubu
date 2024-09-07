package service

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ discord.Command = (*ChannelInformationService)(nil)

const channelCommandChannelOptionName = "channel"

type ChannelInformationService struct{}

func NewChannelInformationService() *ChannelInformationService { return &ChannelInformationService{} }

func (s *ChannelInformationService) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "channel",
		Description: "チャンネルの情報を表示します.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        channelCommandChannelOptionName,
				Description: "情報を表示するチャンネルを指定します.",
				Required:    true,
			},
		},
	}
}
func (s *ChannelInformationService) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, session *discordgo.Session, ic *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("Channel command is called")

		options := ic.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
		for _, option := range options {
			optionMap[option.Name] = option
		}

		channelOption, ok := optionMap[channelCommandChannelOptionName]
		if !ok {
			logger.Error("channel option is not found")
			return
		}
		channel := channelOption.ChannelValue(session)

		embed := &discordgo.MessageEmbed{
			Title:       "チャンネルの情報",
			Description: "頼まれていたチャンネルの情報だよ！",
			Color:       EmbedColor,
			Fields: []*discordgo.MessageEmbedField{
				s.channelName(channel),
				s.channelType(channel),
				s.createdAt(channel),
			},
		}

		if err := session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond", zap.Error(err))
			return
		}
	}
}

func (s *ChannelInformationService) channelName(ch *discordgo.Channel) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "チャンネル名",
		Value:  fmt.Sprintf("<#%s>", ch.ID),
		Inline: true,
	}
}

func (s *ChannelInformationService) channelType(ch *discordgo.Channel) *discordgo.MessageEmbedField {
	var channelType string
	switch ch.Type {
	case discordgo.ChannelTypeGuildText:
		channelType = "Text"
	case discordgo.ChannelTypeGuildVoice:
		channelType = "Voice"
	case discordgo.ChannelTypeGuildCategory:
		channelType = "Category"
	case discordgo.ChannelTypeGuildNews:
		channelType = "Announce"
	case discordgo.ChannelTypeGuildNewsThread:
		channelType = "Announce(Thread)"
	case discordgo.ChannelTypeGuildPublicThread:
		channelType = "Thread(Public)"
	case discordgo.ChannelTypeGuildPrivateThread:
		channelType = "Thread(Private)"
	case discordgo.ChannelTypeGuildStageVoice:
		channelType = "Stage"
	case discordgo.ChannelTypeGuildForum:
		channelType = "Forum"
	default:
		channelType = "Other"
	}
	return &discordgo.MessageEmbedField{
		Name:   "チャンネルタイプ",
		Value:  channelType,
		Inline: true,
	}
}

func (s *ChannelInformationService) createdAt(ch *discordgo.Channel) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "作成日時",
		Value:  fmt.Sprintf("<t:%d>", discord.TimestampFromSnowflake(ch.ID).Unix()),
		Inline: true,
	}
}
