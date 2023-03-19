package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Test_newConfig(t *testing.T) {
	type testCase struct {
		givenOSEnvF func(t *testing.T)
		expErr      error
		expCfg      Config
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
			expCfg: Config{
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
			expCfg: Config{
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
			expCfg: Config{
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
			expErr: fmt.Errorf("name is empty: %w", ErrInvalidConfig),
		},
		"err: empty project": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("project is empty: %w", ErrInvalidConfig),
		},
		"err: empty version": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("version is empty: %w", ErrInvalidConfig),
		},
		"err: empty server instance id": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "production"))
			},
			expErr: fmt.Errorf("server instance ID is empty: %w", ErrInvalidConfig),
		},
		"err: empty env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
			},
			expErr: fmt.Errorf("invalid env: %w", ErrInvalidConfig),
		},
		"err: invalid env": {
			givenOSEnvF: func(t *testing.T) {
				require.NoError(t, os.Setenv(string(semconv.ServiceNameKey), "name"))
				require.NoError(t, os.Setenv(string(semconv.ServiceNamespaceKey), "project"))
				require.NoError(t, os.Setenv(string(semconv.ServiceVersionKey), "v1.0.0"))
				require.NoError(t, os.Setenv(string(semconv.ServiceInstanceIDKey), "instance"))
				require.NoError(t, os.Setenv(string(semconv.DeploymentEnvironmentKey), "abcd"))
			},
			expErr: fmt.Errorf("invalid env: %w", ErrInvalidConfig),
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
			cfg, err := newConfig()

			// Then:
			require.Equal(t, tc.expErr, err)
			require.EqualValues(t, tc.expCfg, cfg)
		})
	}
}
