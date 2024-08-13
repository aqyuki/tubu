package discord

import (
	"context"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type ReadyHandler func(context.Context, *discordgo.Session, *discordgo.Ready)

type Handler struct {
	readyHandler []ReadyHandler
	contextFunc  func() context.Context
}

// HandlerOption is a functional option for Handler.
type HandlerOption func(*Handler)

// WithContextFunc adds a context function to the Handler.
func WithContextFunc(f func() context.Context) HandlerOption {
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

// NewHandler creates a new Handler.
func NewHandler(opts ...HandlerOption) *Handler {
	h := &Handler{
		readyHandler: make([]ReadyHandler, 0),
		contextFunc:  func() context.Context { return context.Background() },
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
	for _, handler := range h.readyHandler {
		wg.Add(1)
		go func(h ReadyHandler) {
			h(ctx, s, r)
			wg.Done()
		}(handler)
	}
	wg.Wait()
}
