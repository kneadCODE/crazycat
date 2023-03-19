package app

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func TestSomething(t *testing.T) {
	require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
	require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
	require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
	require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
	require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "development"))
	defer os.Unsetenv(string(semconv.ServiceNameKey))
	defer os.Unsetenv(string(semconv.ServiceNamespaceKey))
	defer os.Unsetenv(string(semconv.ServiceVersionKey))
	defer os.Unsetenv(string(semconv.ServiceInstanceIDKey))
	defer os.Unsetenv(string(semconv.DeploymentEnvironmentKey))

	ctx, finish, err := Init()
	if err != nil {
		log.Fatal(err)
	}
	defer finish()

	tracer := otelTracerFromContext(ctx)
	ctx, _ = tracer.Start(ctx, "some_span")

	LogDebug(ctx, "hello world")
	TrackInfoEvent(ctx, "hello world")
	TrackWarnEvent(ctx, "hello world")
	TrackErrorEvent(ctx, errors.New("some err"))

	t.Fail()
}
