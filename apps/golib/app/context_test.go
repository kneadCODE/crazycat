package app

import (
	"context"
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func TestConfigFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	cfg := ConfigFromContext(ctx)

	// Then:
	require.EqualValues(t, config.Config{}, cfg)

	// When:
	newCfg := config.Config{Name: "something"}
	ctx = context.WithValue(ctx, configCtxKey, newCfg)

	// When:
	cfg = ConfigFromContext(ctx)

	// Then:
	require.EqualValues(t, newCfg, cfg)
}

func Test_setConfigInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	cfg := ConfigFromContext(ctx)

	// Then:
	require.EqualValues(t, config.Config{}, cfg)

	// When:
	newCfg := config.Config{Name: "something"}
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
	newL := zap.NewExample().Sugar()
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
	newL := zap.NewExample().Sugar()
	ctx = setZapInContext(ctx, newL)

	// When:
	require.EqualValues(t, newL, zapFromContext(ctx))
}

func Test_newRelicFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	nrApp := newRelicFromContext(ctx)

	// Then:
	require.Nil(t, nrApp)

	// When:
	newNRApp := &newrelic.Application{}
	ctx = context.WithValue(ctx, newrelicCtxKey, newNRApp)

	// When:
	nrApp = newRelicFromContext(ctx)

	// Then:
	require.EqualValues(t, newNRApp, nrApp)
}

func Test_setNewRelicInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	nrApp := newRelicFromContext(ctx)

	// Then:
	require.Nil(t, nrApp)

	// When:
	newNRApp := &newrelic.Application{}
	ctx = setNewRelicInContext(ctx, newNRApp)

	// When:
	require.EqualValues(t, newNRApp, newRelicFromContext(ctx))
}

func Test_otelTracerFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	tracer := otelTracerFromContext(ctx)

	// Then:
	require.Nil(t, tracer)

	// When:
	newTracer := sdktrace.NewTracerProvider().Tracer("testing")
	ctx = context.WithValue(ctx, otelTracerCtxKey, newTracer)

	// When:
	tracer = otelTracerFromContext(ctx)

	// Then:
	require.EqualValues(t, newTracer, tracer)
}

func Test_setOTELTracerInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	tracer := otelTracerFromContext(ctx)

	// Then:
	require.Nil(t, tracer)

	// When:
	newTracer := sdktrace.NewTracerProvider().Tracer("testing")
	ctx = setOTELTracerInContext(ctx, newTracer)

	// When:
	require.EqualValues(t, newTracer, otelTracerFromContext(ctx))
}
