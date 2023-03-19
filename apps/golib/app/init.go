package app

import (
	"context"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

// Init initalizes the application by setting up the logger, tracer and error reporter.
func Init() (context.Context, func(), error) {
	ctx := context.Background()

	cfg, err := newConfigFromEnv()
	if err != nil {
		return nil, nil, err
	}
	ctx = setConfigInContext(ctx, cfg)

	// TODO: Add log line to inform that processing has started

	logger, err := newZap(cfg)
	if err != nil {
		return nil, nil, err
	}
	ctx = setZapInContext(ctx, logger)

	tp, err := newTracer(cfg)
	if err != nil {
		return nil, nil, err
	}
	tracer := tp.Tracer("main") // TODO: Check if we need to add some options here.
	ctx = setOTELTracerInContext(ctx, tracer)

	sentryHub, err := initSentry(cfg)
	if err != nil {
		return nil, nil, err
	}
	if sentryHub != nil {
		ctx = sentry.SetHubOnContext(ctx, sentryHub)
	}

	return ctx, func() {
		cancelCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			_ = logger.Sync() // Intentionally ignoring err because we can't do anything about it
		}()

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

		wg.Wait()
	}, nil
}
