package httpserver

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	ProfilingEnabled     bool
	ReadinessHandlerFunc http.HandlerFunc
	RESTRoutes           func(chi.Router)
	GQLHandler           http.Handler
}

func (rtr Router) Handler() (chi.Router, error) {
	r := chi.NewRouter()

	r.Get("/_/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = fmt.Fprintln(w, "ok") // Intentionally ignoring the error as nothing to do once caught.
	})

	if rtr.ReadinessHandlerFunc != nil {
		r.Get("/_/ready", rtr.ReadinessHandlerFunc)
	}

	if rtr.ProfilingEnabled {
		profileRoutes(r)
	}

	rootM, err := newRootMiddlewareStub()
	if err != nil {
		return nil, err
	}

	r.Group(func(r chi.Router) {
		r.Use(rootM)

		if rtr.RESTRoutes != nil {
			r.Group(rtr.RESTRoutes)
		}

		if rtr.GQLHandler != nil {
			r.Handle("/graph", rtr.GQLHandler)
		}
	})

	return r, nil
}

func profileRoutes(r chi.Router) {
	// Based on https: //pkg.go.dev/net/http/pprof
	r.HandleFunc("/_/profile/*", pprof.Index)
	r.HandleFunc("/_/profile/cmdline", pprof.Cmdline)
	r.HandleFunc("/_/profile/profile", pprof.Profile)
	r.HandleFunc("/_/profile/symbol", pprof.Symbol)
	r.HandleFunc("/_/profile/trace", pprof.Trace)
	r.Handle("/_/profile/goroutine", pprof.Handler("goroutine"))
	r.Handle("/_/profile/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/_/profile/mutex", pprof.Handler("mutex"))
	r.Handle("/_/profile/heap", pprof.Handler("heap"))
	r.Handle("/_/profile/block", pprof.Handler("block"))
	r.Handle("/_/profile/allocs", pprof.Handler("allocs"))
}
