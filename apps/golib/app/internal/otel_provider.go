package internal

import (
	"fmt"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var otelInstrumentationScope = instrumentation.Scope{
	Name:      "github.com/kneadCODE/crazycat/apps/golib/app",
	Version:   "v0.0.0",
	SchemaURL: semconv.SchemaURL,
}

func NewOTELPropagator(isSentryEnabled bool) propagation.TextMapPropagator {
	p := []propagation.TextMapPropagator{propagation.TraceContext{}, propagation.Baggage{}}

	if isSentryEnabled {
		p = append(p, sentryotel.NewSentryPropagator())
	}

	return propagation.NewCompositeTextMapPropagator(p...)
}

func NewOTELTraceProvider(res *resource.Resource, isSentryEnabled bool) (*sdktrace.TracerProvider, error) {
	// TODO: Implement the correct trace exporter
	traceExporter, err := stdouttrace.New()
	if err != nil {
		return nil, fmt.Errorf("traceExporter err: %w", err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res), sdktrace.WithBatcher(traceExporter))

	if isSentryEnabled {
		tp.RegisterSpanProcessor(sentryotel.NewSentrySpanProcessor())
	}

	return tp, nil
}

func NewOTELMeterProvider(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	// TODO: Implement the correct metric exporter
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, fmt.Errorf("metricExporter err: %w", err)
	}

	// TODO: Figure out how to silence the initial metrics logging
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)

	return meterProvider, nil
}

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer(
		otelInstrumentationScope.Name,
		trace.WithInstrumentationVersion(otelInstrumentationScope.Version),
	) // TODO: Fill options
}
