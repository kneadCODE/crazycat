package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	type testCase struct {
		givenConfig        config.Config
		givenConfigErr     error
		givenZap           *zap.SugaredLogger
		givenZapErr        error
		givenSentry        *sentry.Hub
		givenSentryErr     error
		givenNR            *newrelic.Application
		givenNRErr         error
		givenOTELShutdownF func(context.Context) error
		givenOTELTracerErr error
		expErr             error
		expCfg             config.Config
		expSentry          bool
		expNR              bool
	}

	tcs := map[string]testCase{
		"ok: basic setup, prod env": {
			givenConfig: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenOTELShutdownF: func(context.Context) error { return nil },
			expCfg: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: basic setup, staging env": {
			givenConfig: config.Config{
				Name:             "name",
				Env:              "staging",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenOTELShutdownF: func(context.Context) error { return nil },
			expCfg: config.Config{
				Name:             "name",
				Env:              "staging",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: basic setup, dev env": {
			givenConfig: config.Config{
				Name:             "name",
				Env:              "development",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenOTELShutdownF: func(context.Context) error { return nil },
			expCfg: config.Config{
				Name:             "name",
				Env:              "development",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: with sentry, prod env": {
			givenConfig: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenSentry:        sentry.NewHub(nil, nil),
			givenOTELShutdownF: func(context.Context) error { return nil },
			expCfg: config.Config{
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
			givenConfig: config.Config{
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
			givenConfig: config.Config{
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
			givenConfig: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenNR:            &newrelic.Application{},
			givenOTELShutdownF: func(context.Context) error { return nil },
			expCfg: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			expNR: true,
		},
		"err: with NR error, prod env": {
			givenConfig: config.Config{
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
			givenConfig: config.Config{
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
			givenConfig: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
			givenZap:           zap.NewExample().Sugar(),
			givenSentry:        sentry.NewHub(nil, nil),
			givenNR:            &newrelic.Application{},
			givenOTELShutdownF: func(context.Context) error { return nil },
			expCfg: config.Config{
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
				newConfigFromEnvF = newConfigFromEnv
				newZapF = internal.NewZap
				newSentryF = internal.NewSentryHub
				newNewRelicF = internal.NewNewRelicApp
				initOTELF = internal.InitOTEL
			}()
			newConfigFromEnvF = func() (config.Config, error) {
				return tc.givenConfig, tc.givenConfigErr
			}
			newZapF = func(cfg config.Config) (*zap.SugaredLogger, error) {
				require.Equal(t, tc.givenConfig, cfg)
				return tc.givenZap, tc.givenZapErr
			}
			newSentryF = func(cfg config.Config) (*sentry.Hub, error) {
				require.Equal(t, tc.givenConfig, cfg)
				return tc.givenSentry, tc.givenSentryErr
			}
			newNewRelicF = func(cfg config.Config) (*newrelic.Application, error) {
				require.Equal(t, tc.givenConfig, cfg)
				return tc.givenNR, tc.givenNRErr
			}
			initOTELF = func(cfg config.Config, isSentryEnabled bool) (func(context.Context) error, error) {
				require.Equal(t, tc.givenConfig, cfg)
				// TODO: Validate isSentryEnabled's impact somehow
				return tc.givenOTELShutdownF, tc.givenOTELTracerErr
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

				finish()
			}
		})
	}
}

func Test_NewConfigFromEnv(t *testing.T) {
	type testCase struct {
		givenOSEnvF func(t *testing.T)
		expErr      error
		expCfg      config.Config
	}

	tcs := map[string]testCase{
		"ok: prod env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expCfg: config.Config{
				Name:             "name",
				Env:              "production",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: staging env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "staging"))
			},
			expCfg: config.Config{
				Name:             "name",
				Env:              "staging",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"ok: dev env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "development"))
			},
			expCfg: config.Config{
				Name:             "name",
				Env:              "development",
				Project:          "project",
				Version:          "v1.0.0",
				ServerInstanceID: "instance",
			},
		},
		"err: empty name": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("name is empty: %w", config.ErrInvalidConfig),
		},
		"err: empty project": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("project is empty: %w", config.ErrInvalidConfig),
		},
		"err: empty version": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("version is empty: %w", config.ErrInvalidConfig),
		},
		"err: empty server instance id": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("server instance ID is empty: %w", config.ErrInvalidConfig),
		},
		"err: empty env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
			},
			expErr: fmt.Errorf("invalid env: %w", config.ErrInvalidConfig),
		},
		"err: invalid env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "abcd"))
			},
			expErr: fmt.Errorf("invalid env: %w", config.ErrInvalidConfig),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			defer os.Unsetenv(string(semconv.ServiceNameKey))
			defer os.Unsetenv(string(semconv.ServiceNamespaceKey))
			defer os.Unsetenv(string(semconv.ServiceVersionKey))
			defer os.Unsetenv(string(semconv.ServiceInstanceIDKey))
			defer os.Unsetenv(string(semconv.DeploymentEnvironmentKey))

			tc.givenOSEnvF(t)

			// When:
			cfg, err := newConfigFromEnv()

			// Then:
			require.Equal(t, tc.expErr, err)
			require.EqualValues(t, tc.expCfg, cfg)
		})
	}
}
