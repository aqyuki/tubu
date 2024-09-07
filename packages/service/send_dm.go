package service

import (
	"context"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ discord.Command = (*SendDMService)(nil)

type SendDMService struct{}

func NewSendDMService() *SendDMService {
	return &SendDMService{}
}

func (s *SendDMService) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:              discordgo.MessageApplicationCommand,
		Name:              "Send DM",
		NameLocalizations: &map[discordgo.Locale]string{discordgo.Japanese: "DMに送る"},
	}
}

func (s *SendDMService) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, session *discordgo.Session, ic *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("send_dm command is called")

		message, ok := ic.ApplicationCommandData().Resolved.Messages[ic.ApplicationCommandData().TargetID]
		if !ok {
			logger.Error("failed to get message from resolved data")
			return
		}

		if message.Content == "" {
			logger.Info("message content is empty. skip")
			return
		}

		dm, err := session.UserChannelCreate(ic.Member.User.ID)
		if err != nil {
			logger.Error("failed to create DM channel", zap.Error(err))
			return
		}

		if _, err := session.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
			Content: message.Content,
		}); err != nil {
			logger.Error("failed to send DM message", zap.Error(err))
			return
		}

		if err := session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "DMにピン留めしておきました．",
			},
		}); err != nil {
			logger.Error("failed to respond to interaction", zap.Error(err))
			return
		}
	}
}
