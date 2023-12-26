package gql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/kneadCODE/crazycat/apps/golib/app"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var _ graphql.RecoverFunc = recoverFunc

func recoverFunc(ctx context.Context, rcv interface{}) error {
	if rcv == nil {
		return nil
	}

	app.RecordError(ctx, fmt.Errorf("httpserver:gql:Recover: PANIC: [%+v]", rcv))

	return &gqlerror.Error{
		Message: "An unknown error occurred",
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": errCodeInternal.String(),
		},
	}
}
