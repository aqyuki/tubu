package command

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/aqyuki/tubu/packages/bot/common"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var _ discord.Command = (*DiceCommand)(nil)

const (
	diceCommandCountOptionName = "count"
	diceCommandFaceOptionName  = "face"

	diceCommandMinCount = 1
	diceCommandMinFace  = 2
	diceCommandMaxCount = 10
	diceCommandMaxFace  = 100
)

type DiceCommand struct{}

func NewDiceCommand() *DiceCommand {
	return &DiceCommand{}
}

func (c *DiceCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "dice",
		Description: "サイコロを振ります.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        diceCommandCountOptionName,
				Description: "振るダイスの個数を指定します.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        diceCommandFaceOptionName,
				Description: "振るダイスの面数を指定します.",
				Required:    true,
			},
		},
	}
}

func (c *DiceCommand) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, s *discordgo.Session, ic *discordgo.InteractionCreate) {
		logger := logging.FromContext(ctx)
		logger.Debug("dice command is called")

		options := ic.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
		for _, option := range options {
			optionMap[option.Name] = option
		}

		countOpt, ok := optionMap[diceCommandCountOptionName]
		if !ok {
			logger.Error("option is not found", zap.String("name", diceCommandCountOptionName))
			s.InteractionRespond(ic.Interaction, c.errorResponse(diceCommandCountOptionName))
			return
		}
		faceOpt, ok := optionMap[diceCommandFaceOptionName]
		if !ok {
			logger.Error("option is not found", zap.String("name", diceCommandFaceOptionName))
			s.InteractionRespond(ic.Interaction, c.errorResponse(diceCommandFaceOptionName))
			return
		}

		count := countOpt.IntValue()
		face := faceOpt.IntValue()

		if count < diceCommandMinCount {
			count = diceCommandMinCount
		} else if count > diceCommandMaxCount {
			count = diceCommandMaxCount
		}

		if face < diceCommandMinFace {
			face = diceCommandMinFace
		} else if face > diceCommandMaxFace {
			face = diceCommandMaxFace
		}

		result := make([]string, 0, count)
		for range count {
			result = append(result, strconv.Itoa(rand.Intn(int(face))+1))
		}

		msg := strings.Join(result, " + ")
		embed := &discordgo.MessageEmbed{
			Title: "ダイスロール",
			Color: common.EmbedColor,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Result",
					Value:  fmt.Sprintf("```\n%dd%d → %s\n```", count, face, msg),
					Inline: true,
				},
			},
		}

		if err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond to the interaction", zap.Error(err))
		}
	}
}

func (c *DiceCommand) errorResponse(name string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "内部エラー",
				Color:       common.EmbedColor,
				Description: "エラーが発生したようです．責任持って修正してください．",
				Fields: []*discordgo.MessageEmbedField{{
					Name:  "Error",
					Value: fmt.Sprintf("Option `%s` is not found", name),
				}},
			}},
		},
	}
}
