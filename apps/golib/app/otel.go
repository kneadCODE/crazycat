package app

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newTracer(cfg Config) (*sdktrace.TracerProvider, error) {
	res, err := getOTELResource(cfg)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),

		// TODO: Check what other options need to be set
	)

	return tp, nil
}

func getOTELResource(cfg Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.Name),
		semconv.ServiceNamespace(cfg.Project),
		semconv.ServiceVersion(cfg.Version),
		semconv.ServiceInstanceID(cfg.ServerInstanceID),
		semconv.DeploymentEnvironment(string(cfg.Env)),

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
