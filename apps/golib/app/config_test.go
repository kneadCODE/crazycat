package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvironment_String(t *testing.T) {
	require.Equal(t, "development", EnvDev.String())
	require.Equal(t, "staging", EnvStaging.String())
	require.Equal(t, "production", EnvProd.String())
	require.Equal(t, "abc", Environment("abc").String())
}

func TestEnvironment_IsValid(t *testing.T) {
	require.NoError(t, EnvDev.IsValid())
	require.NoError(t, EnvStaging.IsValid())
	require.NoError(t, EnvProd.IsValid())
	require.Equal(t, errors.New("invalid env: [abc]"), Environment("abc").IsValid())
	require.Equal(t, errors.New("invalid env: []"), Environment("").IsValid())
}
