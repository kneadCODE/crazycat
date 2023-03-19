package app

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ctxKey string

const configCtxKey = ctxKey("app_config")

const zapCtxKey = ctxKey("app_zap")

const otelTracerCtxKey = ctxKey("app_otel_tracer")

const newrelicCtxKey = ctxKey("app_newrelic")

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

func setZapInContext(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, zapCtxKey, l)
}

func zapFromContext(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(zapCtxKey).(*zap.SugaredLogger); ok {
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

func setNewRelicInContext(ctx context.Context, nrApp *newrelic.Application) context.Context {
	return context.WithValue(ctx, newrelicCtxKey, nrApp)
}

func newRelicFromContext(ctx context.Context) *newrelic.Application {
	if nrApp, ok := ctx.Value(newrelicCtxKey).(*newrelic.Application); ok {
		return nrApp
	}
	return nil
}
