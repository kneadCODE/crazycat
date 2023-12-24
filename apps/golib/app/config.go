package app

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Config holds the application config
type Config struct {
	// Env is the environment in which the application is running
	Env Environment
	res *resource.Resource
}

func newConfigFromEnv(ctx context.Context) (Config, error) {
	res, err := newOTELResourceFromEnvStub(ctx)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{res: res}

	v, _ := res.Set().Value(semconv.DeploymentEnvironmentKey)
	cfg.Env = Environment(v.AsString())
	if err = cfg.Env.IsValid(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Environment denotes the environment where the app is running.
type Environment string

const (
	// EnvProd represents Prod env
	EnvProd = Environment("production")
	// EnvStaging represents Staging environment
	EnvStaging = Environment("staging")
	// EnvDev represents Development environment
	EnvDev = Environment("development")
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// IsValid checks if the environment is valid or not
func (e Environment) IsValid() error {
	if e != EnvDev && e != EnvStaging && e != EnvProd {
		return fmt.Errorf("invalid env: [%s]", e)
	}
	return nil
}
