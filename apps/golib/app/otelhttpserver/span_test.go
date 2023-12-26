package otelhttpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestStartSpan(t *testing.T) {
	defer otel.SetTracerProvider(nil)
	defer otel.SetTextMapPropagator(nil)

	// Given: no provider, no span, no attrs
	r := httptest.NewRequest(http.MethodGet, "/abc", nil)
	rctx := chi.NewRouteContext()
	rctx.RoutePatterns = []string{"/abc"}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// When:
	ctx, end := StartSpan(r, nil)
	end(nil)

	// Then:
	require.NotNil(t, ctx)

	// Given: no span, no attrs
	otel.SetTracerProvider(noop.NewTracerProvider())

	// When:
	ctx, end = StartSpan(r, nil)
	end(nil)

	// Then:
	require.NotNil(t, ctx)

	// When: no span, with attrs, no err
	ctx, end = StartSpan(r, []attribute.KeyValue{attribute.String("k1", "v1")})
	end(nil)

	// Then:
	require.NotNil(t, ctx)

	// When: no span, with attrs, with err
	ctx, end = StartSpan(r, []attribute.KeyValue{attribute.String("k1", "v1")})
	end(errors.New("err"))

	// Then:
	require.NotNil(t, ctx)

	// Given:
	ctx, _ = sdktrace.NewTracerProvider().Tracer("start").Start(r.Context(), "origin")
	prop := internal.NewOTELPropagator(false)
	otel.SetTextMapPropagator(prop)
	prop.Inject(ctx, propagation.HeaderCarrier(r.Header))

	// When:
	ctx, end = StartSpan(r, []attribute.KeyValue{attribute.String("k1", "v1")})
	end(nil)

	// Then:
	require.NotNil(t, ctx)
}

func TestSpanPostProcessing(t *testing.T) {
	// Given:
	ctx := context.Background()
	rbw := &RequestBodyWrapper{}
	rww := &ResponseWriterWrapper{ResponseWriter: httptest.NewRecorder()}

	// When && Then:
	SpanPostProcessing(ctx, rww, rbw)

	// Given:
	rbw.bodySize = 123
	rww.bodySize = 123
	rww.Header().Add("Content-Type", "application/json")
	rww.Header().Add("Content-Length", "1234")

	// When && Then:
	SpanPostProcessing(ctx, rww, rbw)

	// Given:
	rbw.readErr = errors.New("some err")
	rww.writeErr = errors.New("some err")

	// Given:
	ctx, _ = noop.NewTracerProvider().Tracer("testing").Start(ctx, "span")

	// When && Then:
	SpanPostProcessing(ctx, rww, rbw)
}
