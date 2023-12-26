package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kneadCODE/crazycat/apps/golib/app"
	"github.com/kneadCODE/crazycat/apps/golib/app/otelhttpserver"
	"go.opentelemetry.io/otel/attribute"
)

type rootMiddleware struct {
	measure *otelhttpserver.Measure
}

// TODO: Why does it have to be pointer method
func (m *rootMiddleware) serveHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	reqStart := time.Now()
	ctx := r.Context() // Extracting span to pass in to panicHandler

	defer panicHandler(ctx)

	attrs := otelhttpserver.ExtractAttrsFromReq(r) // extract attrs from request

	m.measure.MeasurePreProcessing(ctx, attrs) // OTEL pre-process measuring

	ctx, end := otelhttpserver.StartSpan(r, attrs) // OTEL start span
	defer end(nil)                                 // TODO: See if internal server err should be marking the span as err or not

	bw := &otelhttpserver.RequestBodyWrapper{ReadCloser: r.Body, Ctx: ctx}
	r.Body = bw
	rw := &otelhttpserver.ResponseWriterWrapper{ResponseWriter: w, Ctx: ctx}

	r = r.WithContext(ctx)

	app.RecordInfoEvent(ctx, "START HTTP Request", attribute.String("http.start", reqStart.Format(time.RFC3339)))

	next.ServeHTTP(rw, r) // Actual processing

	reqEnd := time.Now()
	elapsedTime := reqEnd.Sub(reqStart)
	app.RecordInfoEvent(ctx, "END HTTP Request",
		attribute.String("http.duration", fmt.Sprintf("%dms", elapsedTime.Milliseconds())),
		attribute.String("http.end", reqEnd.Format(time.RFC3339)),
	)

	otelhttpserver.SpanPostProcessing(ctx, rw, bw)                   // OTEL post-process span
	m.measure.MeasurePostProcessing(ctx, rw, bw, elapsedTime, attrs) // OTEL post-process measuring
}

func panicHandler(ctx context.Context) {
	rcv := recover()
	if rcv == nil {
		return
	}

	app.RecordError(ctx, fmt.Errorf("httpserver:middleware:RootMiddleware: PANIC: [%+v]", rcv))
}
