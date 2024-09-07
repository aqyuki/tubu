package discord

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

var _ Command = (*MockCommand)(nil)

type MockCommand struct {
	command    *discordgo.ApplicationCommand
	handleFunc InteractionCreateHandler
}

func (m *MockCommand) Command() *discordgo.ApplicationCommand {
	return m.command
}

func (m *MockCommand) Handler() InteractionCreateHandler {
	return m.handleFunc
}

func NewMockCommand(command *discordgo.ApplicationCommand, handleFunc InteractionCreateHandler) *MockCommand {
	return &MockCommand{
		command:    command,
		handleFunc: handleFunc,
	}
}

func TestWithCommand(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() {
		cmd1 := NewMockCommand(&discordgo.ApplicationCommand{Name: "test1"}, nil)
		cmd2 := NewMockCommand(&discordgo.ApplicationCommand{Name: "test1"}, nil)
		NewCommandRouter(WithCommand(cmd1), WithCommand(cmd2))
	})
}

func TestNewCommandRouter(t *testing.T) {
	t.Parallel()

	type except struct {
		exceptCommandCount int
	}

	cases := []struct {
		name   string
		opts   []CommandRouterOption
		except *except
	}{
		{
			name:   "should return a CommandRouter with default values",
			opts:   []CommandRouterOption{},
			except: &except{exceptCommandCount: 0},
		},
		{
			name:   "should return a CommandRouter with a command",
			opts:   []CommandRouterOption{WithCommand(NewMockCommand(&discordgo.ApplicationCommand{Name: "test"}, nil))},
			except: &except{exceptCommandCount: 1},
		},
		{
			name:   "should return a CommandRouter with a command and a context function",
			opts:   []CommandRouterOption{WithCommandContextFunc(func() context.Context { return context.Background() })},
			except: &except{exceptCommandCount: 0},
		},
		{
			name:   "should return a CommandRouter when context function is nil",
			opts:   []CommandRouterOption{WithCommandContextFunc(nil)},
			except: &except{exceptCommandCount: 0},
		},
		{
			name:   "should return a CommandRouter without a invalid context function",
			opts:   []CommandRouterOption{WithCommandContextFunc(func() context.Context { return nil })},
			except: &except{exceptCommandCount: 0},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCommandRouter(tc.opts...)
			assert.NotNil(t, actual)
			assert.Equal(t, tc.except.exceptCommandCount, len(actual.commands))
			assert.NotNil(t, actual.contextFunc)
			assert.NotNil(t, actual.contextFunc())
		})
	}
}

func TestHandleInteractionCreate(t *testing.T) {
	t.Parallel()

	t.Run("should ignore the interaction if it is not a command", func(t *testing.T) {
		t.Parallel()

		mockHandler := func(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) {
			t.Error("should not be called")
		}

		router := NewCommandRouter(WithCommand(NewMockCommand(&discordgo.ApplicationCommand{Name: "test"}, mockHandler)))
		session := &discordgo.Session{}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionMessageComponent,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "test",
				},
			},
		}
		router.HandleInteractionCreate(session, interaction)
	})

	t.Run("should ignore the interaction if the command does not exist", func(t *testing.T) {
		t.Parallel()

		mockHandler := func(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) {
			t.Error("should not be called")
		}

		router := NewCommandRouter(WithCommand(NewMockCommand(&discordgo.ApplicationCommand{Name: "test"}, mockHandler)))
		session := &discordgo.Session{}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "invalid",
				},
			},
		}
		router.HandleInteractionCreate(session, interaction)
	})

	t.Run("should call the handler if the command exists", func(t *testing.T) {
		t.Parallel()

		c := context.Background()
		contextFunc := func() context.Context { return c }
		session := &discordgo.Session{}
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: discordgo.InteractionApplicationCommand,
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "test",
				},
			},
		}

		mockHandler := func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
			assert.Equal(t, c, ctx)
			assert.Equal(t, session, s)
			assert.Equal(t, interaction, i)
		}

		router := NewCommandRouter(
			WithCommandContextFunc(contextFunc),
			WithCommand(NewMockCommand(&discordgo.ApplicationCommand{Name: "test"}, mockHandler)),
		)
		router.HandleInteractionCreate(session, interaction)
	})
}

func TestCommands(t *testing.T) {
	t.Parallel()

	cmd1 := &discordgo.ApplicationCommand{Name: "test1"}
	cmd2 := &discordgo.ApplicationCommand{Name: "test2"}
	except := []*discordgo.ApplicationCommand{cmd1, cmd2}

	router := NewCommandRouter(
		WithCommand(NewMockCommand(cmd1, nil)),
		WithCommand(NewMockCommand(cmd2, nil)),
	)
	actual := router.Commands()
	assert.ElementsMatch(t, except, actual)
}
