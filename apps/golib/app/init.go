// Package app handles the application configuration and is responsible for orchestrating the
// application along with handling all the Metrics, Events, Logs & Traces (MELT/Observability).
package app

import (
	"context"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

// Init initializes the application by setting up the logger, tracer and error reporter.
func Init() (context.Context, func(), error) {
	ctx := context.Background()

	cfg, err := newConfigF()
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

	tp, err := newOTELProviderF(cfg, sentryHub != nil)
	if err != nil {
		return nil, nil, err
	}
	tracer := tp.Tracer("main") // TODO: Check if we need to add some options here.
	ctx = setOTELTracerInContext(ctx, tracer)

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
			_ = tp.Shutdown(cancelCtx) // Intentionally ignoring err because we can't do anything about it
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
	newConfigF       = newConfig
	newZapF          = newZap
	newSentryF       = newSentry
	newNewRelicF     = newNewRelic
	newOTELProviderF = newOTELProvider
)
