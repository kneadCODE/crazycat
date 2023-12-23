package internal

import (
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewOTELPropagator(isSentryEnabled bool) propagation.TextMapPropagator {
	p := []propagation.TextMapPropagator{propagation.TraceContext{}, propagation.Baggage{}}

	if isSentryEnabled {
		p = append(p, sentryotel.NewSentryPropagator())
	}

	return propagation.NewCompositeTextMapPropagator(p...)
}

func NewTraceProvider(res *resource.Resource, isSentryEnabled bool) (*trace.TracerProvider, error) {
	// TODO: Implement the correct trace exporter
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(trace.WithResource(res), trace.WithBatcher(traceExporter))

	if isSentryEnabled {
		tp.RegisterSpanProcessor(sentryotel.NewSentrySpanProcessor())
	}

	return tp, nil
}

func NewMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	// TODO: Implement the correct metric exporter
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	// TODO: Figure out how to silence the initial metrics logging
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)

	return meterProvider, nil
}
