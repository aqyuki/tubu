package discord

import (
	"context"
	"maps"
	"slices"

	"github.com/bwmarrin/discordgo"
)

type InteractionCreateHandler func(context.Context, *discordgo.Session, *discordgo.InteractionCreate)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Handler() InteractionCreateHandler
}

type CommandRouter struct {
	commands    map[string]Command
	contextFunc func() context.Context
}

type CommandRouterOption func(*CommandRouter)

func WithCommand(command Command) CommandRouterOption {
	return func(router *CommandRouter) {
		name := command.Command().Name
		if _, ok := router.commands[name]; ok {
			panic("tried to register a command, but the same name command already exists")
		}
		router.commands[name] = command
	}
}

func WithCommandContextFunc(f func() context.Context) CommandRouterOption {
	return func(router *CommandRouter) {
		if f == nil {
			return
		} else if f() == nil {
			return
		}
		router.contextFunc = f
	}
}

func NewCommandRouter(opts ...CommandRouterOption) *CommandRouter {
	router := &CommandRouter{
		commands:    make(map[string]Command),
		contextFunc: func() context.Context { return context.Background() },
	}
	for _, opt := range opts {
		opt(router)
	}
	return router
}

func (r *CommandRouter) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := r.contextFunc()
	defer recoveryPanic(ctx, i)
	// if the interaction is not a command, ignore it
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	command, ok := r.commands[i.ApplicationCommandData().Name]
	if !ok {
		return
	}
	command.Handler()(ctx, s, i)
}

func (r *CommandRouter) Commands() []*discordgo.ApplicationCommand {
	return slices.Collect(collectCommand(maps.Values(r.commands)))
}
