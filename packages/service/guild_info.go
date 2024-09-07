package service

import (
	"context"
	"fmt"

	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var _ discord.Command = (*GuildInformationService)(nil)

type GuildInformationService struct {
	cache cache.CacheStore[discordgo.Guild]
}

func NewGuildInformationService(cache cache.CacheStore[discordgo.Guild]) *GuildInformationService {
	return &GuildInformationService{
		cache: cache,
	}
}

func (s *GuildInformationService) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "guild",
		Description: "ギルドについて確認します．",
	}
}

func (s *GuildInformationService) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, session *discordgo.Session, ic *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("guild command is called")

		guild, err := s.cache.Get(ctx, ic.GuildID)
		if err != nil {
			g, err := session.Guild(ic.GuildID)
			if err != nil {
				logger.Error("Failed to get guild", zap.Error(err))
				return
			}
			if err := s.cache.Set(ctx, ic.GuildID, lo.FromPtr(g)); err != nil {
				logger.Error("failed to set the guild information to the cache", zap.Error(err))
			}
			guild = g
		}

		embed := &discordgo.MessageEmbed{
			Title:       "ギルドについて",
			Description: "このサーバーについてです．",
			Color:       EmbedColor,
			Fields: []*discordgo.MessageEmbedField{
				s.guildName(guild),
				s.guildOwner(guild),
				s.afkChannel(guild),
				s.channelCount(guild),
				s.emojiCount(guild),
				s.roleCount(guild),
				s.stickerCount(guild),
				s.memberCount(guild),
				s.scale(guild),
				s.createdAt(guild),
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
		logger.Info("guild information is sent")
	}
}

func (s *GuildInformationService) guildName(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "サーバー名",
		Value:  guild.Name,
		Inline: true,
	}
}

func (s *GuildInformationService) guildOwner(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "オーナー",
		Value:  fmt.Sprintf("<@%s>", guild.OwnerID),
		Inline: true,
	}
}

func (s *GuildInformationService) afkChannel(guild *discordgo.Guild) *discordgo.MessageEmbedField {
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

func (s *GuildInformationService) channelCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "チャンネル数",
		Value:  fmt.Sprintf("%d", len(guild.Channels)),
		Inline: true,
	}
}

func (s *GuildInformationService) emojiCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "絵文字数",
		Value:  fmt.Sprintf("%d", len(guild.Emojis)),
		Inline: true,
	}
}

func (s *GuildInformationService) roleCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "ロール数",
		Value:  fmt.Sprintf("%d", len(guild.Roles)),
		Inline: true,
	}
}

func (s *GuildInformationService) stickerCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "ステッカー数",
		Value:  fmt.Sprintf("%d", len(guild.Stickers)),
		Inline: true,
	}
}

func (s *GuildInformationService) memberCount(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "メンバー数",
		Value:  fmt.Sprintf("%d", guild.MemberCount),
		Inline: true,
	}
}

func (s *GuildInformationService) scale(guild *discordgo.Guild) *discordgo.MessageEmbedField {
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

func (s *GuildInformationService) createdAt(guild *discordgo.Guild) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   "作成日時",
		Value:  fmt.Sprintf("<t:%d>", discord.TimestampFromSnowflake(guild.ID).Unix()),
		Inline: true,
	}
}
