package app

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ctxKey string

const configCtxKey = ctxKey("app_config")

const zapCtxKey = ctxKey("app_zap")

const otelTracerCtxKey = ctxKey("app_otel_tracer")

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

func setOTELTracerInContext(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, otelTracerCtxKey, tracer)
}

func otelTracerFromContext(ctx context.Context) trace.Tracer {
	if tracer, ok := ctx.Value(otelTracerCtxKey).(trace.Tracer); ok {
		return tracer
	}
	return nil
}
