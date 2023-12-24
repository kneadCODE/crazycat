package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

func TestConfigFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	cfg := ConfigFromContext(ctx)

	// Then:
	require.EqualValues(t, Config{}, cfg)

	// When:
	newCfg := Config{Env: EnvDev}
	ctx = context.WithValue(ctx, configCtxKey, newCfg)

	// When:
	cfg = ConfigFromContext(ctx)

	// Then:
	require.EqualValues(t, newCfg, cfg)
}

func TestContextWithAttributes(t *testing.T) {
	// Given:
	ctx := context.Background()
	attrs := otelAttrsFromContext(ctx)
	require.Nil(t, attrs)

	newAttrs := []attribute.KeyValue{attribute.String("k1", "v1")}

	// When:
	ctx = ContextWithAttributes(ctx, newAttrs...)

	// Then:
	attrs = otelAttrsFromContext(ctx)
	require.EqualValues(t, newAttrs, attrs)

	// Given:
	ctx, span := noop.NewTracerProvider().Tracer("test").Start(ctx, "span1", trace.WithAttributes(attribute.String("k2", "v2")))

	// When:
	ctx = ContextWithAttributes(ctx, attribute.String("k2", "v2"))

	// Then:
	attrs = otelAttrsFromContext(ctx)
	require.EqualValues(t, append(newAttrs, attribute.String("k2", "v2")), attrs)
	_ = span // TODO: figure out how to verify that attrs were set in the span
}

func TestContextKey_String(t *testing.T) {
	require.Equal(t, "app context value app-config", configCtxKey.String())
	require.Equal(t, "app context value app-zap", zapCtxKey.String())
	require.Equal(t, "app context value app-otel-attrs", otelAttrsCtxKey.String())
	require.Equal(t, "app context value abc", contextKey{"abc"}.String())
}

func Test_setConfigInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	cfg := ConfigFromContext(ctx)

	// Then:
	require.EqualValues(t, Config{}, cfg)

	// When:
	newCfg := Config{Env: EnvDev}
	ctx = setConfigInContext(ctx, newCfg)

	// When:
	require.EqualValues(t, newCfg, ConfigFromContext(ctx))
}

func Test_zapFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	l := zapFromContext(ctx)

	// Then:
	require.Nil(t, l)

	// When:
	newL := zap.NewExample()
	ctx = context.WithValue(ctx, zapCtxKey, newL)

	// When:
	l = zapFromContext(ctx)

	// Then:
	require.EqualValues(t, newL, l)
}

func Test_setZapInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	l := zapFromContext(ctx)

	// Then:
	require.Nil(t, l)

	// When:
	newL := zap.NewExample()
	ctx = setZapInContext(ctx, newL)

	// When:
	require.EqualValues(t, newL, zapFromContext(ctx))
}

func Test_otelAttrsFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	attrs := otelAttrsFromContext(ctx)

	// Then:
	require.Nil(t, attrs)

	// When:
	newAttrs := []attribute.KeyValue{attribute.String("k1", "v1")}
	ctx = context.WithValue(ctx, otelAttrsCtxKey, newAttrs)

	// When:
	attrs = otelAttrsFromContext(ctx)

	// Then:
	require.EqualValues(t, newAttrs, attrs)
}

func Test_setOTELAttrsInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	attrs := otelAttrsFromContext(ctx)

	// Then:
	require.Nil(t, attrs)

	// When:
	newAttrs := []attribute.KeyValue{attribute.String("k1", "v1")}
	ctx = setOTELAttrsInContext(ctx, newAttrs)

	// When:
	require.EqualValues(t, newAttrs, otelAttrsFromContext(ctx))
}
