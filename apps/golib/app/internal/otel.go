package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	// logFieldTraceID is the field used for OTEL trace ID
	logFieldTraceID = "trace_id"
	// logFieldSpanID is the field used for OTEL span ID
	logFieldSpanID = "span_id"
	// logFieldTraceFlags is the field used for OTEL trace flags
	logFieldTraceFlags = "trace_flags"
)

// InitOTEL initializes OTEL providers at global level
func InitOTEL(cfg config.Config, isSentryEnabled bool) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		return err
	}

	res, err := newOTELResource(cfg)
	if err != nil {
		return nil, err
	}

	prop := newOTELPropagator(isSentryEnabled)
	otel.SetTextMapPropagator(prop)

	tp, err := newTraceProvider(res, isSentryEnabled)
	if err != nil {
		return
	}
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
	otel.SetTracerProvider(tp) // Recommended way for OTEL go is to set at global level

	// Set up meter provider.
	mp, err := newMeterProvider(res)
	if err != nil {
		return
	}
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)
	otel.SetMeterProvider(mp) // Recommended way for OTEL go is to set at global level

	// Eventually when logs support is available, this is where to init the logs provider

	return
}

// TODO: Add comprehensive tests for the resource creation
func newOTELResource(cfg config.Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.Name),
		semconv.ServiceNamespace(cfg.Project),
		semconv.ServiceVersion(cfg.Version),
		semconv.ServiceInstanceID(cfg.ServerInstanceID),
		semconv.DeploymentEnvironment(cfg.Env.String()),

		semconv.OSTypeKey.String(runtime.GOOS),

		semconv.ContainerName(os.Getenv(string(semconv.ContainerNameKey))),
		semconv.ContainerID(os.Getenv(string(semconv.ContainerIDKey))),
		semconv.ContainerImageName(os.Getenv(string(semconv.ContainerImageNameKey))),
		semconv.ContainerImageTag(os.Getenv(string(semconv.ContainerImageTagKey))),
		semconv.ContainerRuntime(os.Getenv(string(semconv.ContainerRuntimeKey))),

		semconv.K8SClusterName(os.Getenv(string(semconv.K8SClusterNameKey))),
		semconv.K8SNodeName(os.Getenv(string(semconv.K8SNodeNameKey))),
		semconv.K8SNodeUID(os.Getenv(string(semconv.K8SNodeUIDKey))),
		semconv.K8SNamespaceName(os.Getenv(string(semconv.K8SNamespaceNameKey))),
		semconv.K8SPodName(os.Getenv(string(semconv.K8SPodNameKey))),
		semconv.K8SPodUID(os.Getenv(string(semconv.K8SPodUIDKey))),
		semconv.K8SContainerName(os.Getenv(string(semconv.K8SContainerNameKey))),

		semconv.CloudProviderKey.String(os.Getenv(string(semconv.CloudProviderKey))),
		semconv.CloudRegion(os.Getenv(string(semconv.CloudRegionKey))),
		semconv.CloudAvailabilityZone(os.Getenv(string(semconv.CloudAvailabilityZoneKey))),
		semconv.CloudPlatformKey.String(os.Getenv(string(semconv.CloudPlatformKey))),
	}

	intV, _ := strconv.Atoi(os.Getenv(string(semconv.K8SContainerRestartCountKey))) // Intentionally suppressing the err since nothing to do
	attrs = append(attrs, semconv.K8SContainerRestartCount(intV))

	var v string
	if v = os.Getenv(string(semconv.K8SDeploymentNameKey)); v != "" {
		attrs = append(attrs, semconv.K8SDeploymentName(v))
	}
	if v = os.Getenv(string(semconv.K8SJobNameKey)); v != "" {
		attrs = append(attrs, semconv.K8SJobName(v))
	}
	if v = os.Getenv(string(semconv.K8SCronJobNameKey)); v != "" {
		attrs = append(attrs, semconv.K8SCronJobName(v))
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	))
	if err != nil {
		return nil, fmt.Errorf("err merging OTEL resource: %w", err)
	}
	return res, nil
}

func newOTELPropagator(isSentryEnabled bool) propagation.TextMapPropagator {
	p := []propagation.TextMapPropagator{
		propagation.TraceContext{},
		propagation.Baggage{},
	}

	if isSentryEnabled {
		p = append(p, sentryotel.NewSentryPropagator())
	}

	return propagation.NewCompositeTextMapPropagator(p...)
}

func newTraceProvider(res *resource.Resource, isSentryEnabled bool) (*trace.TracerProvider, error) {
	// TODO: Implement the correct trace exporter
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(traceExporter),
	)

	if isSentryEnabled {
		tp.RegisterSpanProcessor(sentryotel.NewSentrySpanProcessor())
	}

	return tp, nil
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	// TODO: Implement the correct metric exporter
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)

	return meterProvider, nil
}
