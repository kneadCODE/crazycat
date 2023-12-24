package app2

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

// Init initializes the app and returns
func Init() (ctx context.Context, shutdown func(), err error) {
	ctx = context.Background()
	basicLogger := log.New(os.Stdout, "", log.LstdFlags)

	basicLogger.Println("Starting App initialization...")

	basicLogger.Println("Initializing Config from env...")
	cfg, err := newConfigFromEnv(ctx)
	if err != nil {
		return
	}
	setOTELTextMapPropagatorStub(newOTELPropagatorStub(false)) // TODO: Deal with the bool
	basicLogger.Println("Config initialized")

	basicLogger.Println("Initializing Zap...")
	zapLogger, err := newZapStub(cfg.Env == EnvDev, cfg.res)
	if err != nil {
		return
	}
	zapLogger.Info("Zap initialized")

	zapLogger.Info("Initializing OTEL Trace provider...")
	otelTraceP, err := newOTELTraceProviderStub(cfg.res, false) // TODO: Deal with the bool
	if err != nil {
		return
	}
	setOTELTracerProviderStub(otelTraceP)
	zapLogger.Info("OTEL Trace provider initialized")

	zapLogger.Info("Initializing OTEL Meter provider...")
	otelMeterP, err := newOTELMeterProviderStub(cfg.res)
	if err != nil {
		return
	}
	setOTELMeterProviderStub(otelMeterP)
	zapLogger.Info("OTEL Meter provider initialized")

	ctx = setConfigInContext(ctx, cfg)
	ctx = setZapInContext(ctx, zapLogger)
	shutdown = shutdownFunc(zapLogger, otelTraceP, otelMeterP)

	zapLogger.Info("App initialization complete")
	return
}

func shutdownFunc(
	zapLogger *zap.Logger,
	otelTraceP *sdktrace.TracerProvider,
	otelMeterP *sdkmetric.MeterProvider,
) func() {
	return func() {
		zapLogger.Info("Shutting down app...")

		basicLogger := log.New(os.Stdout, "", log.LstdFlags)

		cancelCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			basicLogger.Println("Shutting down zap...")
			_ = zapLogger.Sync() // Intentionally ignoring err because we zap has a bug where it always returns err here
			basicLogger.Println("Zap shutdown complete")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			basicLogger.Println("Shutting down OTEL Trace provider...")
			if err := otelTraceP.Shutdown(cancelCtx); err != nil {
				basicLogger.Printf("OTEL Trace provider shutdown failed: %s", err.Error())
			} else {
				basicLogger.Println("OTEL Trace provider shutdown complete")
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			basicLogger.Println("Shutting down OTEL Meter provider...")
			if err := otelMeterP.Shutdown(cancelCtx); err != nil {
				basicLogger.Printf("OTEL Meter provider shutdown failed: %s", err.Error())
			} else {
				basicLogger.Println("OTEL Meter provider shutdown complete")
			}
		}()

		// if sentryHub != nil {
		// 	wg.Add(1)
		// 	go func() {
		// 		defer wg.Done()
		// 		_ = sentryHub.Flush(10 * time.Second)
		// 	}()
		// }

		// if nrApp != nil {
		// 	wg.Add(1)
		// 	go func() {
		// 		defer wg.Done()
		// 		nrApp.Shutdown(10 * time.Second)
		// 	}()
		// }

		wg.Wait()

		basicLogger.Println("App shutdown complete")
	}
}
