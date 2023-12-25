package otelhttpserver

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type RequestBodyWrapper struct {
	io.ReadCloser

	Ctx context.Context

	bodySize int64
	readErr  error
}

func (w *RequestBodyWrapper) Read(b []byte) (int, error) {
	n, err := w.ReadCloser.Read(b)
	n1 := int64(n)
	w.bodySize += n1
	w.readErr = err
	trace.SpanFromContext(w.Ctx).AddEvent("http.request.read", trace.WithAttributes(
		attribute.Int64("http.request.read_bytes", n1),
	))
	return n, err
}

func ExtractOTELAttrsFromReq(r *http.Request) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(r.Method),
		semconv.HTTPRoute(chi.RouteContext(r.Context()).RoutePattern()),

		semconv.NetProtocolName("http"),
		// Protocol Version filled in later

		semconv.URLScheme("http"),
		semconv.URLFull(r.RequestURI),
		semconv.URLPath(r.URL.Path),
		semconv.URLQuery(r.URL.RawQuery),
		semconv.URLScheme("http"),
		semconv.URLFragment(r.URL.Fragment),

		semconv.UserAgentOriginal(r.UserAgent()),

		semconv.ServerAddress(r.URL.Host),
		// Server port filled in later

		semconv.ClientAddress(r.Header.Get("X-Forwarded-For")),
	}

	if v, err := strconv.Atoi(r.URL.Port()); err != nil {
		attrs = append(attrs, semconv.ServerPort(v))
	}

	_, protoVersion, _ := strings.Cut(r.Proto, "/")
	attrs = append(attrs, semconv.NetProtocolVersion(protoVersion))

	if v := r.Header.Get("Content-Type"); v != "" {
		attrs = append(attrs, attribute.String("http.request.header.content-type", v))
	}
	if v, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64); err != nil {
		attrs = append(attrs, attribute.Int64("http.request.header.content-length", v))
	}
	if v := r.Header.Values("X-Forwarded-For"); len(v) != 0 {
		attrs = append(attrs, attribute.StringSlice("http.request.header.x-forwarded-for", v))
	}

	return attrs
}
