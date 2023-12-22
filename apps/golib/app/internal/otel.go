package internal

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	// logFieldTraceID is the field used for OTEL trace ID
	logFieldTraceID = "trace_id"
	// logFieldSpanID is the field used for OTEL span ID
	logFieldSpanID = "span_id"
	// logFieldTraceFlags is the field used for OTEL trace flags
	logFieldTraceFlags = "trace_flags"
)

type OTELMetadata struct {
	AppName          string
	AppEnv           string
	AppVersion       string
	Project          string
	ServerInstanceID string
}

// NewOTELProvider returns a new instance of OTEL provider
func NewOTELProvider(metadata OTELMetadata, isSentryEnabled bool) (*sdktrace.TracerProvider, error) {
	res, err := newOTELResource(metadata)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		// TODO: Check what other options need to be set
	)

	if isSentryEnabled {
		tp.RegisterSpanProcessor(sentryotel.NewSentrySpanProcessor())
	}

	// otel.SetTextMapPropagator(sentryotel.NewSentryPropagator()) // TODO: Look into writing custom propagator that can combine all propagators

	return tp, nil
}

// TODO: Add comprehensive tests for the resource creation
func newOTELResource(m OTELMetadata) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(m.AppName),
		semconv.ServiceNamespace(m.Project),
		semconv.ServiceVersion(m.AppVersion),
		semconv.ServiceInstanceID(m.ServerInstanceID),
		semconv.DeploymentEnvironment(m.AppEnv),

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
