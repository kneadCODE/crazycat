package internal

import (
	"testing"

	"github.com/kneadCODE/crazycat/apps/golib/app/config"
	"github.com/stretchr/testify/require"
)

func Test_newOTELProvider(t *testing.T) {
	// Given && When:
	tp, err := InitOTEL(config.Config{}, false)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, tp)

	// Given && When:
	tp, err = InitOTEL(config.Config{}, true)

	// Then:
	require.NoError(t, err)
	require.NotNil(t, tp)

	// TODO: Figure out how to write proper tests for OTEL configs. Only choice I see now is using interfaces :(
}
