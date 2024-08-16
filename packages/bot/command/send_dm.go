package command

import (
	"context"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
)

var _ discord.Command = (*SendDMCommand)(nil)

type SendDMCommand struct{}

func NewSendDMCommand() *SendDMCommand {
	return &SendDMCommand{}
}

func (c *SendDMCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:              discordgo.MessageApplicationCommand,
		Name:              "Send DM",
		NameLocalizations: &map[discordgo.Locale]string{discordgo.Japanese: "DMに送る"},
	}
}

func (c *SendDMCommand) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, s *discordgo.Session, ic *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("send_dm command is called")

		message, ok := ic.ApplicationCommandData().Resolved.Messages[ic.ApplicationCommandData().TargetID]
		if !ok {
			logger.Errorf("failed to get message from resolved data")
			return
		}

		if message.Content == "" {
			logger.Info("message content is empty. skip")
			return
		}

		dm, err := s.UserChannelCreate(ic.Member.User.ID)
		if err != nil {
			logger.Error("failed to create DM channel", err)
			return
		}

		if _, err := s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
			Content: message.Content,
		}); err != nil {
			logger.Error("failed to send DM message", err)
			return
		}

		if err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "DMにピン留めしておいたよ！",
			},
		}); err != nil {
			logger.Error("failed to respond to interaction", err)
			return
		}
	}
}
