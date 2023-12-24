package app

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel/attribute"
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
	return internal.ContextWithAttributes(ctx, attrs...)
}

var configCtxKey = internal.ContextKey{Name: "app-config"}

func setConfigInContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}
