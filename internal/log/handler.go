package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/aLieexe/tsukatsuki/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

var LevelStyleMap = map[slog.Level]lipgloss.Style{
	slog.LevelDebug: ui.DebugStyle,
	slog.LevelInfo:  ui.InfoStyle,
	slog.LevelWarn:  ui.WarnStyle,
	slog.LevelError: ui.ErrorStyle,
}

type Handler struct {
	opts  slog.HandlerOptions
	mu    *sync.Mutex
	w     io.Writer
	attrs []slog.Attr
	group string
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	return &Handler{
		opts: *opts,
		mu:   new(sync.Mutex),
		w:    w,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	style, ok := LevelStyleMap[r.Level]
	if !ok {
		style = lipgloss.NewStyle()
	}

	msg := r.Message

	var attrs strings.Builder

	// add handler's own attributes first
	for _, attr := range h.attrs {
		key := attr.Key
		if h.group != "" {
			key = h.group + "." + key
		}
		attrs.WriteString(fmt.Sprintf(" %s=%v", key, attr.Value.Any()))
	}

	// add record attributes
	r.Attrs(func(a slog.Attr) bool {
		key := a.Key
		if h.group != "" {
			key = h.group + "." + key
		}
		attrs.WriteString(fmt.Sprintf(" %s=%v", key, a.Value.Any()))
		return true
	})

	output := strings.TrimRight(style.Render(msg+attrs.String()), " \t\r\n")
	_, err := h.w.Write([]byte(output + "\n"))
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	newHandler := *h
	if h.group != "" {
		newHandler.group = h.group + "." + name
	} else {
		newHandler.group = name
	}
	return &newHandler
}

func Init(w io.Writer, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := NewHandler(w, opts)

	return slog.New(handler)
}
