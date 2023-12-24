package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func Test_NewZap(t *testing.T) {
	// Given:
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.DeploymentEnvironment("development"),
	)

	// When:
	l, err := NewZap(true, res)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)
	l.Info("testing")

	// When:
	l, err = NewZap(false, res)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)
	l.Info("testing")
}
