package otelhttpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(r *http.Request) (ctx context.Context, end func(error)) {
	ctx = r.Context()

	attrs := extractOTELAttrsFromReq(r)

	traceOpts := []trace.SpanStartOption{
		trace.WithNewRoot(),
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	}

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	if s := trace.SpanContextFromContext(ctx); s.IsValid() && s.IsRemote() {
		traceOpts = append(traceOpts, trace.WithLinks(trace.Link{SpanContext: s}))
	}

	ctx, span := internal.GetTracer().Start(
		ctx,
		fmt.Sprintf("%s_%s", r.Method, chi.RouteContext(ctx).RoutePattern()),
		traceOpts...,
	)

	ctx = internal.SetOTELAttrsInContext(ctx, attrs)

	return ctx, func(err error) {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}
}

func PostHandle(ctx context.Context, rw *OTELResponseWrapper, bw *OTELRequestBodyWrapper) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		semconv.HTTPRequestBodySize(int(bw.BodySize)),
		semconv.HTTPResponseBodySize(int(rw.BodySize)),
		semconv.HTTPResponseStatusCode(rw.StatusCode),
	)
	if bw.ReadErr != nil {
		span.SetAttributes(readErrorKey.String(bw.ReadErr.Error()))
	}
	if rw.WriteErr != nil {
		span.SetAttributes(writeErrorKey.String(rw.WriteErr.Error()))
	}
}
