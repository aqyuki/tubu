package discord

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	type except struct {
		readyHandlerCount         int
		messageCreateHandlerCount int
	}

	cases := []struct {
		name   string
		opts   []HandlerOption
		except *except
	}{
		{
			name: "should be return default handler",
			opts: []HandlerOption{},
			except: &except{
				readyHandlerCount:         0,
				messageCreateHandlerCount: 0,
			},
		},
		{
			name: "should be return handler with ready handler",
			opts: []HandlerOption{WithReadyHandler(func(context.Context, *discordgo.Session, *discordgo.Ready) {})},
			except: &except{
				readyHandlerCount:         1,
				messageCreateHandlerCount: 0,
			},
		},
		{
			name: "should be return handler with message create handler",
			opts: []HandlerOption{WithMessageCreateHandler(func(context.Context, *discordgo.Session, *discordgo.MessageCreate) {})},
			except: &except{
				readyHandlerCount:         0,
				messageCreateHandlerCount: 1,
			},
		},
		{
			name: "should be return handler with context function",
			opts: []HandlerOption{WithHandlerContextFunc(func() context.Context { return context.Background() })},
			except: &except{
				readyHandlerCount:         0,
				messageCreateHandlerCount: 0,
			},
		},
		{
			name: "should be return handler without invalid context function",
			opts: []HandlerOption{WithHandlerContextFunc(nil)},
			except: &except{
				readyHandlerCount:         0,
				messageCreateHandlerCount: 0,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := NewHandler(tc.opts...)
			assert.NotNil(t, actual)
			assert.Equal(t, tc.except.readyHandlerCount, len(actual.readyHandler))
			assert.Equal(t, tc.except.messageCreateHandlerCount, len(actual.messageCreateHandler))
			assert.NotNil(t, actual.contextFunc)
			assert.NotNil(t, actual.contextFunc())
		})
	}
}

func TestHandleReady(t *testing.T) {
	t.Parallel()

	c := context.Background()
	contextFunc := func() context.Context {
		return c
	}
	session := &discordgo.Session{}
	ready := &discordgo.Ready{}

	mockReadyHandler := func(ctx context.Context, s *discordgo.Session, r *discordgo.Ready) {
		assert.Equal(t, c, ctx)
		assert.Equal(t, session, s)
		assert.Equal(t, ready, r)
	}
	handler := NewHandler(
		WithHandlerContextFunc(contextFunc),
		WithReadyHandler(mockReadyHandler),
		WithReadyHandler(mockReadyHandler),
	)
	handler.HandleReady(session, ready)
}

func TestHandleMessageCreate(t *testing.T) {
	t.Parallel()

	c := context.Background()
	contextFunc := func() context.Context {
		return c
	}
	session := &discordgo.Session{}
	messageCreate := &discordgo.MessageCreate{}

	mockMessageCreateHandler := func(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
		assert.Equal(t, c, ctx)
		assert.Equal(t, session, s)
		assert.Equal(t, messageCreate, m)
	}
	handler := NewHandler(
		WithHandlerContextFunc(contextFunc),
		WithMessageCreateHandler(mockMessageCreateHandler),
		WithMessageCreateHandler(mockMessageCreateHandler),
	)
	handler.HandleMessageCreate(session, messageCreate)
}
