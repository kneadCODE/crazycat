package app

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a new span and returns the context with the span and the end func
// If a span already exists inside the given ctx, the new span is created as a child of the parent span.
// If async is set to true, then the newCtx is separated from the old ctx's signals.
// Intentionally didn't use an options pattern for async to force devs to pay attention to when to use sync/async.
func StartSpan(
	ctx context.Context,
	name string,
	async bool,
	attrs ...attribute.KeyValue,
) (newCtx context.Context, end func(error)) {

	if async {
		newCtx = internal.SetZapInContext(context.Background(), internal.ZapFromContext(ctx).With())
		newCtx = trace.ContextWithSpan(newCtx, trace.SpanFromContext(ctx))
	} else {
		newCtx = ctx
	}

	newCtx, span := internal.GetTracer().Start(newCtx, name, trace.WithAttributes(attrs...)) // TODO: Fill options

	newCtx = internal.SetOTELAttrsInContext(newCtx, attrs)

	return newCtx, func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}
}
