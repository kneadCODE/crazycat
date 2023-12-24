package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func TestNewOTELResourceFromEnv(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	res, err := NewOTELResourceFromEnv(ctx)

	// Then:
	require.NoError(t, err)
	require.Equal(t, semconv.SchemaURL, res.SchemaURL())

	v, ok := res.Set().Value(semconv.ServiceNameKey)
	require.True(t, ok)
	require.Equal(t, "golib", v.AsString()) // TODO: Get value from .env file

	v, ok = res.Set().Value(semconv.ServiceVersionKey)
	require.True(t, ok)
	require.Equal(t, "v0.0.0", v.AsString()) // TODO: Get value from .env file

	v, ok = res.Set().Value(semconv.ServiceNamespaceKey)
	require.True(t, ok)
	require.Equal(t, "crazycat", v.AsString()) // TODO: Get value from .env file

	v, ok = res.Set().Value(semconv.DeploymentEnvironmentKey)
	require.True(t, ok)
	require.Equal(t, "development", v.AsString()) // TODO: Get value from .env file

	// TODO: Verify the rest of the resource attrs are set correctly are not
}
