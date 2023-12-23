package app2

import (
	"context"

	"github.com/kneadCODE/crazycat/apps/golib/app2/internal"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Config struct {
	Env Environment
	res *resource.Resource
}

func newConfigFromEnv(ctx context.Context) (Config, error) {
	res, err := internal.NewOTELResourceFromEnv(ctx)
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
