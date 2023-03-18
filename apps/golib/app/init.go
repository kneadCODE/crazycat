package app

import (
	"context"
)

// Init initalizes the application by setting up the logger, tracer and error reporter.
func Init() (context.Context, func(), error) {
	ctx := context.Background()

	cfg, err := newConfigFromEnv()
	if err != nil {
		return nil, nil, err
	}

	logger, err := newZap(cfg)
	if err != nil {
		return nil, nil, err
	}

	ctx = setConfigInContext(ctx, cfg)
	ctx = setZapInContext(ctx, logger)

	return ctx, func() {
		// TODO: Add cleanup/flush/sync logic here
	}, nil
}
