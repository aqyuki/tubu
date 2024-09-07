package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var _ discord.MessageCreateHandler = (*CitationService)(nil).Citation

type CitationService struct {
	rgx   *regexp.Regexp
	cache cache.CacheStore[discordgo.Channel]
}

func NewCitationService(cache cache.CacheStore[discordgo.Channel]) *CitationService {
	return &CitationService{
		rgx:   regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord(app)?\.com/channels/(\d+)/(\d+)/(\d+)`),
		cache: cache,
	}
}

func (s *CitationService) Citation(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) {
	logger := logging.FromContext(ctx)
	if message.Author.Bot {
		logger.Debug("skip the processing because the message was created by the bot.")
		return
	}

	links := s.extractMessageLinks(message.Content)
	if len(links) == 0 {
		logger.Debug("skip the processing because there is no message link in the message.")
		return
	}

	ids, err := s.extractMessageInfo(links[0])
	if err != nil {
		logger.Error("failed to extract the id from the message link")
		return
	}

	if ids.guild != message.GuildID {
		logger.Debug("skip the processing because the message is not in the same guild.")
		return
	}

	channel, err := s.cache.Get(ctx, ids.channel)
	if err != nil {
		ch, err := session.Channel(ids.channel)
		if err != nil {
			logger.Error("failed to get the channel", zap.Error(err))
			return
		}
		if err := s.cache.Set(ctx, ids.channel, lo.FromPtr(ch)); err != nil {
			logger.Warn("failed to set the channel information to the cache", zap.Error(err))
		}
		channel = ch
	}
	if channel.NSFW {
		logger.Debug("skip the processing because the channel is NSFW.")
		return
	}

	msg, err := session.ChannelMessage(ids.channel, ids.message)
	if err != nil {
		logger.Error("failed to get the message", zap.Error(err))
		return
	}

	if msg.Content == "" {
		logger.Debug("skip the processing because the message content is empty.")
		return
	}

	var image *discordgo.MessageEmbedImage
	if len(msg.Attachments) > 0 {
		image = &discordgo.MessageEmbedImage{
			URL: msg.Attachments[0].URL,
		}
	}

	embed := &discordgo.MessageEmbed{
		Image: image,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    msg.Author.Username,
			IconURL: msg.Author.AvatarURL("64"),
		},
		Color:       EmbedColor,
		Description: msg.Content,
		Timestamp:   msg.Timestamp.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: channel.Name,
		},
	}

	replyMsg := discordgo.MessageSend{
		Embed:     embed,
		Reference: message.Reference(),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			RepliedUser: true,
		},
	}
	if _, err := session.ChannelMessageSendComplex(message.ChannelID, &replyMsg); err != nil {
		logger.Error("failed to send message", zap.Error(err))
		return
	}
	logger.Info("citation message has been sent", zap.String("channel_id", message.ChannelID), zap.String("message_id", message.ID))
}

func (s *CitationService) extractMessageLinks(text string) []string {
	return s.rgx.FindAllString(text, -1)
}

type message struct {
	guild   string
	channel string
	message string
}

// extractMessageInfo extracts the channel ID and message ID from the message link.
func (s *CitationService) extractMessageInfo(link string) (info message, err error) {
	segments := strings.Split(link, "/")
	if len(segments) < 4 {
		return message{}, errors.New("invalid message link")
	}
	return message{
		guild:   segments[len(segments)-3],
		channel: segments[len(segments)-2],
		message: segments[len(segments)-1],
	}, nil
}
