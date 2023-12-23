package app2

import (
	"fmt"
)

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
