package app

import (
	"context"
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
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
	attrs := internal.OTELAttrsFromContext(ctx)
	require.Nil(t, attrs)

	newAttrs := []attribute.KeyValue{attribute.String("k1", "v1")}

	// When:
	ctx = ContextWithAttributes(ctx, newAttrs...)

	// Then:
	attrs = internal.OTELAttrsFromContext(ctx)
	require.EqualValues(t, newAttrs, attrs)

	// Given:
	ctx, span := noop.NewTracerProvider().Tracer("test").Start(ctx, "span1", trace.WithAttributes(attribute.String("k2", "v2")))

	// When:
	ctx = ContextWithAttributes(ctx, attribute.String("k2", "v2"))

	// Then:
	attrs = internal.OTELAttrsFromContext(ctx)
	require.EqualValues(t, append(newAttrs, attribute.String("k2", "v2")), attrs)
	_ = span // TODO: figure out how to verify that attrs were set in the span
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
