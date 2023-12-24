package internal

import (
	"context"
	"testing"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func TestNewOTELPropagator(t *testing.T) {
	// Given && When:
	props := NewOTELPropagator(false)

	// Then:
	require.Equal(t, propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	), props)

	// When:
	props = NewOTELPropagator(true)

	// Then:
	require.Equal(t, propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}, sentryotel.NewSentryPropagator(),
	), props)
}

func TestNewOTELTraceProvider(t *testing.T) {
	// Given:
	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development"))

	// When:
	tp, err := NewOTELTraceProvider(res, false)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, tp)
	require.NotNil(t, tp.Tracer("test"))

	// Given && When:
	tp, err = NewOTELTraceProvider(res, true)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, tp)
	require.NotNil(t, tp.Tracer("test"))

	// TODO: Figure out how to write proper tests for OTEL configs. Only choice I see now is using interfaces :(
}

func TestNewOTELMeterProvider(t *testing.T) {
	// Given:
	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development"))

	// When:
	tp, err := NewOTELMeterProvider(res)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, tp)
	require.NotNil(t, tp.Meter("test"))

	// TODO: Figure out how to write proper tests for OTEL configs. Only choice I see now is using interfaces :(
}

func TestGetOTELTracer(t *testing.T) {
	// Given && When:
	tracer := GetTracer()

	// Then:
	require.NotNil(t, tracer)
	ctx, span := tracer.Start(context.Background(), "span 1")
	require.NotNil(t, span)
	require.NotNil(t, ctx)
}
