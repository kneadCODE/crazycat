package otelhttpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestRequestBodyWrapper_Read(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("testing")
	validJSONFileData, err := os.ReadFile("./testdata/valid.json")
	require.NoError(t, err)
	invalidJSONFileData, err := os.ReadFile("./testdata/invalid.json")
	require.NoError(t, err)
	corruptedJSONFileData, err := os.ReadFile("./testdata/corrupted.json")
	require.NoError(t, err)

	type body struct {
		StrField string `json:"str_field"`
		IntField int    `json:"int_field"`
	}
	var invalidBody body
	jsonUnmarshallErr := json.Unmarshal(invalidJSONFileData, &invalidBody)

	type testCase struct {
		givenCtx   func() context.Context
		givenReq   func() *http.Request
		expErr     error
		expBytes   int64
		expReadErr error
	}
	tcs := map[string]testCase{
		"GET: no body": {
			givenCtx: func() context.Context { return context.Background() },
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			expErr:     io.EOF,
			expReadErr: io.EOF,
		},
		"GET: no body, with span": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			expErr:     io.EOF,
			expReadErr: io.EOF,
		},
		"GET: with body": {
			givenCtx: func() context.Context { return context.Background() },
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(validJSONFileData))
			},
			expBytes: int64(len(validJSONFileData)),
		},
		"GET: with body, with span": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(validJSONFileData))
			},
			expBytes: int64(len(validJSONFileData)),
		},
		"GET: with body, with span, but invalid body": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(invalidJSONFileData))
			},
			expErr:   jsonUnmarshallErr,
			expBytes: int64(len(invalidJSONFileData)),
		},
		"GET: with body, with span, but corrupted body": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(corruptedJSONFileData))
			},
			expErr:     io.ErrUnexpectedEOF,
			expBytes:   int64(len(corruptedJSONFileData)),
			expReadErr: io.EOF,
		},
		"POST: no body": {
			givenCtx: func() context.Context {
				return context.Background()
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", nil)
			},
			expErr:     io.EOF,
			expReadErr: io.EOF,
		},
		"POST: no body, with span": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", nil)
			},
			expErr:     io.EOF,
			expReadErr: io.EOF,
		},
		"POST: with body": {
			givenCtx: func() context.Context { return context.Background() },
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(validJSONFileData))
			},
			expBytes: int64(len(validJSONFileData)),
		},
		"POST: with body, with span": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(validJSONFileData))
			},
			expBytes: int64(len(validJSONFileData)),
		},
		"POST: with body, with span, but invalid body": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidJSONFileData))
			},
			expErr:   jsonUnmarshallErr,
			expBytes: int64(len(invalidJSONFileData)),
		},
		"POST: with body, with span, but corrupted body": {
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			givenReq: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(corruptedJSONFileData))
			},
			expErr:     io.ErrUnexpectedEOF,
			expBytes:   int64(len(corruptedJSONFileData)),
			expReadErr: io.EOF,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			ctx := tc.givenCtx()

			r := tc.givenReq()

			rbw := RequestBodyWrapper{
				ReadCloser: r.Body,
				Ctx:        ctx,
			}
			r.Body = &rbw

			// When:
			var b body
			err := json.NewDecoder(r.Body).Decode(&b)

			// Then:
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expBytes, rbw.bodySize)
			require.Equal(t, tc.expReadErr, rbw.readErr)
		})
	}
}

