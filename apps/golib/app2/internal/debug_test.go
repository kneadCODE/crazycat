package internal

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func TestGetOTELErrorAttrs(t *testing.T) {
	// Given && When:
	attrs := GetOTELErrorAttrs(0)

	// Then:
	require.Equal(t, 4, len(attrs))

	// Given:
	defer func() {
		stackTraceStub = debug.Stack
	}()
	stackTraceStub = func() []byte {
		return []byte("stack_trace")
	}

	// When:
	attrs = GetOTELErrorAttrs(0)

	// Then:
	dir, err := os.Getwd()
	require.NoError(t, err)

	require.Equal(t, []attribute.KeyValue{
		semconv.ExceptionStacktrace("stack_trace"),
		semconv.CodeFunction("github.com/kneadCODE/crazycat/apps/golib/app2/internal.GetOTELErrorAttrs"),
		semconv.CodeFilepath(fmt.Sprintf("%s/debug.go", dir)),
		semconv.CodeLineNumber(16),
	}, attrs)
}
