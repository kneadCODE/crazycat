package app

import (
	"os"

	"github.com/getsentry/sentry-go"
)

func initSentry(cfg Config) (*sentry.Hub, error) {
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN == "" {
		return nil, nil
	}

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              sentryDSN,
		AttachStacktrace: true,
		SampleRate:       1, // Send everything
		EnableTracing:    true,
		TracesSampleRate: 1, // Send everything
		SendDefaultPII:   false,
		Integrations:     nil,
		DebugWriter:      nil,
		ServerName:       cfg.ServerInstanceID,
		Release:          cfg.Version,
		Dist:             cfg.Version,
		Environment:      string(cfg.Env),
	})
	if err != nil {
		return nil, err
	}

	return sentry.NewHub(client, nil), nil
}
