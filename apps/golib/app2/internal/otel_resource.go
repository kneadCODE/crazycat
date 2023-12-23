package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// NewOTELResourceFromEnv returns a new instance of OTEL resource from env
func NewOTELResourceFromEnv(ctx context.Context) (*resource.Resource, error) {
	attrs, err := loadServiceResourceFromEnv()
	if err != nil {
		return nil, err
	}

	attrs = append(attrs, loadDeploymentResourceFromEnv()...)
	attrs = append(attrs, loadContainerResourceFromEnv()...)
	attrs = append(attrs, loadK8sResourceFromEnv()...)
	attrs = append(attrs, loadCloudResourceFromEnv()...)

	res, err := resource.New(
		ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithTelemetrySDK(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithContainer(),
		resource.WithProcess(),
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return nil, fmt.Errorf("err creating OTEL resource: %w", err)
	}

	return res, nil
}

func loadServiceResourceFromEnv() ([]attribute.KeyValue, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceInstanceID(os.Getenv(string(semconv.ServiceInstanceIDKey))),
	}

	if v := os.Getenv(string(semconv.ServiceNameKey)); v == "" {
		return nil, errors.New("otel:svc name not provided")
	} else {
		attrs = append(attrs, semconv.ServiceName(v))
	}

	if v := os.Getenv(string(semconv.ServiceNamespaceKey)); v == "" {
		// OTEL considers this optional, but we will consider it mandatory to avoid mistakes
		return nil, errors.New("otel:svc namespace not provided")
	} else {
		attrs = append(attrs, semconv.ServiceNamespace(v))
	}

	if v := os.Getenv(string(semconv.ServiceVersionKey)); v == "" {
		// OTEL considers this optional, but we will consider it mandatory to avoid mistakes
		return nil, errors.New("otel:svc version not provided")
	} else {
		attrs = append(attrs, semconv.ServiceNamespace(v))
	}

	return attrs, nil
}

func loadDeploymentResourceFromEnv() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.DeploymentEnvironment(os.Getenv(string(semconv.DeploymentEnvironmentKey))),
	}
}

func loadContainerResourceFromEnv() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ContainerName(os.Getenv(string(semconv.ContainerNameKey))),
		semconv.ContainerImageName(os.Getenv(string(semconv.ContainerImageNameKey))),
		semconv.ContainerImageTag(os.Getenv(string(semconv.ContainerImageTagKey))),
		semconv.ContainerRuntime(os.Getenv(string(semconv.ContainerRuntimeKey))),
	}
}

func loadK8sResourceFromEnv() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.K8SClusterName(os.Getenv(string(semconv.K8SClusterNameKey))),
		semconv.K8SNodeName(os.Getenv(string(semconv.K8SNodeNameKey))),
		semconv.K8SNodeUID(os.Getenv(string(semconv.K8SNodeUIDKey))),
		semconv.K8SNamespaceName(os.Getenv(string(semconv.K8SNamespaceNameKey))),
		semconv.K8SPodName(os.Getenv(string(semconv.K8SPodNameKey))),
		semconv.K8SPodUID(os.Getenv(string(semconv.K8SPodUIDKey))),
		semconv.K8SContainerName(os.Getenv(string(semconv.K8SContainerNameKey))),
		semconv.K8SDeploymentName(os.Getenv(string(semconv.K8SDeploymentNameKey))),
		semconv.K8SJobName(os.Getenv(string(semconv.K8SJobNameKey))),
		semconv.K8SCronJobName(os.Getenv(string(semconv.K8SCronJobNameKey))),
	}

	intV, _ := strconv.Atoi(os.Getenv(string(semconv.K8SContainerRestartCountKey))) // Intentionally suppressing the err since nothing to do
	attrs = append(attrs, semconv.K8SContainerRestartCount(intV))

	return attrs
}

func loadCloudResourceFromEnv() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.CloudProviderKey.String(os.Getenv(string(semconv.CloudProviderKey))),
		semconv.CloudRegion(os.Getenv(string(semconv.CloudRegionKey))),
		semconv.CloudAvailabilityZone(os.Getenv(string(semconv.CloudAvailabilityZoneKey))),
		semconv.CloudPlatformKey.String(os.Getenv(string(semconv.CloudPlatformKey))),
	}
}
