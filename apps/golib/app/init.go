// Package app handles the application configuration and is responsible for orchestrating the
// application along with handling all the Metrics, Events, Logs & Traces (MELT/Observability).
package app

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/kneadCODE/crazycat/apps/golib/app/internal"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Init initializes the application by setting up the logger, tracer and error reporter.
func Init() (context.Context, func(), error) {
	ctx := context.Background()

	cfg, err := newConfigFromEnvF()
	if err != nil {
		return nil, nil, err
	}
	ctx = setConfigInContext(ctx, cfg)

	// TODO: Add log line to inform that processing has started

	logger, err := newZapF(cfg)
	if err != nil {
		return nil, nil, err
	}
	ctx = setZapInContext(ctx, logger)

	sentryHub, err := newSentryF(cfg)
	if err != nil {
		return nil, nil, err
	}
	if sentryHub != nil {
		ctx = sentry.SetHubOnContext(ctx, sentryHub)
	}

	nrApp, err := newNewRelicF(cfg)
	if err != nil {
		return nil, nil, err
	}
	if nrApp != nil {
		ctx = setNewRelicInContext(ctx, nrApp)
	}

	otelShutdown, err := initOTELF(cfg, sentryHub != nil)
	if err != nil {
		return nil, nil, err
	}

	return ctx, func() {
		cancelCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = logger.Sync() // Intentionally ignoring err because we can't do anything about it
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = otelShutdown(cancelCtx) // Intentionally ignoring err because we can't do anything about it
		}()

		if sentryHub != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = sentryHub.Flush(10 * time.Second)
			}()
		}

		if nrApp != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				nrApp.Shutdown(10 * time.Second)
			}()
		}

		wg.Wait()
	}, nil
}

// The following are for stubbing in tests
var (
	newZapF      = internal.NewZap
	newSentryF   = internal.NewSentryHub
	newNewRelicF = internal.NewNewRelicApp
	initOTELF    = internal.InitOTEL
)

// newConfigFromEnv initalizes a new config from environment
func newConfigFromEnv() (config.Config, error) {
	cfg := config.Config{
		Name:             os.Getenv(string(semconv.ServiceNameKey)),
		Project:          os.Getenv(string(semconv.ServiceNamespaceKey)),
		Env:              config.Environment(os.Getenv(string(semconv.DeploymentEnvironmentKey))),
		Version:          os.Getenv(string(semconv.ServiceVersionKey)),
		ServerInstanceID: os.Getenv(string(semconv.ServiceInstanceIDKey)),
	}

	if err := cfg.IsValid(); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

var newConfigFromEnvF = newConfigFromEnv
