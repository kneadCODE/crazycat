package otelhttpserver

import (
	"context"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ResponseWriterWrapper is a wrapper around http.ResponseWriter to help with OTEL instrumentation
type ResponseWriterWrapper struct {
	http.ResponseWriter

	Ctx context.Context

	bodySize   int64
	statusCode int
	writeErr   error

	wroteHeader bool
}

// WriteHeader satisfies the interface and records the status code
func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.statusCode = statusCode
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write satisfies the interface and records the status code
func (w *ResponseWriterWrapper) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(p)
	n1 := int64(n)
	w.bodySize += n1
	w.writeErr = err
	trace.SpanFromContext(w.Ctx).AddEvent("http.response.write", trace.WithAttributes(
		attribute.Int64("http.response.wrote_bytes", n1),
	))
	return n, err
}

func (w *ResponseWriterWrapper) getHeaderAttrs() []attribute.KeyValue {
	var attrs []attribute.KeyValue

	if v := w.Header().Get("Content-Type"); v != "" {
		attrs = append(attrs, attribute.String("http.response.header.content-type", v))
	}
	if v, err := strconv.ParseInt(w.Header().Get("Content-Length"), 10, 64); err != nil {
		attrs = append(attrs, attribute.Int64("http.response.header.content-length", v))
	}

	// TODO: Check if any other header needs to be added

	return attrs
}
