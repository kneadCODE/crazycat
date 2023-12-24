package app2

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app2/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

// RecordDebugEvent records a debug event in both the logs and OTEL span
func RecordDebugEvent(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	recordCommon(ctx, zapcore.DebugLevel, msg, attrs)
}

// RecordInfoEvent records an info event in both the logs and OTEL span
func RecordInfoEvent(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	recordCommon(ctx, zapcore.InfoLevel, msg, attrs)
}

// RecordWarnEvent records a warning event in both the logs and OTEL span
func RecordWarnEvent(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	recordCommon(ctx, zapcore.WarnLevel, msg, attrs)
}

// RecordError records an error in both the logs and OTEL span
func RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	attrs = append(attrs, internal.GetOTELErrorAttrs(2)...)

	span := trace.SpanFromContext(ctx)
	span.RecordError(err, trace.WithAttributes(attrs...))

	zapL := zapFromContext(ctx)
	if zapL == nil {
		return
	}

	internal.ZapLogEnriched(zapL, zapcore.ErrorLevel, err.Error(), span, append(attrs, otelAttrsFromContext(ctx)...))
}

func recordCommon(ctx context.Context, level zapcore.Level, msg string, attrs []attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(msg, trace.WithAttributes(attrs...))

	zapL := zapFromContext(ctx)
	if zapL == nil {
		return
	}

	internal.ZapLogEnriched(zapL, level, msg, span, append(attrs, otelAttrsFromContext(ctx)...))
}

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
	tracer := internal.GetTracer()

	if async {
		newCtx = setZapInContext(context.Background(), zapFromContext(ctx).With())
		newCtx = trace.ContextWithSpan(newCtx, trace.SpanFromContext(ctx))
	} else {
		newCtx = ctx
	}

	newCtx, span := tracer.Start(newCtx, name, trace.WithAttributes(attrs...)) // TODO: Fill options

	newCtx = setOTELAttrsInContext(newCtx, attrs)

	return newCtx, func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}
}
