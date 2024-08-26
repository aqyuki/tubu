package handler

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/aqyuki/tubu/packages/bot/common"
	"github.com/aqyuki/tubu/packages/cache"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ExpandHandler struct {
	rgx   *regexp.Regexp
	cache cache.CacheStore[discordgo.Channel]
}

func NewExpandHandler(cache cache.CacheStore[discordgo.Channel]) *ExpandHandler {
	return &ExpandHandler{
		rgx:   regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord(app)?\.com/channels/(\d+)/(\d+)/(\d+)`),
		cache: cache,
	}
}

// ExpandHandler is a handler that expands the link in the message.
func (h *ExpandHandler) Expand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := logging.FromContext(ctx)
	if m.Author.Bot {
		logger.Infof("skip the processing because the message was created by the bot.")
		return
	}

	links := h.extractMessageLinks(m.Content)
	if len(links) == 0 {
		logger.Infof("skip the processing because there is no message link in the message.")
		return
	}

	ids, err := h.extractMessageInfo(links[0])
	if err != nil {
		logger.Error("failed to extract the id from the message link")
		return
	}

	if ids.guild != m.GuildID {
		logger.Infof("skip the processing because the message is not in the same guild.")
		return
	}

	channel, err := h.cache.Get(ctx, ids.channel)
	if err != nil {
		ch, err := s.Channel(ids.channel)
		if err != nil {
			logger.Error("failed to get the channel")
			return
		}
		if err := h.cache.Set(ctx, ids.channel, lo.FromPtr(ch)); err != nil {
			logger.Errorf("failed to set the channel information to the cache: %v", err)
		}
		channel = ch
	}
	if channel.NSFW {
		logger.Infof("skip the processing because the channel is NSFW.")
		return
	}

	msg, err := s.ChannelMessage(ids.channel, ids.message)
	if err != nil {
		logger.Error("failed to get the message")
		return
	}

	if msg.Content == "" {
		logger.Info("skip the processing because the message content is empty.")
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
		Color:       common.EmbedColor,
		Description: msg.Content,
		Timestamp:   msg.Timestamp.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: channel.Name,
		},
	}

	replyMsg := discordgo.MessageSend{
		Embed:     embed,
		Reference: m.Reference(),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			RepliedUser: true,
		},
	}
	if _, err := s.ChannelMessageSendComplex(m.ChannelID, &replyMsg); err != nil {
		logger.Error("failed to send message", zap.Error(err))
		return
	}
	logger.Info("message link expanded")
}

func (h *ExpandHandler) extractMessageLinks(s string) []string {
	return h.rgx.FindAllString(s, -1)
}

type message struct {
	guild   string
	channel string
	message string
}

// extractMessageInfo extracts the channel ID and message ID from the message link.
func (h *ExpandHandler) extractMessageInfo(link string) (info message, err error) {
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
