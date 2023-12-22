package internal

import (
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/kneadCODE/crazycat/apps/golib/app/config"
)

// NewSentryHub returns a new instance of sentry.Hub
func NewSentryHub(cfg config.Config) (*sentry.Hub, error) {
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
		Environment:      cfg.Env.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("init sentry failed: %w", err)
	}

	return sentry.NewHub(client, nil), nil
}
