package app

import (
	"context"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

// TrackDebugEvent logs the given message at debug level
func TrackDebugEvent(ctx context.Context, msg string, fields ...interface{}) {
	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	l.Debugw(msg, injectOTELSpanFields(trace.SpanFromContext(ctx), fields...)...)
}

// TrackInfoEvent logs the given message at info level and adds an event to the span (if exists)
func TrackInfoEvent(ctx context.Context, msg string, fields ...interface{}) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(msg)

	if nrApp := newRelicFromContext(ctx); nrApp != nil {
		nrApp.RecordLog(newrelic.LogData{
			Timestamp: time.Now().UnixMilli(),
			Severity:  zapcore.InfoLevel.CapitalString(),
			Message:   msg,
		})
	}

	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	l.Infow(msg, injectOTELSpanFields(span, fields...)...)
}

// TrackWarnEvent logs the given message at warn level and adds an event to the span (if exists)
func TrackWarnEvent(ctx context.Context, msg string, fields ...interface{}) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(msg)

	if nrApp := newRelicFromContext(ctx); nrApp != nil {
		nrApp.RecordLog(newrelic.LogData{
			Timestamp: time.Now().UnixMilli(),
			Severity:  zapcore.WarnLevel.CapitalString(),
			Message:   msg,
		})
	}

	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	l.Warnw(msg, injectOTELSpanFields(span, fields...)...)
}

// TrackErrorEvent logs the given message at error level and also reports the error
func TrackErrorEvent(ctx context.Context, err error, fields ...interface{}) {
	stackTraceSlice := make([]byte, 2048)
	n := runtime.Stack(stackTraceSlice, false)
	stackTrace := string(stackTraceSlice[0:n])

	span := trace.SpanFromContext(ctx)

	attrs := []attribute.KeyValue{semconv.ExceptionStacktraceKey.String(stackTrace)}
	if fn, file, line, ok := runtimeCaller(1); ok {
		if fn != "" {
			attrs = append(attrs, semconv.CodeFunctionKey.String(fn))
		}
		if file != "" {
			attrs = append(attrs, semconv.CodeFilepathKey.String(file))
			attrs = append(attrs, semconv.CodeLineNumberKey.Int(line))
		}
	}
	span.RecordError(err, trace.WithAttributes(attrs...))

	if sentryHub := sentry.GetHubFromContext(ctx); sentryHub != nil {
		client, scope := sentryHub.Client(), sentryHub.Scope()
		// scope.SetTags(fields) // TODO: Figure out how to handle tags in a performant way
		_ = client.CaptureException(
			err,
			&sentry.EventHint{Context: ctx},
			scope,
		)
	}

	if nrApp := newRelicFromContext(ctx); nrApp != nil {
		nrApp.RecordLog(newrelic.LogData{
			Timestamp: time.Now().UnixMilli(),
			Severity:  zapcore.ErrorLevel.CapitalString(),
			Message:   err.Error(),
		})
	}

	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	l.Errorw(err.Error(), injectOTELSpanFields(span, fields...)...)
}

func injectOTELSpanFields(span trace.Span, fields ...interface{}) []interface{} {
	if v := span.SpanContext(); v.IsValid() {
		fields = append(fields,
			logFieldTraceID, v.TraceID(),
			logFieldSpanID, v.SpanID(),
			logFieldTraceFlags, v.TraceFlags(),
		)
	}
	return fields
}

func runtimeCaller(skip int) (fn, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame.Function, frame.File, frame.Line, frame.PC != 0
}
