package app

import (
	"fmt"
	"os"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Config holds the app configuration
type Config struct {
	Name             string
	Env              Environment
	Project          string
	Version          string
	ServerInstanceID string
}

func newConfigFromEnv() (Config, error) {
	cfg := Config{
		Name:    os.Getenv(string(semconv.ServiceNameKey)),
		Project: os.Getenv(string(semconv.ServiceNamespaceKey)),
		Env:     Environment(os.Getenv(string(semconv.DeploymentEnvironmentKey))),
		ServerInstanceID: fmt.Sprintf("%s:%s",
			os.Getenv(os.Getenv(string(semconv.K8SPodUIDKey))),
			os.Getenv(os.Getenv(string(semconv.K8SNodeUIDKey))),
		),
	}

	// TODO: Add validation

	return cfg, nil
}
