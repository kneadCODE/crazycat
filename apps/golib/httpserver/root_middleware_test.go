package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Test_newRootMiddleware(t *testing.T) {
	defer otel.SetMeterProvider(nil)

	// Given && When:
	m, err := newRootMiddleware()

	// Then:
	require.NoError(t, err)
	require.NotNil(t, m)

	// Given:
	otel.SetMeterProvider(sdkmetric.NewMeterProvider())

	// When:
	m, err = newRootMiddleware()

	// Then:
	require.NoError(t, err)
	require.NotNil(t, m)
}

func Test_rootMiddleware_serveHTTP(t *testing.T) {
	defer otel.SetTracerProvider(nil)
	defer otel.SetMeterProvider(nil)

	otel.SetMeterProvider(sdkmetric.NewMeterProvider())
	otel.SetTracerProvider(sdktrace.NewTracerProvider())

	m, err := newRootMiddleware()
	require.NoError(t, err)

	type testCase struct {
		givenReq  func() *http.Request
		givenHF   func(http.ResponseWriter, *http.Request)
		expStatus int
	}
	tcs := map[string]testCase{
		"panic": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abc", nil)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
				return r
			},
			givenHF:   func(http.ResponseWriter, *http.Request) { panic("some err") },
			expStatus: http.StatusInternalServerError,
		},
		"GET": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/abc", nil)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
				return r
			},
			givenHF: func(_ http.ResponseWriter, r *http.Request) {
				v, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, "", string(v))
			},
			expStatus: http.StatusOK,
		},
		"POST": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/abc", nil)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
				return r
			},
			givenHF: func(_ http.ResponseWriter, r *http.Request) {
				v, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, "", string(v))
			},
			expStatus: http.StatusOK,
		},
		"POST with non-json body": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/abc", bytes.NewReader([]byte("abc")))
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
				return r
			},
			givenHF: func(_ http.ResponseWriter, r *http.Request) {
				v, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, "abc", string(v))
			},
			expStatus: http.StatusOK,
		},
		"POST with json body": {
			givenReq: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/abc", bytes.NewReader([]byte(`{"key":"val"}`)))
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
				return r
			},
			givenHF: func(_ http.ResponseWriter, r *http.Request) {
				var res struct {
					Key string `json:"key"`
				}
				err := json.NewDecoder(r.Body).Decode(&res)
				require.NoError(t, err)
				require.Equal(t, "val", res.Key)
			},
			expStatus: http.StatusOK,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			r := tc.givenReq()
			w := httptest.NewRecorder()

			// When:
			m(http.HandlerFunc(tc.givenHF)).ServeHTTP(w, r)

			// Then:
			require.Equal(t, tc.expStatus, w.Code)
		})
	}
}
