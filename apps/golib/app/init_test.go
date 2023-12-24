package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	type testCase struct {
		givenSentryEnabled                    bool
		mockRes                               *resource.Resource
		mockResErr                            error
		mockDebugMode                         bool
		mockPropagators                       propagation.TextMapPropagator
		mockZap                               *zap.Logger
		mockZapErr                            error
		mockTraceProv                         *sdktrace.TracerProvider
		mockTraceProvErr                      error
		mockMeterProv                         *sdkmetric.MeterProvider
		mockMeterProvErr                      error
		expCfg                                Config
		expErr                                error
		expNewOTELResourceFromEnvStubCalled   bool
		expNewZapStubCalled                   bool
		expNewOTELPropagatorStubCalled        bool
		expNewOTELTraceProviderStubCalled     bool
		expNewOTELMeterProviderStubCalled     bool
		expSetOTELTextMapPropagatorStubCalled bool
		expSetOTELTracerProviderStubCalled    bool
		expSetOTELMeterProviderStubCalled     bool
	}

	tcs := map[string]testCase{
		"res err": {
			mockResErr:                          errors.New("some err"),
			expErr:                              errors.New("some err"),
			expNewOTELResourceFromEnvStubCalled: true,
		},
		"invalid env": {
			mockRes:                             resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("dev")),
			expErr:                              errors.New("invalid env: [dev]"),
			expNewOTELResourceFromEnvStubCalled: true,
		},
		"zap err": {
			mockRes:                               resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development")),
			mockDebugMode:                         true,
			mockPropagators:                       propagation.NewCompositeTextMapPropagator(),
			mockZapErr:                            errors.New("some err"),
			expErr:                                errors.New("some err"),
			expNewOTELResourceFromEnvStubCalled:   true,
			expNewOTELPropagatorStubCalled:        true,
			expSetOTELTextMapPropagatorStubCalled: true,
			expNewZapStubCalled:                   true,
		},
		"trace provider err": {
			mockRes:                               resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development")),
			mockDebugMode:                         true,
			mockPropagators:                       propagation.NewCompositeTextMapPropagator(),
			mockZap:                               zap.NewExample(),
			mockTraceProvErr:                      errors.New("some err"),
			expErr:                                errors.New("some err"),
			expNewOTELResourceFromEnvStubCalled:   true,
			expNewOTELPropagatorStubCalled:        true,
			expSetOTELTextMapPropagatorStubCalled: true,
			expNewZapStubCalled:                   true,
			expNewOTELTraceProviderStubCalled:     true,
		},
		"meter provider err": {
			mockRes:                               resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development")),
			mockDebugMode:                         true,
			mockPropagators:                       propagation.NewCompositeTextMapPropagator(),
			mockZap:                               zap.NewExample(),
			mockTraceProv:                         sdktrace.NewTracerProvider(),
			mockMeterProvErr:                      errors.New("some err"),
			expErr:                                errors.New("some err"),
			expNewOTELResourceFromEnvStubCalled:   true,
			expNewOTELPropagatorStubCalled:        true,
			expSetOTELTextMapPropagatorStubCalled: true,
			expNewZapStubCalled:                   true,
			expNewOTELTraceProviderStubCalled:     true,
			expSetOTELTracerProviderStubCalled:    true,
			expNewOTELMeterProviderStubCalled:     true,
		},
		"success": {
			mockRes:                               resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development")),
			mockDebugMode:                         true,
			mockPropagators:                       propagation.NewCompositeTextMapPropagator(),
			mockZap:                               zap.NewExample(),
			mockTraceProv:                         sdktrace.NewTracerProvider(),
			mockMeterProv:                         sdkmetric.NewMeterProvider(),
			expCfg:                                Config{Env: EnvDev, res: resource.NewWithAttributes(semconv.SchemaURL, semconv.DeploymentEnvironment("development"))},
			expNewOTELResourceFromEnvStubCalled:   true,
			expNewOTELPropagatorStubCalled:        true,
			expSetOTELTextMapPropagatorStubCalled: true,
			expNewZapStubCalled:                   true,
			expNewOTELTraceProviderStubCalled:     true,
			expSetOTELTracerProviderStubCalled:    true,
			expNewOTELMeterProviderStubCalled:     true,
			expSetOTELMeterProviderStubCalled:     true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			defer resetStubs()
			var newOTELResourceFromEnvStubCalled bool
			newOTELResourceFromEnvStub = func(ctx context.Context) (*resource.Resource, error) {
				newOTELResourceFromEnvStubCalled = true
				return tc.mockRes, tc.mockResErr
			}
			var newZapStubCalled bool
			newZapStub = func(debugMode bool, res *resource.Resource) (*zap.Logger, error) {
				newZapStubCalled = true
				require.Equal(t, tc.mockRes, res)
				require.Equal(t, tc.mockDebugMode, debugMode)
				return tc.mockZap, tc.mockZapErr
			}
			var newOTELPropagatorStubCalled bool
			newOTELPropagatorStub = func(isSentryEnabled bool) propagation.TextMapPropagator {
				newOTELPropagatorStubCalled = true
				require.Equal(t, tc.givenSentryEnabled, isSentryEnabled)
				return tc.mockPropagators
			}
			var newOTELTraceProviderStubCalled bool
			newOTELTraceProviderStub = func(res *resource.Resource, isSentryEnabled bool) (*sdktrace.TracerProvider, error) {
				newOTELTraceProviderStubCalled = true
				require.Equal(t, tc.givenSentryEnabled, isSentryEnabled)
				require.Equal(t, tc.mockRes, res)
				return tc.mockTraceProv, tc.mockTraceProvErr
			}
			var newOTELMeterProviderStubCalled bool
			newOTELMeterProviderStub = func(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
				newOTELMeterProviderStubCalled = true
				require.Equal(t, tc.mockRes, res)
				return tc.mockMeterProv, tc.mockMeterProvErr
			}
			var setOTELTextMapPropagatorStubCalled bool
			setOTELTextMapPropagatorStub = func(propagator propagation.TextMapPropagator) {
				setOTELTextMapPropagatorStubCalled = true
				require.Equal(t, tc.mockPropagators, propagator)
			}
			var setOTELTracerProviderStubCalled bool
			setOTELTracerProviderStub = func(tp trace.TracerProvider) {
				setOTELTracerProviderStubCalled = true
				require.Equal(t, tc.mockTraceProv, tp)
			}
			var setOTELMeterProviderStubCalled bool
			setOTELMeterProviderStub = func(mp metric.MeterProvider) {
				setOTELMeterProviderStubCalled = true
				require.Equal(t, tc.mockMeterProv, mp)
			}

			// When:
			ctx, finish, err := Init()

			// Then:
			require.Equal(t, tc.expNewOTELResourceFromEnvStubCalled, newOTELResourceFromEnvStubCalled)
			require.Equal(t, tc.expNewZapStubCalled, newZapStubCalled)
			require.Equal(t, tc.expNewOTELPropagatorStubCalled, newOTELPropagatorStubCalled)
			require.Equal(t, tc.expNewOTELTraceProviderStubCalled, newOTELTraceProviderStubCalled)
			require.Equal(t, tc.expNewOTELMeterProviderStubCalled, newOTELMeterProviderStubCalled)
			require.Equal(t, tc.expSetOTELTextMapPropagatorStubCalled, setOTELTextMapPropagatorStubCalled)
			require.Equal(t, tc.expSetOTELTracerProviderStubCalled, setOTELTracerProviderStubCalled)
			require.Equal(t, tc.expSetOTELMeterProviderStubCalled, setOTELMeterProviderStubCalled)

			if tc.expErr != nil {
				require.Equal(t, tc.expErr, err)
				require.Nil(t, finish)
				require.EqualValues(t, tc.expCfg, ConfigFromContext(ctx))
			} else {
				require.Nil(t, err)
				require.NotNil(t, finish)

				cfg := ConfigFromContext(ctx)
				require.EqualValues(t, tc.expCfg, cfg)
				require.Equal(t, tc.mockZap, zapFromContext(ctx))

				finish()
			}
		})
	}
}
