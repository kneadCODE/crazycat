package internal

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ContextWithAttributes adds the given attributes to the context and the span in the context (if any)
func ContextWithAttributes(ctx context.Context, attrs ...attribute.KeyValue) context.Context {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)

	ctx = SetOTELAttrsInContext(ctx, append(OTELAttrsFromContext(ctx), attrs...))

	return ctx
}

func SetOTELAttrsInContext(ctx context.Context, attrs []attribute.KeyValue) context.Context {
	return context.WithValue(ctx, otelAttrsCtxKey, attrs)
}

func OTELAttrsFromContext(ctx context.Context) []attribute.KeyValue {
	if v, ok := ctx.Value(otelAttrsCtxKey).([]attribute.KeyValue); ok {
		return v
	}
	return nil
}

// ContextKey implementation is referenced from go stdlib:
// https://github.com/golang/go/blob/2184a394777ccc9ce9625932b2ad773e6e626be0/src/net/http/http.go#L42
type ContextKey struct {
	Name string
}

func (k ContextKey) String() string { return "app context value " + k.Name }

var (
	zapCtxKey       = ContextKey{"app-zap"}
	otelAttrsCtxKey = ContextKey{"app-otel-attrs"}
	// newrelicCtxKey   = ContextKey{"app_newrelic"}
)

func SetZapInContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, zapCtxKey, l)
}

func ZapFromContext(ctx context.Context) *zap.Logger {
	if v, ok := ctx.Value(zapCtxKey).(*zap.Logger); ok {
		return v
	}
	return nil
}
