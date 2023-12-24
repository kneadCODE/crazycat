package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kneadCODE/crazycat/apps/golib/app2"
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

		// TODO: Figure out why app pkg does not have a way to put fields in context
		// logger = logger.With(
		// 	slog.String("http.req.method", r.Method),
		// 	slog.String("http.req.path", r.URL.Path),
		// 	slog.String("http.req.host", r.URL.Host),
		// 	slog.String("http.req.user-agent", r.UserAgent()),
		// 	slog.String("http.req.referer", r.Referer()),
		// 	slog.String("http.req.remote_addr", r.RemoteAddr),
		// )

		app2.RecordInfoEvent(ctx, "START HTTP Request",
			attribute.String("http.req.content-type", r.Header.Get("Content-Type")),
			attribute.String("http.req.proto", r.Proto),
			attribute.String("http.req.start", reqStart.Format(time.RFC3339)),
		)

		writer := &respWriter{ResponseWriter: w}

		next.ServeHTTP(writer, r)

		reqEnd := time.Now()
		app2.RecordInfoEvent(ctx, "END HTTP Request",
			attribute.String("http.resp.status", strconv.Itoa(writer.statusCode)),
			attribute.String("http.resp.total-duration", fmt.Sprintf("%dms", time.Since(reqEnd).Milliseconds())),
			attribute.String("http.resp.content-length", writer.Header().Get("Content-Length")),
			attribute.String("http.resp.end", reqEnd.Format(time.RFC3339)),
		)
	})
}

func panicHandler(ctx context.Context) {
	rcv := recover()
	if rcv == nil {
		return
	}

	// TODO: Add additional log fields if necessary.
	app2.RecordError(ctx, fmt.Errorf(
		"httpserver:middleware:RootMiddleware: PANIC: [%+v]", rcv))
}

type respWriter struct {
	http.ResponseWriter

	statusCode int
}

func (w *respWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
