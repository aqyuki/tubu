package command

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/bot/common"
	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

var _ discord.Command = (*GuildCommand)(nil)

type GuildCommand struct {
	cache cache.CacheStore[discordgo.Guild]
}

func NewGuildCommand(cache cache.CacheStore[discordgo.Guild]) *GuildCommand {
	return &GuildCommand{
		cache: cache,
	}
}

func (c *GuildCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "guild",
		Description: "ギルドの情報を表示します.",
	}
}

func (c *GuildCommand) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, s *discordgo.Session, ic *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("guild command is called")

		guild, err := c.cache.Get(ctx, ic.GuildID)
		if err != nil {
			g, err := s.Guild(ic.GuildID)
			if err != nil {
				logger.Error("Failed to get guild", err)
				return
			}
			if err := c.cache.Set(ctx, ic.GuildID, lo.FromPtr(g)); err != nil {
				logger.Errorf("failed to set the guild information to the cache: %v", err)
			}
			guild = g
		}

		embed := &discordgo.MessageEmbed{
			Title:       "拠点情報",
			Description: "このサーバーの情報だよ！",
			Color:       common.EmbedColor,
			Fields: []*discordgo.MessageEmbedField{
				c.guildName(guild),
				c.guildOwner(guild),
				c.afkChannel(guild),
				c.channelCount(guild),
				c.emojiCount(guild),
				c.roleCount(guild),
				c.stickerCount(guild),
				c.memberCount(guild),
				c.scale(guild),
				c.createdAt(guild),
			},
		}
		if err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Errorf("failed to respond: %v", err)
			return
		}
	}
}

func (c *GuildCommand) guildName(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "サーバー名",
		Value:  guild.Name,
		Inline: true,
	}
}

func (c *GuildCommand) guildOwner(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "オーナー",
		Value:  fmt.Sprintf("<@%s>", guild.OwnerID),
		Inline: true,
	}
}

func (c *GuildCommand) afkChannel(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	var afk string
	if guild.AfkChannelID == "" {
		afk = "なし"
	} else {
		afk = fmt.Sprintf("<#%s>(%d)", guild.AfkChannelID, guild.AfkTimeout)
	}

	return &discordgo.MessageEmbedField{
		Name:   "AFKチャンネル",
		Value:  afk,
		Inline: true,
	}
}

func (c *GuildCommand) channelCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "チャンネル数",
		Value:  fmt.Sprintf("%d", len(guild.Channels)),
		Inline: true,
	}
}

func (c *GuildCommand) emojiCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "絵文字数",
		Value:  fmt.Sprintf("%d", len(guild.Emojis)),
		Inline: true,
	}
}

func (c *GuildCommand) roleCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "ロール数",
		Value:  fmt.Sprintf("%d", len(guild.Roles)),
		Inline: true,
	}
}

func (c *GuildCommand) stickerCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "ステッカー数",
		Value:  fmt.Sprintf("%d", len(guild.Stickers)),
		Inline: true,
	}
}

func (c *GuildCommand) memberCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "メンバー数",
		Value:  fmt.Sprintf("%d", guild.MemberCount),
		Inline: true,
	}
}

func (c *GuildCommand) scale(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	var text string
	if guild.Large {
		text = "大規模"
	} else {
		text = "小規模"
	}
	return &discordgo.MessageEmbedField{
		Name:   "スケール",
		Value:  text,
		Inline: true,
	}
}

func (c *GuildCommand) createdAt(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "作成日時",
		Value:  fmt.Sprintf("<t:%d>", discord.TimestampFromSnowflake(guild.ID).Unix()),
		Inline: true,
	}
}
