package logging

import (
	"context"
	"log/slog"
	"slices"
)

type attrsKey struct{}

// With returns a copy of ctx that carries attrs. Loggers made by New add them
// to every line logged with that context, through methods like InfoContext.
func With(ctx context.Context, attrs ...slog.Attr) context.Context {
	prev, _ := ctx.Value(attrsKey{}).([]slog.Attr)
	return context.WithValue(ctx, attrsKey{}, slices.Concat(prev, attrs))
}

// contextHandler adds the attributes stored by With to each record, and then
// passes it on to the real handler.
type contextHandler struct {
	slog.Handler
}

func (h contextHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {

	// 1. Get the value stored in the context.
	value := ctx.Value(attrsKey{})

	// 2. Check whether it's a slice of slog.Attr.
	attributes, ok := value.([]slog.Attr)

	// 3. Add those attributes to this log record.
	if ok {
		record.AddAttrs(attributes...)
	}

	// 4. Pass the updated record to the real handler.
	innerHandler := h.Handler

	return innerHandler.Handle(ctx, record)
}

func (h contextHandler) WithAttrs(
	attrs []slog.Attr,
) slog.Handler {

	// 1. Create a new inner handler with the attributes.
	newInnerHandler := h.Handler.WithAttrs(attrs)

	// 2. Wrap it again to preserve context-aware logging.
	newWrapper := contextHandler{
		Handler: newInnerHandler,
	}

	return newWrapper
}

func (h contextHandler) WithGroup(
	name string,
) slog.Handler {

	// 1. Create a new inner handler with the group.
	newInnerHandler := h.Handler.WithGroup(name)

	// 2. Wrap it again.
	return contextHandler{
		Handler: newInnerHandler,
	}
}
