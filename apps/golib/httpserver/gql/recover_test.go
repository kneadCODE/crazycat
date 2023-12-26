package gql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Test_recover(t *testing.T) {
	// Given: no err
	ctx := context.Background()

	// When:
	err := recoverFunc(ctx, nil)

	// Then:
	require.NoError(t, err)

	// When: non-err err
	err = recoverFunc(ctx, "abc")

	// Then:
	require.Equal(t, &gqlerror.Error{
		Message: "An unknown error occurred",
		Path:    nil,
		Extensions: map[string]interface{}{
			"code": errCodeInternal.String(),
		},
	}, err)

	// When: err err
	err = recoverFunc(ctx, errors.New("some err"))

	// Then:
	require.Equal(t, &gqlerror.Error{
		Message: "An unknown error occurred",
		Path:    nil,
		Extensions: map[string]interface{}{
			"code": errCodeInternal.String(),
		},
	}, err)
}
