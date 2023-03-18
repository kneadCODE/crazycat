package app

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const configCtxKey = ctxKey("app_config")

const zapCtxKey = ctxKey("app_zap")

// ConfigFromContext retrieves the Config from context if exists else return a new Config
func ConfigFromContext(ctx context.Context) Config {
	if v, ok := ctx.Value(configCtxKey).(Config); ok {
		return v
	}
	return Config{}
}

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
