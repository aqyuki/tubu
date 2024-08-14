package discord

import (
	"iter"

	"github.com/bwmarrin/discordgo"
)

func collectCommand(commands iter.Seq[Command]) iter.Seq[*discordgo.ApplicationCommand] {
	return func(yield func(*discordgo.ApplicationCommand) bool) {
		for command := range commands {
			if !yield(command.Command()) {
				return
			}
		}
	}
}
