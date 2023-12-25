package otelhttpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts the span. If span is available
func StartSpan(r *http.Request, attrs []attribute.KeyValue) (ctx context.Context, end func(error)) {
	ctx = r.Context()

	traceOpts := []trace.SpanStartOption{
		trace.WithNewRoot(),
		trace.WithSpanKind(trace.SpanKindServer),
	}

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	if s := trace.SpanContextFromContext(ctx); s.IsValid() && s.IsRemote() {
		traceOpts = append(traceOpts, trace.WithLinks(trace.Link{SpanContext: s}))
		// NOTE: We are assuming that all spans within the svc will be child spans and across svcs will be connected
		// via links. Since HTTP server is for inbound requests to the svc, we assume that every call is from another
		// svc and so we only connect the spans via SpanLink instead of as a child.
	}

	ctx, span := internal.GetTracer().Start(
		ctx,
		fmt.Sprintf("%s_%s", r.Method, chi.RouteContext(ctx).RoutePattern()),
		traceOpts...,
	)

	ctx = internal.SetOTELAttrsInContext(ctx, attrs)

	end = func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}

	return
}

// SpanPostProcessing adds post-processing attributes to the span
func SpanPostProcessing(ctx context.Context, rww *ResponseWriterWrapper, rbw *RequestBodyWrapper) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		semconv.HTTPRequestBodySize(int(rbw.bodySize)),
		semconv.HTTPResponseBodySize(int(rww.bodySize)),
		semconv.HTTPResponseStatusCode(rww.statusCode),
	)

	span.SetAttributes(rww.getHeaderAttrs()...)
	if rbw.readErr != nil && rbw.readErr != io.EOF {
		span.SetAttributes(attribute.String("http.request.read_error", rbw.readErr.Error()))
	}
	if rww.writeErr != nil && rww.writeErr != io.EOF {
		span.SetAttributes(attribute.String("http.response.write_error", rww.writeErr.Error()))
	}
}
