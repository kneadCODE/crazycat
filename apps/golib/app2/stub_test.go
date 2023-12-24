package app2

import (
	"github.com/kneadCODE/crazycat/apps/golib/app2/internal"
	"go.opentelemetry.io/otel"
)

func resetStubs() {
	newZapStub = internal.NewZap
	newOTELResourceFromEnvStub = internal.NewOTELResourceFromEnv
	newOTELPropagatorStub = internal.NewOTELPropagator
	newOTELTraceProviderStub = internal.NewOTELTraceProvider
	newOTELMeterProviderStub = internal.NewOTELMeterProvider
	setOTELTextMapPropagatorStub = otel.SetTextMapPropagator
	setOTELTracerProviderStub = otel.SetTracerProvider
	setOTELMeterProviderStub = otel.SetMeterProvider
}
