package app

// Environment denotes the environment where the app is running.
type Environment string

const (
	// EnvProd represents Prod env
	EnvProd = Environment("production")
	// EnvStaging represents Staging environment
	EnvStaging = Environment("staging")
	// EnvDev represents Development environment
	EnvDev = Environment("dev")
)
