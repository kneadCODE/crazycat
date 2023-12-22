package internal

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// NewNewRelicApp returns a new instance of newrelic.Application
func NewNewRelicApp(cfg config.Config) (*newrelic.Application, error) {
	license := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if license == "" {
		return nil, nil
	}

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigLicense(license),
		newrelic.ConfigAppName(cfg.Name),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		func(c *newrelic.Config) {
			c.ErrorCollector.IgnoreStatusCodes = append(
				c.ErrorCollector.IgnoreStatusCodes,
				http.StatusBadRequest,
				http.StatusForbidden,
				http.StatusMethodNotAllowed,
				http.StatusServiceUnavailable,
			)

			c.Labels[string(semconv.ServiceNameKey)] = cfg.Name
			c.Labels[string(semconv.ServiceNamespaceKey)] = cfg.Project
			c.Labels[string(semconv.ServiceVersionKey)] = cfg.Version
			c.Labels[string(semconv.ServiceInstanceIDKey)] = cfg.ServerInstanceID
			c.Labels[string(semconv.DeploymentEnvironmentKey)] = string(cfg.Env)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init new relic failed: %w", err)
	}

	return nrApp, nil
}
