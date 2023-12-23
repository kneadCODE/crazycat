package app2

import (
	"context"

	"go.uber.org/zap"
)

// ConfigFromContext retrieves the Config from context if exists else return a new Config
func ConfigFromContext(ctx context.Context) Config {
	if v, ok := ctx.Value(configCtxKey).(Config); ok {
		return v
	}
	return Config{}
}

// contextKey implementation is referenced from go stdlib:
// https://github.com/golang/go/blob/2184a394777ccc9ce9625932b2ad773e6e626be0/src/net/http/http.go#L42
type contextKey struct {
	name string
}

func (k contextKey) String() string { return "app context value " + k.name }

var (
	configCtxKey = contextKey{"app-config"}
	zapCtxKey    = contextKey{"app-zap"}
	// otelTracerCtxKey = contextKey{"app-otel-tracer"}
	// newrelicCtxKey   = contextKey{"app_newrelic"}
)

func setConfigInContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}

func setZapInContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, zapCtxKey, l)
}

func zapFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(zapCtxKey).(*zap.Logger); ok {
		return l
	}
	return nil
}