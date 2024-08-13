package command

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/aqyuki/tubu/packages/bot/common"
	"github.com/aqyuki/tubu/packages/logging"
	"github.com/aqyuki/tubu/packages/platform/discord"
	"github.com/bwmarrin/discordgo"
)

var _ discord.Command = (*DiceCommand)(nil)

const (
	diceCommandCountOptionName = "count"
	diceCommandFaceOptionName  = "face"
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
			return
		}
		faceOpt, ok := optionMap[diceCommandFaceOptionName]
		if !ok {
			return
		}

		count := countOpt.IntValue()
		face := faceOpt.IntValue()

		if count < 1 {
			count = 1
		}

		if face < 2 {
			face = 2
		} else if face > math.MaxInt {
			face = math.MaxInt
		}

		result := make([]string, count)
		for range count {
			result = append(result, strconv.Itoa(rand.Intn(int(face))+1))
		}

		msg := strings.Join(result, " ")
		embed := &discordgo.MessageEmbed{
			Title:       "サイコロ",
			Color:       common.EmbedColor,
			Description: fmt.Sprintf("Result\n```\n%dd%d → %s\n```", count, face, msg),
		}

		if err := s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}); err != nil {
			logger.Error("failed to respond to the interaction", "error", err)
		}
	}
}
