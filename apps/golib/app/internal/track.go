package internal

import (
	"runtime"

	"go.opentelemetry.io/otel/trace"
)

// AppendOTELSpanFields injects OTEL span fields into the given set of fields
func AppendOTELSpanFields(span trace.Span, fields ...any) []any {
	if v := span.SpanContext(); v.IsValid() {
		fields = append(fields,
			logFieldTraceID, v.TraceID(),
			logFieldSpanID, v.SpanID(),
			logFieldTraceFlags, v.TraceFlags(),
		)
	}
	return fields
}

// RuntimeCaller gets the caller code for a particular call
func RuntimeCaller(skip int) (fn, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame.Function, frame.File, frame.Line, frame.PC != 0
}
