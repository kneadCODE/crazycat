package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kneadCODE/crazycat/apps/golib/app"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	log.New(os.Stdout, "", log.LstdFlags).Println("Hello from cazycat golib")

	ctx, finish, err := app.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer finish()

	var _ = ctx

	app.RecordDebugEvent(ctx, "Testing debug", attribute.String("kkkk1", "vvvv1"))
	app.RecordInfoEvent(ctx, "Testing info", attribute.String("kkkk1", "vvvv1"))
	app.RecordWarnEvent(ctx, "Testing warn", attribute.String("kkkk1", "vvvv1"))
	app.RecordError(ctx, errors.New("testing err"), attribute.String("kkkk1", "vvvv1"))
	f1(ctx)

	ctx, cancel := context.WithCancel(ctx)
	ctx, end := app.StartSpan(ctx, "span1", false, attribute.String("sp1", "v1"))
	app.RecordInfoEvent(ctx, "Span1 info 1", attribute.String("sk1", "sv1"))
	time.Sleep(time.Second)
	ctx = app.ContextWithAttributes(ctx, attribute.String("NEWK1", "NEWV1"), attribute.String("sk1", "overriden_sv1"))
	app.RecordInfoEvent(ctx, "Span1 info 2", attribute.String("sk2", "sv2"))
	end(nil)

	ctx, end = app.StartSpan(ctx, "span2", false)
	app.RecordInfoEvent(ctx, "Span2 info 1", attribute.String("sk1", "sv1"))
	time.Sleep(time.Second)
	app.RecordInfoEvent(ctx, "Span2 info 2", attribute.String("sk2", "sv2"))
	end(errors.New("some err"))

	newCtx, end := app.StartSpan(ctx, "span3", true)
	app.RecordInfoEvent(newCtx, "Span3 info 1", attribute.String("sk1", "sv1"))
	time.Sleep(time.Second)
	app.RecordInfoEvent(newCtx, "Span3 info 2", attribute.String("sk2", "sv2"))
	end(nil)

	cancel()

	fmt.Println(ctx.Err())
	fmt.Println(newCtx.Err())
}

func f1(ctx context.Context) {
	app.RecordError(ctx, errors.New("testing err again"), attribute.String("kkkk1", "vvvv1"))
}
