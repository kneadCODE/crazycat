package app

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ConfigFromContext retrieves the Config from context if exists else return a new Config
func ConfigFromContext(ctx context.Context) config.Config {
	if v, ok := ctx.Value(configCtxKey).(config.Config); ok {
		return v
	}
	return config.Config{}
}

// contextKey implementation is referenced from go stdlib:
// https://github.com/golang/go/blob/2184a394777ccc9ce9625932b2ad773e6e626be0/src/net/http/http.go#L42
type contextKey struct {
	name string
}

func (k contextKey) String() string { return "app context value " + k.name }

var (
	configCtxKey     = contextKey{"app-config"}
	zapCtxKey        = contextKey{"app-zap"}
	otelTracerCtxKey = contextKey{"app-otel-tracer"}
	newrelicCtxKey   = contextKey{"app_newrelic"}
)

func setConfigInContext(ctx context.Context, cfg config.Config) context.Context {
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
