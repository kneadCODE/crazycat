package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/kneadCODE/crazycat/apps/golib/app2"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	log.New(os.Stdout, "", log.LstdFlags).Println("Hello from cazycat golib")

	ctx, finish, err := app2.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer finish()

	var _ = ctx

	app2.RecordDebugEvent(ctx, "Testing debug", attribute.String("kkkk1", "vvvv1"))
	app2.RecordInfoEvent(ctx, "Testing info", attribute.String("kkkk1", "vvvv1"))
	app2.RecordWarnEvent(ctx, "Testing warn", attribute.String("kkkk1", "vvvv1"))
	app2.RecordError(ctx, errors.New("testing err"), attribute.String("kkkk1", "vvvv1"))
	f1(ctx)
}

func f1(ctx context.Context) {
	app2.RecordError(ctx, errors.New("testing err again"), attribute.String("kkkk1", "vvvv1"))
}
