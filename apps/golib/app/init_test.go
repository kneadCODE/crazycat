package app

import (
	"errors"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	type testCase struct {
		givenConfig        Config
		givenConfigErr     error
		givenZap           *zap.SugaredLogger
		givenZapErr        error
		givenSentry        *sentry.Hub
		givenSentryErr     error
		givenNR            *newrelic.Application
		givenNRErr         error
		givenOTELTracer    *sdktrace.TracerProvider
		givenOTELTracerErr error
		expErr             error
		expCfg             Config
		expSentry          bool
		expNR              bool
	}

	tcs := map[string]testCase{
		"ok: basic setup, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:        zap.NewExample().Sugar(),
			givenOTELTracer: sdktrace.NewTracerProvider(),
			expCfg: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: basic setup, staging env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "staging",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:        zap.NewExample().Sugar(),
			givenOTELTracer: sdktrace.NewTracerProvider(),
			expCfg: Config{
				Name:             "name",
				Env:              "staging",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: basic setup, dev env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "development",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:        zap.NewExample().Sugar(),
			givenOTELTracer: sdktrace.NewTracerProvider(),
			expCfg: Config{
				Name:             "name",
				Env:              "development",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: with sentry, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:        zap.NewExample().Sugar(),
			givenSentry:     sentry.NewHub(nil, nil),
			givenOTELTracer: sdktrace.NewTracerProvider(),
			expCfg: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			expSentry: true,
		},
		"err: with cfg error, prod env": {
			givenConfigErr: errors.New("some cfg err"),
			expErr:         errors.New("some cfg err"),
		},
		"err: with zap error, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZapErr: errors.New("some zap err"),
			expErr:      errors.New("some zap err"),
		},
		"err: with sentry error, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:       zap.NewExample().Sugar(),
			givenSentryErr: errors.New("some sentry err"),
			expErr:         errors.New("some sentry err"),
		},
		"ok: with NR, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:        zap.NewExample().Sugar(),
			givenNR:         &newrelic.Application{},
			givenOTELTracer: sdktrace.NewTracerProvider(),
			expCfg: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			expNR: true,
		},
		"err: with NR error, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:   zap.NewExample().Sugar(),
			givenNRErr: errors.New("some NR err"),
			expErr:     errors.New("some NR err"),
		},
		"err: with OTEL error, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenOTELTracerErr: errors.New("some OTEL err"),
			expErr:             errors.New("some OTEL err"),
		},
		"ok: with everything, prod env": {
			givenConfig: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:        zap.NewExample().Sugar(),
			givenSentry:     sentry.NewHub(nil, nil),
			givenNR:         &newrelic.Application{},
			givenOTELTracer: sdktrace.NewTracerProvider(),
			expCfg: Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			expSentry: true,
			expNR:     true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			defer func() {
				newConfigF = newConfig
				newZapF = newZap
				newSentryF = newSentry
				newNewRelicF = newNewRelic
				newOTELProviderF = newOTELProvider
			}()
			newConfigF = func() (Config, error) {
				return tc.givenConfig, tc.givenConfigErr
			}
			newZapF = func(cfg Config) (*zap.SugaredLogger, error) {
				require.Equal(t, tc.givenConfig, cfg)
				return tc.givenZap, tc.givenZapErr
			}
			newSentryF = func(cfg Config) (*sentry.Hub, error) {
				require.Equal(t, tc.givenConfig, cfg)
				return tc.givenSentry, tc.givenSentryErr
			}
			newNewRelicF = func(cfg Config) (*newrelic.Application, error) {
				require.Equal(t, tc.givenConfig, cfg)
				return tc.givenNR, tc.givenNRErr
			}
			newOTELProviderF = func(cfg Config, isSentryEnabled bool) (*sdktrace.TracerProvider, error) {
				require.Equal(t, tc.givenConfig, cfg)
				require.Equal(t, tc.givenSentry != nil, isSentryEnabled)
				// TODO: Validate isSentryEnabled's impact somehow
				return tc.givenOTELTracer, tc.givenOTELTracerErr
			}

			// When:
			ctx, finish, err := Init()

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.Nil(t, finish)
			} else {
				require.Nil(t, err)
				require.NotNil(t, finish)

				cfg := ConfigFromContext(ctx)
				require.EqualValues(t, tc.expCfg, cfg)

				require.Equal(t, tc.givenZap, zapFromContext(ctx))
				require.Equal(t, tc.givenSentry, sentry.GetHubFromContext(ctx))
				require.Equal(t, tc.givenNR, newRelicFromContext(ctx))
				require.Equal(t, tc.givenOTELTracer.Tracer("main"), otelTracerFromContext(ctx))

				finish()
			}
		})
	}
}
