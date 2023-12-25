package app

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ConfigFromContext retrieves the Config from context if exists else return a new Config
func ConfigFromContext(ctx context.Context) Config {
	if v, ok := ctx.Value(configCtxKey).(Config); ok {
		return v
	}
	return Config{}
}

// ContextWithAttributes adds the given attributes to the context and the span in the context (if any)
func ContextWithAttributes(ctx context.Context, attrs ...attribute.KeyValue) context.Context {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)

	ctx = internal.SetOTELAttrsInContext(ctx, append(internal.OTELAttrsFromContext(ctx), attrs...))

	return ctx
}

// CloneNewContext returns a new context void of the signals of the given context but inclusive of Config, trace.Span,
// Zap and Attrs
func CloneNewContext(ctx context.Context) context.Context {
	newCtx := context.Background()

	newCtx = setConfigInContext(newCtx, ConfigFromContext(ctx))
	newCtx = trace.ContextWithSpan(newCtx, trace.SpanFromContext(ctx))
	newCtx = internal.SetZapInContext(newCtx, internal.ZapFromContext(ctx))
	newCtx = internal.SetOTELAttrsInContext(newCtx, internal.OTELAttrsFromContext(ctx))

	return newCtx
}

var configCtxKey = internal.ContextKey{Name: "app-config"}

func setConfigInContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}
