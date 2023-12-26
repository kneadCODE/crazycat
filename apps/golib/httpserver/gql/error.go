package gql

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/kneadCODE/crazycat/apps/golib/app"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ErrorCode represents the code returned to the client
type ErrorCode string

// String returns the string representation of the ErrorCode
func (e ErrorCode) String() string {
	return string(e)
}

// Based on Apollo's spec - https://www.apollographql.com/docs/apollo-server/data/errors
var (
	// errCodeBadRequest means a bad request was sent
	errCodeBadRequest = ErrorCode("BAD_REQUEST")
	// errCodeInternal means an internal error occurred
	errCodeInternal = ErrorCode("INTERNAL_SERVER_ERROR")
	// errCodeUnauthenticated means the request was not authenticated
	errCodeUnauthenticated = ErrorCode("UNAUTHENTICATED")
	// errCodeForbidden means the request was not authorized
	errCodeForbidden = ErrorCode("FORBIDDEN")
)

// ConvertBadRequestError converts the known error into *gqlerror.Error
func ConvertBadRequestError(ctx context.Context, cause string, message string) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: message,
		Extensions: map[string]interface{}{
			"code":  errCodeBadRequest.String(),
			"cause": cause,
		},
	}
}

// ConvertUnexpectError converts the given unexpected error into *gqlerror.Error
func ConvertUnexpectError(ctx context.Context, err error) *gqlerror.Error {
	if err == nil {
		return nil
	}
	gerr := gqlerror.WrapPath(graphql.GetPath(ctx), err)
	gerr.Message = "An unknown error occurred"
	gerr.Extensions = map[string]interface{}{
		"code": errCodeInternal.String(),
	}
	return gerr
}

func errorPresenter(isIntrospectionEnabled bool) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		if err == nil {
			return nil
		}

		var gerr *gqlerror.Error
		if !errors.As(err, &gerr) {
			gerr = ConvertUnexpectError(ctx, err)
		}

		// Don't expose any schema-identifiable info when introspection is disabled
		if !isIntrospectionEnabled {
			gerr.Locations = nil
			gerr.Path = nil
		}

		if underlyingErr := gerr.Unwrap(); underlyingErr != nil {
			app.RecordError(ctx, underlyingErr)
		}

		return gerr
	}
}
