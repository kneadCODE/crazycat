package gql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	require.NotNil(t, Handler(nil, false))
	require.NotNil(t, Handler(nil, true))
	// TODO: Figure out how to write proper unit tests for this
}
