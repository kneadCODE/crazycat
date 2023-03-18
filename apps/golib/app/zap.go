package app

import (
	"fmt"

	"go.uber.org/zap"
)

func newZap(cfg Config) (*zap.Logger, error) {
	var l *zap.Logger
	var err error

	switch cfg.Env {
	case EnvDev:
		l, err = zap.NewDevelopment()
	default:
		l, err = zap.NewProduction()
	}

	// TODO: Set up Zap to use the OTEL log spec

	if err != nil {
		return nil, fmt.Errorf("err initializing zap: %w", err)
	}

	return l, nil
}
