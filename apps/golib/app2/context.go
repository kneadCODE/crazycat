package app2

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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

	ctx = setOTELAttrsInContext(ctx, append(otelAttrsFromContext(ctx), attrs...))

	return ctx
}

// contextKey implementation is referenced from go stdlib:
// https://github.com/golang/go/blob/2184a394777ccc9ce9625932b2ad773e6e626be0/src/net/http/http.go#L42
type contextKey struct {
	name string
}

func (k contextKey) String() string { return "app context value " + k.name }

var (
	configCtxKey    = contextKey{"app-config"}
	zapCtxKey       = contextKey{"app-zap"}
	otelAttrsCtxKey = contextKey{"app-otel-attrs"}
	// newrelicCtxKey   = contextKey{"app_newrelic"}
)

func setConfigInContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}

func setZapInContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, zapCtxKey, l)
}

func zapFromContext(ctx context.Context) *zap.Logger {
	if v, ok := ctx.Value(zapCtxKey).(*zap.Logger); ok {
		return v
	}
	return nil
}

func setOTELAttrsInContext(ctx context.Context, attrs []attribute.KeyValue) context.Context {
	return context.WithValue(ctx, otelAttrsCtxKey, attrs)
}

func otelAttrsFromContext(ctx context.Context) []attribute.KeyValue {
	if v, ok := ctx.Value(otelAttrsCtxKey).([]attribute.KeyValue); ok {
		return v
	}
	return nil
}
