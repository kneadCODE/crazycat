package app2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestStartSpan(t *testing.T) {
	// Given:
	ctx, cancel := context.WithCancel(context.Background())

	// When && Then:
	secondCtx, endSecond := StartSpan(ctx, "span1", false, attribute.String("k1", "v1"))
	endSecond(nil)

	// When && Then:
	secondCtx, endSecond = StartSpan(ctx, "span1", false, attribute.String("k1", "v1"))
	endSecond(errors.New("some err"))

	// When && Then:
	secondCtx, endSecond = StartSpan(ctx, "span1", false, attribute.String("k1", "v1"))
	thirdCtx, endThird := StartSpan(secondCtx, "span2", false, attribute.String("k2", "v2"))
	fourthCtx, endFourth := StartSpan(secondCtx, "span3", true, attribute.String("k3", "v3"))
	endThird(errors.New("some err"))
	endSecond(nil)
	endFourth(nil)

	cancel()

	require.Equal(t, context.Canceled, secondCtx.Err())
	require.Equal(t, context.Canceled, thirdCtx.Err())
	require.NoError(t, fourthCtx.Err())
}
