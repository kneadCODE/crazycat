package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

	ctx, cancel := context.WithCancel(ctx)
	ctx, end := app2.StartSpan(ctx, "span1", false)
	app2.RecordInfoEvent(ctx, "Span1 info 1", attribute.String("sk1", "sv1"))
	time.Sleep(time.Second)
	app2.RecordInfoEvent(ctx, "Span1 info 2", attribute.String("sk2", "sv2"))
	end(nil)

	ctx, end = app2.StartSpan(ctx, "span2", false)
	app2.RecordInfoEvent(ctx, "Span2 info 1", attribute.String("sk1", "sv1"))
	time.Sleep(time.Second)
	app2.RecordInfoEvent(ctx, "Span2 info 2", attribute.String("sk2", "sv2"))
	end(errors.New("some err"))

	newCtx, end := app2.StartSpan(ctx, "span3", true)
	app2.RecordInfoEvent(newCtx, "Span3 info 1", attribute.String("sk1", "sv1"))
	time.Sleep(time.Second)
	app2.RecordInfoEvent(newCtx, "Span3 info 2", attribute.String("sk2", "sv2"))
	end(nil)

	cancel()

	fmt.Println(ctx.Err())
	fmt.Println(newCtx.Err())
}

func f1(ctx context.Context) {
	app2.RecordError(ctx, errors.New("testing err again"), attribute.String("kkkk1", "vvvv1"))
}
