package otelhttpserver

import (
	"context"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type OTELResponseWrapper struct {
	http.ResponseWriter

	Ctx context.Context

	BodySize   int64
	StatusCode int
	WriteErr   error

	wroteHeader bool
}

func (w *OTELResponseWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.StatusCode = statusCode
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *OTELResponseWrapper) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(p)
	n1 := int64(n)
	w.BodySize += n1
	w.WriteErr = err
	trace.SpanFromContext(w.Ctx).AddEvent("http.response.write", trace.WithAttributes(wroteBytesKey.Int64(n1)))
	return n, err
}

func (w OTELResponseWrapper) GetSpanAttrs() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.HTTPResponseBodySize(int(w.BodySize)),
		semconv.HTTPResponseStatusCode(w.StatusCode),
	}

	if v := w.Header().Get("Content-Type"); v != "" {
		attrs = append(attrs, attribute.String("http.response.header.content-type", v))
	}
	if v, err := strconv.ParseInt(w.Header().Get("Content-Length"), 10, 64); err != nil {
		attrs = append(attrs, attribute.Int64("http.response.header.content-length", v))
	}

	return attrs
}
