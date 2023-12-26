package httpserver

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func TestNew(t *testing.T) {
	defer otel.SetMeterProvider(nil)
	defer otel.SetTracerProvider(nil)
	otel.SetMeterProvider(metricnoop.NewMeterProvider())
	otel.SetTracerProvider(tracenoop.NewTracerProvider())

	rtr := Router{}
	handler, err := rtr.Handler()
	require.NoError(t, err)

	type testCase struct {
		givenOpts []ServerOption
		expErr    error
		expSrv    *Server
	}
	tcs := map[string]testCase{
		"default": {
			expSrv: &Server{
				srv: &http.Server{
					Handler:      handler,
					Addr:         ":9000",
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				gracefulShutdownTimeout: 10 * time.Second,
			},
		},
		"override port": {
			givenOpts: []ServerOption{WithServerPort(3000)},
			expSrv: &Server{
				srv: &http.Server{
					Handler:      handler,
					Addr:         ":3000",
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  120 * time.Second,
				},
				gracefulShutdownTimeout: 10 * time.Second,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given && When:
			srv, err := New(context.Background(), rtr, tc.givenOpts...)

			// Then:
			require.Equal(t, tc.expErr, err)
			if tc.expSrv != nil {
				require.EqualValues(t, tc.expSrv.srv.Addr, srv.srv.Addr)
				require.EqualValues(t, tc.expSrv.srv.ReadTimeout, srv.srv.ReadTimeout)
				require.EqualValues(t, tc.expSrv.srv.ReadHeaderTimeout, srv.srv.ReadHeaderTimeout)
				require.EqualValues(t, tc.expSrv.srv.WriteTimeout, srv.srv.WriteTimeout)
				require.EqualValues(t, tc.expSrv.srv.MaxHeaderBytes, srv.srv.MaxHeaderBytes)
				require.EqualValues(t, tc.expSrv.srv.ErrorLog, srv.srv.ErrorLog)
				require.NotNil(t, srv.srv.BaseContext)
				require.Nil(t, srv.srv.ConnContext)
				require.EqualValues(t, tc.expSrv.gracefulShutdownTimeout, srv.gracefulShutdownTimeout)
			} else {
				require.Nil(t, srv)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	defer otel.SetMeterProvider(nil)
	defer otel.SetTracerProvider(nil)
	otel.SetMeterProvider(metricnoop.NewMeterProvider())
	otel.SetTracerProvider(tracenoop.NewTracerProvider())

	// Given:
	srv, err := New(context.Background(), Router{})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// When:
	go func() {
		defer wg.Done()
		err = srv.Start(ctx)
	}()

	// Then:
	time.Sleep(2 * time.Second)

	cancel()

	wg.Wait()

	require.NoError(t, err)
}
