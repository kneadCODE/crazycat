package otelhttpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Measure measures server's inbound HTTP call performance
type Measure struct {
	requestBytesCounter  metric.Int64Histogram
	responseBytesCounter metric.Int64Histogram
	serverLatencyMeasure metric.Float64Histogram
	activeRequestCounter metric.Int64UpDownCounter
}

// MeasurePreProcessing records pre processing metrics
// TODO: Why does this have to be a pointer?
func (m *Measure) MeasurePreProcessing(ctx context.Context, attrs []attribute.KeyValue) {
	m.activeRequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// MeasurePostProcessing records post processing metrics
// TODO: Why does this have to be a pointer?
func (m *Measure) MeasurePostProcessing(
	ctx context.Context,
	rww *ResponseWriterWrapper,
	rbw *RequestBodyWrapper,
	elapsedTime time.Duration,
	attrs []attribute.KeyValue,
) {
	opts := metric.WithAttributes(attrs...)
	m.requestBytesCounter.Record(ctx, rbw.bodySize, opts)
	m.responseBytesCounter.Record(ctx, rww.bodySize, opts)
	m.serverLatencyMeasure.Record(ctx, elapsedTime.Seconds(), opts)
	m.activeRequestCounter.Add(ctx, -1, opts)
}

// NewMeasure returns a new instance of Measure
func NewMeasure() (*Measure, error) {
	meter := internal.GetMeter()

	// https://opentelemetry.io/docs/specs/semconv/http/http-metrics/#metric-httpserverrequestbodysize
	requestBytesCounter, err := meter.Int64Histogram(
		"http.server.request.body.size",
		metric.WithUnit("By"),
		metric.WithDescription("Size of HTTP server request bodies"),
	)
	if err != nil {
		return nil, fmt.Errorf("reqBytesCounter meter creation failed: %w", err)
	}

	// https://opentelemetry.io/docs/specs/semconv/http/http-metrics/#metric-httpserverresponsebodysize
	responseBytesCounter, err := meter.Int64Histogram(
		"http.server.response.body.size",
		metric.WithUnit("By"),
		metric.WithDescription("Size of HTTP server response bodies"),
	)
	if err != nil {
		return nil, fmt.Errorf("respBytesCounter meter creation failed: %w", err)
	}

	// https://opentelemetry.io/docs/specs/semconv/http/http-metrics/#metric-httpserverrequestduration
	serverLatencyMeasure, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithUnit("s"),
		metric.WithDescription("Duration of HTTP server requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("reqBytesCounter meter creation failed: %w", err)
	}

	// https://opentelemetry.io/docs/specs/semconv/http/http-metrics/#metric-httpserveractive_requests
	activeRequestCounter, err := meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithUnit("{request}"),
		metric.WithDescription("Number of active HTTP server requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("activeRequestCounter meter creation failed: %w", err)
	}

	return &Measure{
		requestBytesCounter:  requestBytesCounter,
		responseBytesCounter: responseBytesCounter,
		serverLatencyMeasure: serverLatencyMeasure,
		activeRequestCounter: activeRequestCounter,
	}, nil
}
