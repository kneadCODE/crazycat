package gql

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestConvertBadRequestError(t *testing.T) {
	fieldPtr := "str"
	indexPtr := 1
	type testCase struct {
		givenCause   string
		givenMsg     string
		givenPathCtx *graphql.PathContext
		expErr       *gqlerror.Error
	}
	tcs := map[string]testCase{
		"ok without path": {
			givenCause: "invalid_name",
			givenMsg:   "Invalid name given",
			expErr: &gqlerror.Error{
				Message: "Invalid name given",
				Extensions: map[string]interface{}{
					"code":  errCodeBadRequest.String(),
					"cause": "invalid_name",
				},
			},
		},
		"ok with path": {
			givenCause: "invalid_name",
			givenMsg:   "Invalid name given",
			givenPathCtx: &graphql.PathContext{
				Field: &fieldPtr,
				Index: &indexPtr,
			},
			expErr: &gqlerror.Error{
				Message: "Invalid name given",
				Path: []ast.PathElement{
					ast.PathIndex(1),
				},
				Extensions: map[string]interface{}{
					"code":  errCodeBadRequest.String(),
					"cause": "invalid_name",
				},
			},
		},
		"empty cause & msg": {
			expErr: &gqlerror.Error{
				Extensions: map[string]interface{}{
					"code":  errCodeBadRequest.String(),
					"cause": "",
				},
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			ctx := context.Background()
			if tc.givenPathCtx != nil {
				ctx = graphql.WithPathContext(ctx, tc.givenPathCtx)
			}

			// When:
			err := ConvertBadRequestError(ctx, tc.givenCause, tc.givenMsg)

			// Then:
			require.EqualValues(t, tc.expErr, err)
		})
	}
}

func TestConvertUnexpectedError(t *testing.T) {
	fieldPtr := "str"
	indexPtr := 1

	type testCase struct {
		givenErr     error
		givenPathCtx *graphql.PathContext
		expErr       func() *gqlerror.Error
	}
	tcs := map[string]testCase{
		"ok without path": {
			givenErr: errors.New("some err"),
			expErr: func() *gqlerror.Error {
				err := gqlerror.WrapPath(ast.Path(nil), errors.New("some err"))
				err.Message = "An unknown error occurred"
				err.Extensions = map[string]interface{}{
					"code": errCodeInternal.String(),
				}
				return err
			},
		},
		"ok with path": {
			givenErr: errors.New("some err"),
			givenPathCtx: &graphql.PathContext{
				Field: &fieldPtr,
				Index: &indexPtr,
			},
			expErr: func() *gqlerror.Error {
				err := gqlerror.WrapPath(
					[]ast.PathElement{ast.PathIndex(1)},
					errors.New("some err"),
				)
				err.Message = "An unknown error occurred"
				err.Extensions = map[string]interface{}{
					"code": errCodeInternal.String(),
				}
				return err
			},
		},
		"nil err": {
			expErr: func() *gqlerror.Error {
				return nil
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			ctx := context.Background()
			if tc.givenPathCtx != nil {
				ctx = graphql.WithPathContext(ctx, tc.givenPathCtx)
			}

			// When:
			err := ConvertUnexpectError(ctx, tc.givenErr)

			// Then:
			require.EqualValues(t, tc.expErr(), err)
		})
	}
}

func TestErrorPresenter(t *testing.T) {
	type testCase struct {
		givenIntrospectionEnabled bool
		givenErr                  error
		expErr                    func() *gqlerror.Error
	}
	tcs := map[string]testCase{
		"nil err with introspection": {
			expErr: func() *gqlerror.Error {
				return nil
			},
		},
		"nil err without introspection": {
			givenIntrospectionEnabled: true,
			expErr: func() *gqlerror.Error {
				return nil
			},
		},
		"unexp err without introspection": {
			givenErr: ConvertUnexpectError(context.Background(), errors.New("some err")),
			expErr: func() *gqlerror.Error {
				return ConvertUnexpectError(context.Background(), errors.New("some err"))
			},
		},
		"unexp err with introspection": {
			givenIntrospectionEnabled: true,
			givenErr:                  ConvertUnexpectError(context.Background(), errors.New("some err")),
			expErr: func() *gqlerror.Error {
				err := ConvertUnexpectError(context.Background(), errors.New("some err"))
				err.Path = nil
				err.Locations = nil
				return err
			},
		},
		"badreq err without introspection": {
			givenErr: ConvertBadRequestError(context.Background(), "cause", "msg"),
			expErr: func() *gqlerror.Error {
				return ConvertBadRequestError(context.Background(), "cause", "msg")
			},
		},
		"badreq err with introspection": {
			givenIntrospectionEnabled: true,
			givenErr:                  ConvertBadRequestError(context.Background(), "cause", "msg"),
			expErr: func() *gqlerror.Error {
				return ConvertBadRequestError(context.Background(), "cause", "msg")
			},
		},
		"plain err without introspection": {
			givenErr: errors.New("some err"),
			expErr: func() *gqlerror.Error {
				return ConvertUnexpectError(context.Background(), errors.New("some err"))
			},
		},
		"plain err with introspection": {
			givenIntrospectionEnabled: true,
			givenErr:                  errors.New("some err"),
			expErr: func() *gqlerror.Error {
				err := ConvertUnexpectError(context.Background(), errors.New("some err"))
				err.Path = nil
				err.Locations = nil
				return err
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given && When:
			err := errorPresenter(tc.givenIntrospectionEnabled)(context.Background(), tc.givenErr)

			// Then:
			require.EqualValues(t, tc.expErr(), err)
		})
	}
}
