package app

import (
	"errors"
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

// ErrInvalidConfig represents an invalid config error
var ErrInvalidConfig = errors.New("invalid config")

func newConfig() (Config, error) {
	cfg := Config{
		Name:             os.Getenv(string(semconv.ServiceNameKey)),
		Project:          os.Getenv(string(semconv.ServiceNamespaceKey)),
		Env:              Environment(os.Getenv(string(semconv.DeploymentEnvironmentKey))),
		Version:          os.Getenv(string(semconv.ServiceVersionKey)),
		ServerInstanceID: os.Getenv(string(semconv.ServiceInstanceIDKey)),
	}

	if cfg.Name == "" {
		return Config{}, fmt.Errorf("name is empty: %w", ErrInvalidConfig)
	}

	if cfg.Project == "" {
		return Config{}, fmt.Errorf("project is empty: %w", ErrInvalidConfig)
	}

	if err := cfg.Env.IsValid(); err != nil {
		return Config{}, err
	}

	if cfg.Version == "" {
		return Config{}, fmt.Errorf("version is empty: %w", ErrInvalidConfig)
	}

	if cfg.ServerInstanceID == "" {
		return Config{}, fmt.Errorf("server instance ID is empty: %w", ErrInvalidConfig)
	}

	return cfg, nil
}
