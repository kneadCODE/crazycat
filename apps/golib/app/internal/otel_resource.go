package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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
		semconv.ServiceInstanceID(getOTELEnvVar(semconv.ServiceInstanceIDKey)),
	}

	if v := getOTELEnvVar(semconv.ServiceNameKey); v == "" {
		return nil, errors.New("otel:svc name not provided")
	} else {
		attrs = append(attrs, semconv.ServiceName(v))
	}

	if v := getOTELEnvVar(semconv.ServiceNamespaceKey); v == "" {
		// OTEL considers this optional, but we will consider it mandatory to avoid mistakes
		return nil, errors.New("otel:svc namespace not provided")
	} else {
		attrs = append(attrs, semconv.ServiceNamespace(v))
	}

	if v := getOTELEnvVar(semconv.ServiceVersionKey); v == "" {
		// OTEL considers this optional, but we will consider it mandatory to avoid mistakes
		return nil, errors.New("otel:svc version not provided")
	} else {
		attrs = append(attrs, semconv.ServiceVersion(v))
	}

	return attrs, nil
}

func loadDeploymentResourceFromEnv() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.DeploymentEnvironment(getOTELEnvVar(semconv.DeploymentEnvironmentKey)),
	}
}

func loadContainerResourceFromEnv() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ContainerName(getOTELEnvVar(semconv.ContainerNameKey)),
		semconv.ContainerImageName(getOTELEnvVar(semconv.ContainerImageNameKey)),
		semconv.ContainerImageTag(getOTELEnvVar(semconv.ContainerImageTagKey)),
		semconv.ContainerRuntime(getOTELEnvVar(semconv.ContainerRuntimeKey)),
	}
}

func loadK8sResourceFromEnv() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.K8SClusterName(getOTELEnvVar(semconv.K8SClusterNameKey)),
		semconv.K8SNodeName(getOTELEnvVar(semconv.K8SNodeNameKey)),
		semconv.K8SNodeUID(getOTELEnvVar(semconv.K8SNodeUIDKey)),
		semconv.K8SNamespaceName(getOTELEnvVar(semconv.K8SNamespaceNameKey)),
		semconv.K8SPodName(getOTELEnvVar(semconv.K8SPodNameKey)),
		semconv.K8SPodUID(getOTELEnvVar(semconv.K8SPodUIDKey)),
		semconv.K8SContainerName(getOTELEnvVar(semconv.K8SContainerNameKey)),
		semconv.K8SDeploymentName(getOTELEnvVar(semconv.K8SDeploymentNameKey)),
		semconv.K8SJobName(getOTELEnvVar(semconv.K8SJobNameKey)),
		semconv.K8SCronJobName(getOTELEnvVar(semconv.K8SCronJobNameKey)),
	}

	intV, _ := strconv.Atoi(getOTELEnvVar(semconv.K8SContainerRestartCountKey)) // Intentionally suppressing the err since nothing to do
	attrs = append(attrs, semconv.K8SContainerRestartCount(intV))

	return attrs
}

func loadCloudResourceFromEnv() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.CloudProviderKey.String(getOTELEnvVar(semconv.CloudProviderKey)),
		semconv.CloudRegion(getOTELEnvVar(semconv.CloudRegionKey)),
		semconv.CloudAvailabilityZone(getOTELEnvVar(semconv.CloudAvailabilityZoneKey)),
		semconv.CloudPlatformKey.String(getOTELEnvVar(semconv.CloudPlatformKey)),
	}
}

func getOTELEnvVar(key attribute.Key) string {
	// Convert key into OTEL_<envvar> format and replace dots with underscore
	return os.Getenv(
		fmt.Sprintf("OTEL_%s",
			strings.ToUpper(
				strings.Replace(string(key), ".", "_", -1),
			),
		),
	)
}
