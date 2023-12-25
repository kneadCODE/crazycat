package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func TestContextKey_String(t *testing.T) {
	require.Equal(t, "app context value app-zap", zapCtxKey.String())
	require.Equal(t, "app context value app-otel-attrs", otelAttrsCtxKey.String())
	require.Equal(t, "app context value abc", ContextKey{"abc"}.String())
}

func TestZapFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	l := ZapFromContext(ctx)

	// Then:
	require.Nil(t, l)

	// When:
	newL := zap.NewExample()
	ctx = context.WithValue(ctx, zapCtxKey, newL)

	// When:
	l = ZapFromContext(ctx)

	// Then:
	require.EqualValues(t, newL, l)
}

func TestSetZapInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	l := ZapFromContext(ctx)

	// Then:
	require.Nil(t, l)

	// When:
	newL := zap.NewExample()
	ctx = SetZapInContext(ctx, newL)

	// When:
	require.EqualValues(t, newL, ZapFromContext(ctx))
}

func TestOTELAttrsFromContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	attrs := OTELAttrsFromContext(ctx)

	// Then:
	require.Nil(t, attrs)

	// When:
	newAttrs := []attribute.KeyValue{attribute.String("k1", "v1")}
	ctx = context.WithValue(ctx, otelAttrsCtxKey, newAttrs)

	// When:
	attrs = OTELAttrsFromContext(ctx)

	// Then:
	require.EqualValues(t, newAttrs, attrs)
}

func TestSetOTELAttrsInContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When:
	attrs := OTELAttrsFromContext(ctx)

	// Then:
	require.Nil(t, attrs)

	// When:
	newAttrs := []attribute.KeyValue{attribute.String("k1", "v1")}
	ctx = SetOTELAttrsInContext(ctx, newAttrs)

	// When:
	require.EqualValues(t, newAttrs, OTELAttrsFromContext(ctx))
}
