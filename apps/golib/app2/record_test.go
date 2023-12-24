package app2

import (
	"context"
	"errors"
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app2/internal"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestRecordDebugEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	RecordDebugEvent(ctx, "message")
	RecordDebugEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	l, err := internal.NewZap(false, &resource.Resource{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	RecordDebugEvent(ctx, "message")
	RecordDebugEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	ctx, _ = internal.GetTracer().Start(ctx, "testing span")
	// When && Then:
	RecordDebugEvent(ctx, "message")
	RecordDebugEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)
}

func TestTrackInfoEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	RecordInfoEvent(ctx, "message")
	RecordInfoEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	l, err := internal.NewZap(false, &resource.Resource{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	RecordInfoEvent(ctx, "message")
	RecordInfoEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	ctx, _ = internal.GetTracer().Start(ctx, "testing span")
	// When && Then:
	RecordInfoEvent(ctx, "message")
	RecordInfoEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)
}

func TestTrackWarnEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	RecordWarnEvent(ctx, "message")
	RecordWarnEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	l, err := internal.NewZap(false, &resource.Resource{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	RecordWarnEvent(ctx, "message")
	RecordWarnEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	ctx, _ = internal.GetTracer().Start(ctx, "testing span")
	// When && Then:
	RecordWarnEvent(ctx, "message")
	RecordWarnEvent(ctx, "message",
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)
}

func TestTrackErrorEvent(t *testing.T) {
	// Given:
	ctx := context.Background()
	// When && Then:
	RecordError(ctx, errors.New("some err"))
	RecordError(ctx, errors.New("some err"),
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	l, err := internal.NewZap(false, &resource.Resource{})
	require.NoError(t, err)
	ctx = setZapInContext(ctx, l)
	// When && Then:
	RecordError(ctx, errors.New("some err"))
	RecordError(ctx, errors.New("some err"),
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)

	// Given:
	ctx, _ = internal.GetTracer().Start(ctx, "testing span")
	// When && Then:
	RecordError(ctx, errors.New("some err"))
	RecordError(ctx, errors.New("some err"),
		attribute.String("k1", "v1"),
		attribute.Int("k2", 2),
		attribute.Float64("k3", 3.0),
		attribute.Bool("k4", true),
	)
}
