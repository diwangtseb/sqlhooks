package othooks

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Hook struct {
	tracer trace.Tracer
}

func New(tracer trace.Tracer) *Hook {
	return &Hook{tracer: tracer}
}

func (h *Hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	fmt.Println("caocoacoaocao")
	ctx, span := h.tracer.Start(ctx, "sql")
	span.SetAttributes(attribute.Key("query").String(fmt.Sprintf("%v", query)))
	span.SetAttributes(attribute.Key("args").String(fmt.Sprintf("%v", args)))
	return ctx, nil
}

func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		defer span.End()
	}
	return ctx, nil
}

func (h *Hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		defer span.End()
		span.SetAttributes(attribute.Key("error").String(fmt.Sprintf("%v", err)))
	}

	return err
}
