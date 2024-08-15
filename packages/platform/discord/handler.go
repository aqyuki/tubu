package discord

import (
	"context"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type ReadyHandler func(context.Context, *discordgo.Session, *discordgo.Ready)

type MessageCreateHandler func(context.Context, *discordgo.Session, *discordgo.MessageCreate)

type Handler struct {
	readyHandler         []ReadyHandler
	messageCreateHandler []MessageCreateHandler
	contextFunc          func() context.Context
}

// HandlerOption is a functional option for Handler.
type HandlerOption func(*Handler)

// WithHandlerContextFunc adds a context function to the Handler.
func WithHandlerContextFunc(f func() context.Context) HandlerOption {
	return func(h *Handler) {
		if f == nil {
			return
		}
		h.contextFunc = f
	}
}

// WithReadyHandler adds a ReadyHandler to the Handler.
func WithReadyHandler(handler ReadyHandler) HandlerOption {
	return func(h *Handler) {
		h.readyHandler = append(h.readyHandler, handler)
	}
}

// WithMessageCreateHandler adds a MessageCreateHandler to the Handler.
func WithMessageCreateHandler(handler MessageCreateHandler) HandlerOption {
	return func(h *Handler) {
		h.messageCreateHandler = append(h.messageCreateHandler, handler)
	}
}

// NewHandler creates a new Handler.
func NewHandler(opts ...HandlerOption) *Handler {
	h := &Handler{
		readyHandler:         make([]ReadyHandler, 0),
		messageCreateHandler: make([]MessageCreateHandler, 0),
		contextFunc:          func() context.Context { return context.Background() },
	}

	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HandleReady handles the Ready event.
func (h *Handler) HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	var wg sync.WaitGroup
	ctx := h.contextFunc()
	wg.Add(len(h.readyHandler))
	for _, handler := range h.readyHandler {
		go func(h ReadyHandler) {
			defer recoveryPanic(ctx, r)
			h(ctx, s, r)
			wg.Done()
		}(handler)
	}
	wg.Wait()
}

// HandleMessageCreate handles the MessageCreate event.
func (h *Handler) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var wg sync.WaitGroup
	ctx := h.contextFunc()
	wg.Add(len(h.messageCreateHandler))
	for _, handler := range h.messageCreateHandler {
		go func(h MessageCreateHandler) {
			defer recoveryPanic(ctx, m)
			h(ctx, s, m)
			wg.Done()
		}(handler)
	}
	wg.Wait()
}
