package app

import (
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

func newConfig() (Config, error) {
	cfg := Config{
		Name:             os.Getenv(string(semconv.ServiceNameKey)),
		Project:          os.Getenv(string(semconv.ServiceNamespaceKey)),
		Env:              Environment(os.Getenv(string(semconv.DeploymentEnvironmentKey))),
		Version:          os.Getenv(string(semconv.ServiceVersionKey)),
		ServerInstanceID: os.Getenv(string(semconv.ServiceInstanceIDKey)),
	}

	// TODO: Add validation

	return cfg, nil
}
