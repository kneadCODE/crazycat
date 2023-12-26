package otelhttpserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestResponseWriterWrapper_WriteHeader(t *testing.T) {
	type testCase struct {
		givenStatus int
		expStatus   int
	}
	tcs := map[string]testCase{
		"200": {
			givenStatus: http.StatusOK,
			expStatus:   http.StatusOK,
		},
		"400": {
			givenStatus: http.StatusBadRequest,
			expStatus:   http.StatusBadRequest,
		},
		"500": {
			givenStatus: http.StatusInternalServerError,
			expStatus:   http.StatusInternalServerError,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			ctx := context.Background()

			w := httptest.NewRecorder()

			rbw := ResponseWriterWrapper{
				ResponseWriter: w,
				Ctx:            ctx,
			}

			// When:
			rbw.WriteHeader(tc.givenStatus)

			// Then:
			require.Equal(t, tc.expStatus, w.Result().StatusCode)
			require.NoError(t, rbw.writeErr)
			require.Equal(t, int64(0), rbw.bodySize)
			require.Equal(t, tc.expStatus, rbw.statusCode)
		})
	}
}

func TestResponseWriterWrapper_Write(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("testing")
	validJSONFileData, err := os.ReadFile("./testdata/valid.json")
	require.NoError(t, err)

	type testCase struct {
		givenPayload []byte
		givenCtx     func() context.Context
		expErr       error
		expBytes     int64
		expStatus    int
		expWriteErr  error
		expBody      []byte
	}
	tcs := map[string]testCase{
		"success json": {
			givenPayload: validJSONFileData,
			givenCtx:     func() context.Context { return context.Background() },
			expBytes:     int64(len(validJSONFileData)),
			expStatus:    http.StatusOK,
			expBody:      validJSONFileData,
		},
		"success, empty payload": {
			givenPayload: []byte(""),
			givenCtx:     func() context.Context { return context.Background() },
			expBytes:     0,
			expStatus:    http.StatusOK,
			expBody:      []byte(""),
		},
		"success, non-json": {
			givenPayload: []byte("abcd"),
			givenCtx:     func() context.Context { return context.Background() },
			expBytes:     4,
			expStatus:    http.StatusOK,
			expBody:      []byte("abcd"),
		},
		"success json, with span": {
			givenPayload: validJSONFileData,
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			expBytes:  int64(len(validJSONFileData)),
			expStatus: http.StatusOK,
			expBody:   validJSONFileData,
		},
		"success, empty payload, with span": {
			givenPayload: []byte(""),
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			expBytes:  0,
			expStatus: http.StatusOK,
			expBody:   []byte(""),
		},
		"success, non-json, with span": {
			givenPayload: []byte("abcd"),
			givenCtx: func() context.Context {
				ctx, _ := tracer.Start(context.Background(), "test")
				return ctx
			},
			expBytes:  4,
			expStatus: http.StatusOK,
			expBody:   []byte("abcd"),
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			ctx := tc.givenCtx()

			w := httptest.NewRecorder()

			rbw := ResponseWriterWrapper{
				ResponseWriter: w,
				Ctx:            ctx,
			}

			// When:
			b, err := rbw.Write(tc.givenPayload)

			// Then:
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expBytes, rbw.bodySize)
			require.Equal(t, tc.expStatus, rbw.statusCode)
			require.Equal(t, tc.expWriteErr, rbw.writeErr)
			require.Equal(t, rbw.bodySize, int64(b))
			require.Equal(t, tc.expStatus, w.Result().StatusCode)
			if tc.expErr == nil {
				result, err := io.ReadAll(w.Body)
				require.NoError(t, err)
				require.Equal(t, tc.expBody, result)

				result, err = io.ReadAll(w.Result().Body)
				require.NoError(t, err)
				require.Equal(t, tc.expBody, result)
			}
		})
	}
}