func TestExtractAttrsFromReq(t *testing.T) {
	type testCase struct {
		givenReq func() *http.Request
		expAttrs []attribute.KeyValue
	}
	tcs := map[string]testCase{
		"plain req": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd", nil)
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd"),
				semconv.URLPath("/abcd"),
				semconv.URLQuery(""),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal(""),
				semconv.ServerAddress(""),
				semconv.ClientAddress(""),
				semconv.NetProtocolVersion("1.1"),
			},
		},
		"req with ua": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd"),
				semconv.URLPath("/abcd"),
				semconv.URLQuery(""),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress(""),
				semconv.ClientAddress(""),
				semconv.NetProtocolVersion("1.1"),
			},
		},
		"req with ua, route pattern": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd/123", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd/{uid}"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd/{uid}"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd/123"),
				semconv.URLPath("/abcd/123"),
				semconv.URLQuery(""),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress(""),
				semconv.ClientAddress(""),
				semconv.NetProtocolVersion("1.1"),
			},
		},
		"req with ua, route pattern, query param": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd/123?a=b&c=d", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd/{uid}"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd/{uid}"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd/123?a=b&c=d"),
				semconv.URLPath("/abcd/123"),
				semconv.URLQuery("a=b&c=d"),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress(""),
				semconv.ClientAddress(""),
				semconv.NetProtocolVersion("1.1"),
			},
		},
		"req with ua, route pattern, query param, content-type": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd/123?a=b&c=d", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				r.Header.Set("Content-Type", "application/json")
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd/{uid}"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd/{uid}"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd/123?a=b&c=d"),
				semconv.URLPath("/abcd/123"),
				semconv.URLQuery("a=b&c=d"),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress(""),
				semconv.ClientAddress(""),
				semconv.NetProtocolVersion("1.1"),
				attribute.String("http.request.header.content-type", "application/json"),
			},
		},
		"req with ua, route pattern, query param, content-type, content-length": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd/123?a=b&c=d", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				r.Header.Set("Content-Type", "application/json")
				r.Header.Set("Content-Length", "1234")
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd/{uid}"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd/{uid}"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd/123?a=b&c=d"),
				semconv.URLPath("/abcd/123"),
				semconv.URLQuery("a=b&c=d"),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress(""),
				semconv.ClientAddress(""),
				semconv.NetProtocolVersion("1.1"),
				attribute.String("http.request.header.content-type", "application/json"),
				attribute.Int64("http.request.header.content-length", 1234),
			},
		},
		"req with ua, route pattern, query param, content-type, content-length, x-forwarded-for": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd/123?a=b&c=d", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				r.Header.Set("Content-Type", "application/json")
				r.Header.Set("Content-Length", "1234")
				r.Header.Set("X-Forwarded-For", "0.0.0.0,1.1.1.1,2.2.2.2")
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd/{uid}"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd/{uid}"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd/123?a=b&c=d"),
				semconv.URLPath("/abcd/123"),
				semconv.URLQuery("a=b&c=d"),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress(""),
				semconv.ClientAddress("0.0.0.0,1.1.1.1,2.2.2.2"),
				semconv.NetProtocolVersion("1.1"),
				attribute.String("http.request.header.content-type", "application/json"),
				attribute.Int64("http.request.header.content-length", 1234),
				attribute.String("http.request.header.x-forwarded-for", "0.0.0.0,1.1.1.1,2.2.2.2"),
			},
		},
		"req with ua, route pattern, query param, content-type, content-length, x-forwarded-for, srv addr and port": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abcd/123?a=b&c=d", nil)
				r.Header.Set("User-Agent", "crazycat/v0.0.0")
				r.Header.Set("Content-Type", "application/json")
				r.Header.Set("Content-Length", "1234")
				r.Header.Set("X-Forwarded-For", "0.0.0.0,1.1.1.1,2.2.2.2")
				r.URL.Host = "localhost:3000"
				rctx := chi.NewRouteContext()
				rctx.RoutePatterns = []string{"/abcd/{uid}"}
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
				return r
			},
			expAttrs: []attribute.KeyValue{
				semconv.HTTPRequestMethodKey.String("GET"),
				semconv.HTTPRoute("/abcd/{uid}"),
				semconv.NetProtocolName("http"),
				semconv.URLScheme("http"),
				semconv.URLFull("/abcd/123?a=b&c=d"),
				semconv.URLPath("/abcd/123"),
				semconv.URLQuery("a=b&c=d"),
				semconv.URLScheme("http"),
				semconv.URLFragment(""),
				semconv.UserAgentOriginal("crazycat/v0.0.0"),
				semconv.ServerAddress("localhost"),
				semconv.ClientAddress("0.0.0.0,1.1.1.1,2.2.2.2"),
				semconv.ServerPort(3000),
				semconv.NetProtocolVersion("1.1"),
				attribute.String("http.request.header.content-type", "application/json"),
				attribute.Int64("http.request.header.content-length", 1234),
				attribute.String("http.request.header.x-forwarded-for", "0.0.0.0,1.1.1.1,2.2.2.2"),
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			r := tc.givenReq()

			// When:
			attrs := ExtractAttrsFromReq(r)

			// Then:
			require.EqualValues(t, tc.expAttrs, attrs)
		})
	}
}
