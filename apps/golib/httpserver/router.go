package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kneadCODE/crazycat/apps/golib/app"
	"github.com/kneadCODE/crazycat/apps/golib/app/otelhttpserver"
	"go.opentelemetry.io/otel/attribute"
)

type Router struct {
	ProfilingEnabled     bool
	ReadinessHandlerFunc http.HandlerFunc
	RESTRoutes           func(chi.Router)
	GQLHandler           http.Handler
}

func (rt Router) Handler() chi.Router {
	rtr := chi.NewRouter()

	rtr.Get("/_/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = fmt.Fprintln(w, "ok") // Intentionally ignoring the error as nothing to do once caught.
	})

	// TODO: Deal with this
	// if rt.ReadinessHandlerFunc != nil {
	// 	rtr.Get("/_/ready", rt.ReadinessHandlerFunc)
	// }

	if rt.ProfilingEnabled {
		// Based on https: //pkg.go.dev/net/http/pprof
		rtr.HandleFunc("/_/profile/*", pprof.Index)
		rtr.HandleFunc("/_/profile/cmdline", pprof.Cmdline)
		rtr.HandleFunc("/_/profile/profile", pprof.Profile)
		rtr.HandleFunc("/_/profile/symbol", pprof.Symbol)
		rtr.HandleFunc("/_/profile/trace", pprof.Trace)
		rtr.Handle("/_/profile/goroutine", pprof.Handler("goroutine"))
		rtr.Handle("/_/profile/threadcreate", pprof.Handler("threadcreate"))
		rtr.Handle("/_/profile/mutex", pprof.Handler("mutex"))
		rtr.Handle("/_/profile/heap", pprof.Handler("heap"))
		rtr.Handle("/_/profile/block", pprof.Handler("block"))
		rtr.Handle("/_/profile/allocs", pprof.Handler("allocs"))
	}

	rtr.Group(func(r chi.Router) {
		r.Use(rootMiddleware)

		if rt.RESTRoutes != nil {
			r.Group(rt.RESTRoutes)
		}

		if rt.GQLHandler != nil {
			r.Handle("/graph", rt.GQLHandler)
		}
	})

	return rtr
}

func rootMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqStart := time.Now()
		ctx := r.Context()

		defer panicHandler(ctx)

		ctx, end := otelhttpserver.StartSpan(r)
		defer end(nil) // TODO: See if internal server err should be marking the span as err or not

		bw := &otelhttpserver.OTELRequestBodyWrapper{ReadCloser: r.Body, Ctx: ctx}
		r.Body = bw
		rw := &otelhttpserver.OTELResponseWrapper{ResponseWriter: w, Ctx: ctx}

		r = r.WithContext(ctx)

		app.RecordInfoEvent(ctx, "START HTTP Request",
			attribute.String("http.req.start", reqStart.Format(time.RFC3339)),
		)

		next.ServeHTTP(rw, r)

		otelhttpserver.PostHandle(ctx, rw, bw)

		reqEnd := time.Now()
		app.RecordInfoEvent(ctx, "END HTTP Request",
			attribute.String("http.resp.total-duration", fmt.Sprintf("%dms", time.Since(reqEnd).Milliseconds())),
			attribute.String("http.resp.end", reqEnd.Format(time.RFC3339)),
		)
	})
}

func panicHandler(ctx context.Context) {
	rcv := recover()
	if rcv == nil {
		return
	}

	app.RecordError(ctx, fmt.Errorf(
		"httpserver:middleware:RootMiddleware: PANIC: [%+v]", rcv))
}
