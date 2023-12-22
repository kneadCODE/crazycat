package internal

import (
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/stretchr/testify/require"
)

func Test_NewZap(t *testing.T) {
	// Given: && When:
	l, err := NewZap(config.Config{Name: "name", Env: config.EnvDev})

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)

	// Given: && When:
	l, err = NewZap(config.Config{Name: "name", Env: config.EnvStaging})

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)

	// Given: && When:
	l, err = NewZap(config.Config{Name: "name", Env: config.EnvProd})

	// Then:
	require.NoError(t, err)
	require.NotNil(t, l)
}
