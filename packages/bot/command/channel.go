package command

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aqyuki/tubu/packages/bot/common"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
)

var _ discord.Command = (*ChannelCommand)(nil)

const (
	channelCommandChannelOptionName = "channel"

	// Discord Epochは，Discord APIのSnowflakeで使用される基準日時
	// 基準時刻は，2015/01/01 00:00:00 UTC
	// 上位42ビットはタイムスタンプ，次の5ビットはデータセンターID，次の5ビットはワーカーID，最後の12ビットはシーケンス番号
	// ref: https://discord.com/developers/docs/reference#snowflakes
	discordEpoch = int64(1420070400000)
)

type ChannelCommand struct{}

func NewChannelCommand() *ChannelCommand {
	return &ChannelCommand{}
}

func (c *ChannelCommand) Command() *discordgo.ApplicationCommand {
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

func (c *ChannelCommand) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, s *discordgo.Session, ic *discordgo.InteractionCreate) {
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
		channel := channelOption.ChannelValue(s)

		embed := &discordgo.MessageEmbed{
			Title:       "チャンネルの情報",
			Description: "頼まれていたチャンネルの情報だよ！",
			Color:       common.EmbedColor,
			Fields: []*discordgo.MessageEmbedField{
				c.channelName(channel),
				c.channelType(channel),
				c.createdAt(channel),
			},
		}

		if err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond", "error", err)
			return
		}
	}
}

func (c *ChannelCommand) channelName(ch *discordgo.Channel) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "`チャンネル名`",
		Value:  fmt.Sprintf("<#%s>", ch.ID),
		Inline: true,
	}
}

func (c *ChannelCommand) channelType(ch *discordgo.Channel) *discordgo.MessageEmbedField {
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
		Name:   "`チャンネルタイプ`",
		Value:  fmt.Sprintf("`%s`", channelType),
		Inline: true,
	}
}

func (c *ChannelCommand) createdAt(ch *discordgo.Channel) *discordgo.MessageEmbedField {
	// IDは，確実にsnowflakeで有るため，簡略化の為にエラーチェックを省略
	snowflake, _ := strconv.ParseInt(ch.ID, 10, 64)
	createdAt := time.Unix(0, ((snowflake>>22)+discordEpoch)*int64(time.Millisecond))
	return &discordgo.MessageEmbedField{
		Name:   "`作成日時`",
		Value:  fmt.Sprintf("`<t:%d>`", createdAt.Unix()),
		Inline: true,
	}
}
