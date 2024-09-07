package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

var _ discord.Command = (*DiceRollService)(nil)

const (
	diceCommandCountOptionName = "count"
	diceCommandFaceOptionName  = "face"

	diceCommandMinCount = 1
	diceCommandMinFace  = 2
	diceCommandMaxCount = 10
	diceCommandMaxFace  = 100
)

type DiceRollService struct{}

func NewDiceService() *DiceRollService {
	return &DiceRollService{}
}

func (s *DiceRollService) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "dice",
		Description: "サイコロを振ります.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        diceCommandCountOptionName,
				Description: "振るサイコロの個数を指定してください.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        diceCommandFaceOptionName,
				Description: "振るサイコロの面数を指定してください.",
				Required:    true,
			},
		},
	}
}

func (s *DiceRollService) Handler() discord.InteractionCreateHandler {
	return func(ctx context.Context, session *discordgo.Session, ic *discordgo.InteractionCreate) {
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
			session.InteractionRespond(ic.Interaction, s.errorResponse(diceCommandCountOptionName))
			return
		}
		faceOpt, ok := optionMap[diceCommandFaceOptionName]
		if !ok {
			logger.Error("option is not found", zap.String("name", diceCommandFaceOptionName))
			session.InteractionRespond(ic.Interaction, s.errorResponse(diceCommandFaceOptionName))
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
			Color: EmbedColor,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Result",
					Value:  fmt.Sprintf("```\n%dd%d → %s\n```", count, face, msg),
					Inline: true,
				},
			},
		}

		if err := session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond to the interaction", zap.Error(err))
		}
	}
}

func (s *DiceRollService) errorResponse(name string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "内部エラー",
				Color:       EmbedColor,
				Description: "サイコロを振るのに失敗しちゃいました．",
				Fields: []*discordgo.MessageEmbedField{{
					Name:  "Error",
					Value: fmt.Sprintf("Option `%s` is not found", name),
				}},
			}},
		},
	}
}
