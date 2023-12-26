package otelhttpserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/noop"
)

func TestNewMeasure(t *testing.T) {
	// Given && When:
	m, err := NewMeasure()
	require.NoError(t, err)

	require.NotNil(t, m.activeRequestCounter)
	require.NotNil(t, m.serverLatencyMeasure)
	require.NotNil(t, m.requestBytesCounter)
	require.NotNil(t, m.responseBytesCounter)

	// Given:
	otel.SetMeterProvider(noop.NewMeterProvider())

	// When:
	m, err = NewMeasure()
	require.NoError(t, err)

	require.NotNil(t, m.activeRequestCounter)
	require.NotNil(t, m.serverLatencyMeasure)
	require.NotNil(t, m.requestBytesCounter)
	require.NotNil(t, m.responseBytesCounter)
}

func TestMeasure_MeasurePreProcessing(t *testing.T) {
	// Given:
	otel.SetMeterProvider(noop.NewMeterProvider())
	m, err := NewMeasure()
	require.NoError(t, err)
	ctx := context.Background()

	// When && Then:
	m.MeasurePreProcessing(ctx, nil)

	// When && Then:
	m.MeasurePreProcessing(ctx, []attribute.KeyValue{attribute.String("k1", "v1")})
}

func TestMeasure_MeasurePostProcessing(t *testing.T) {
	// Given:
	otel.SetMeterProvider(noop.NewMeterProvider())
	m, err := NewMeasure()
	require.NoError(t, err)
	ctx := context.Background()
	rbw := &RequestBodyWrapper{}
	rww := &ResponseWriterWrapper{}

	// When && Then:
	m.MeasurePostProcessing(ctx, rww, rbw, time.Duration(0), nil)

	// Given:
	rbw.bodySize = 123
	rww.bodySize = 123

	// When && Then:
	m.MeasurePostProcessing(ctx, rww, rbw, time.Duration(123), []attribute.KeyValue{attribute.String("k1", "v1")})
}
