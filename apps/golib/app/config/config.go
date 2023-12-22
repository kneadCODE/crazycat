package config

import (
	"errors"
	"fmt"
)

// Config holds the app configuration
type Config struct {
	Name             string
	Env              Environment
	Project          string
	Version          string
	ServerInstanceID string
}

// IsValid validates if the config is valid or not and if not throws an error.
// TODO: Add tests for this
func (cfg Config) IsValid() error {
	if cfg.Name == "" {
		return fmt.Errorf("name is empty: %w", ErrInvalidConfig)
	}

	if cfg.Project == "" {
		return fmt.Errorf("project is empty: %w", ErrInvalidConfig)
	}

	if err := cfg.Env.IsValid(); err != nil {
		return err
	}

	if cfg.Version == "" {
		return fmt.Errorf("version is empty: %w", ErrInvalidConfig)
	}

	if cfg.ServerInstanceID == "" {
		return fmt.Errorf("server instance ID is empty: %w", ErrInvalidConfig)
	}

	return nil
}

// ErrInvalidConfig represents an invalid config error
var ErrInvalidConfig = errors.New("invalid config")
