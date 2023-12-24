package internal

import (
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// GetOTELErrorAttrs gets error related OTEL attributes
func GetOTELErrorAttrs(callSkipLevels int) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.ExceptionStacktrace(string(stackTraceStub())),
	}

	if fn, file, line, ok := runtimeCaller(callSkipLevels + 1); ok {
		if fn != "" {
			attrs = append(attrs, semconv.CodeFunction(fn))
		}
		if file != "" {
			attrs = append(attrs, semconv.CodeFilepath(file))
			attrs = append(attrs, semconv.CodeLineNumber(line))
		}
	}

	return attrs
}

// RuntimeCaller gets the caller code for a particular call
func runtimeCaller(callSkipLevels int) (fn, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(callSkipLevels+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame.Function, frame.File, frame.Line, frame.PC != 0
}
