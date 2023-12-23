package app2

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app2/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

// RecordDebugEvent records a debug event in both the logs and OTEL span
func RecordDebugEvent(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	recordCommon(ctx, zapcore.DebugLevel, msg, attrs...)
}

// RecordInfoEvent records an info event in both the logs and OTEL span
func RecordInfoEvent(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	recordCommon(ctx, zapcore.InfoLevel, msg, attrs...)
}

// RecordWarnEvent records a warning event in both the logs and OTEL span
func RecordWarnEvent(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	recordCommon(ctx, zapcore.WarnLevel, msg, attrs...)
}

// RecordError records an error in both the logs and OTEL span
func RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	attrs = append(attrs, internal.GetOTELErrorAttrs(2)...)

	span := trace.SpanFromContext(ctx)
	span.RecordError(err, trace.WithAttributes(attrs...))

	zapLogger := zapFromContext(ctx)
	if zapLogger == nil {
		return
	}

	internal.ZapLogEnriched(zapLogger, zapcore.ErrorLevel, err.Error(), span, attrs...)
}

func recordCommon(ctx context.Context, level zapcore.Level, msg string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(msg, trace.WithAttributes(attrs...))

	zapLogger := zapFromContext(ctx)
	if zapLogger == nil {
		return
	}

	internal.ZapLogEnriched(zapLogger, level, msg, span, attrs...)
}
