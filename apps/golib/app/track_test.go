package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestTrackDebugEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	TrackDebugEvent(ctx, "message")
	TrackDebugEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err := internal.NewZap(config.Config{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackDebugEvent(ctx, "message")
	TrackDebugEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err = internal.NewZap(config.Config{Env: config.EnvDev})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackDebugEvent(ctx, "message")
	TrackDebugEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)
}

func TestTrackInfoEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	TrackInfoEvent(ctx, "message")
	TrackInfoEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err := internal.NewZap(config.Config{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackInfoEvent(ctx, "message")
	TrackInfoEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err = internal.NewZap(config.Config{Env: config.EnvDev})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackInfoEvent(ctx, "message")
	TrackInfoEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx, _ := trace.NewNoopTracerProvider().Tracer("testing").Start(ctx, "testing span")
	// When && Then:
	TrackInfoEvent(newCtx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx, _ = sdktrace.NewTracerProvider().Tracer("testing").Start(ctx, "testing span")
	// When && Then:
	TrackInfoEvent(newCtx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx = setNewRelicInContext(newCtx, &newrelic.Application{})
	// When && Then:
	TrackInfoEvent(newCtx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)
}

func TestTrackWarnEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	TrackWarnEvent(ctx, "message")
	TrackWarnEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err := internal.NewZap(config.Config{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackWarnEvent(ctx, "message")
	TrackWarnEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err = internal.NewZap(config.Config{Env: config.EnvDev})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackWarnEvent(ctx, "message")
	TrackWarnEvent(ctx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx, _ := trace.NewNoopTracerProvider().Tracer("testing").Start(ctx, "testing span")
	// When && Then:
	TrackWarnEvent(newCtx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx, _ = sdktrace.NewTracerProvider().Tracer("testing").Start(ctx, "testing span")
	// When && Then:
	TrackWarnEvent(newCtx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx = setNewRelicInContext(newCtx, &newrelic.Application{})
	// When && Then:
	TrackWarnEvent(newCtx, "message", "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)
}

func TestTrackErrorEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	TrackErrorEvent(ctx, errors.New("some err"))
	TrackErrorEvent(ctx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err := internal.NewZap(config.Config{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackErrorEvent(ctx, errors.New("some err"))
	TrackErrorEvent(ctx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	l, err = internal.NewZap(config.Config{Env: config.EnvDev})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	TrackErrorEvent(ctx, errors.New("some err"))
	TrackErrorEvent(ctx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx, _ := trace.NewNoopTracerProvider().Tracer("testing").Start(ctx, "testing span")
	// When && Then:
	TrackErrorEvent(newCtx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx, _ = sdktrace.NewTracerProvider().Tracer("testing").Start(ctx, "testing span")
	// When && Then:
	TrackErrorEvent(newCtx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// Given:
	newCtx = setNewRelicInContext(newCtx, &newrelic.Application{})
	// When && Then:
	TrackErrorEvent(newCtx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)

	// // Given: // TODO: Figure out how to test the sentry integration when we can't use a real Sentry key. Again, trying to avoid interfaces
	// newCtx = sentry.SetHubOnContext(newCtx, sentry.NewHub(&sentry.Client{}, sentry.NewScope()))
	// // When && Then:
	// TrackErrorEvent(newCtx, errors.New("some err"), "k1", "v1", "k2", 2, "k3", 3.0, "k4", true, "k5", 100*time.Second)
}
