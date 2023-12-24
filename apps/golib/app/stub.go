package app

import (
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	"go.opentelemetry.io/otel"
)

// stubs for testing
var newZapStub = internal.NewZap
var newOTELResourceFromEnvStub = internal.NewOTELResourceFromEnv
var newOTELPropagatorStub = internal.NewOTELPropagator
var newOTELTraceProviderStub = internal.NewOTELTraceProvider
var newOTELMeterProviderStub = internal.NewOTELMeterProvider
var setOTELTextMapPropagatorStub = otel.SetTextMapPropagator
var setOTELTracerProviderStub = otel.SetTracerProvider
var setOTELMeterProviderStub = otel.SetMeterProvider
