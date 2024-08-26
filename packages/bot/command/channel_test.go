package command

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_channelName(t *testing.T) {
	t.Parallel()
	ch := &discordgo.Channel{
		ID: "1234567890",
	}
	expected := &discordgo.MessageEmbedField{
		Name:   "チャンネル名",
		Value:  "<#1234567890>",
		Inline: true,
	}
	assert.Equal(t, expected, (&ChannelCommand{}).channelName(ch))
}

func Test_channelType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		channelType discordgo.ChannelType
	}{
		{name: "Text", channelType: discordgo.ChannelTypeGuildText},
		{name: "Voice", channelType: discordgo.ChannelTypeGuildVoice},
		{name: "Category", channelType: discordgo.ChannelTypeGuildCategory},
		{name: "Announce", channelType: discordgo.ChannelTypeGuildNews},
		{name: "Announce(Thread)", channelType: discordgo.ChannelTypeGuildNewsThread},
		{name: "Thread(Public)", channelType: discordgo.ChannelTypeGuildPublicThread},
		{name: "Thread(Private)", channelType: discordgo.ChannelTypeGuildPrivateThread},
		{name: "Stage", channelType: discordgo.ChannelTypeGuildStageVoice},
		{name: "Forum", channelType: discordgo.ChannelTypeGuildForum},
		{name: "Other", channelType: discordgo.ChannelTypeGuildStore},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ch := &discordgo.Channel{
				Type: tc.channelType,
			}
			expected := &discordgo.MessageEmbedField{
				Name:   "チャンネルタイプ",
				Value:  tc.name,
				Inline: true,
			}
			assert.Equal(t, expected, (&ChannelCommand{}).channelType(ch))
		})
	}
}

func Test_createdAt(t *testing.T) {
	t.Parallel()
	ch := &discordgo.Channel{
		ID: "0",
	}
	expected := &discordgo.MessageEmbedField{
		Name:   "作成日時",
		Value:  "<t:1420070400>",
		Inline: true,
	}
	assert.Equal(t, expected, (&ChannelCommand{}).createdAt(ch))
}
