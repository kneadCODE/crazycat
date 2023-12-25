package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/kneadCODE/crazycat/apps/golib/app"
	"github.com/kneadCODE/crazycat/apps/golib/httpserver/router"
)

// New returns a new instance of Server.
func New(ctx context.Context, rtr router.Router, options ...ServerOption) (*Server, error) {
	handler, err := rtr.Handler()
	if err != nil {
		return nil, err
	}

	s := &Server{
		srv: &http.Server{
			Addr:         ":9000",
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			BaseContext: func(net.Listener) context.Context {
				return app.CloneNewContext(ctx)
			},
		},
		gracefulShutdownTimeout: 10 * time.Second,
	}

	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Server is the server instance
type Server struct {
	srv                     *http.Server
	gracefulShutdownTimeout time.Duration
}

// Start starts the server and is context aware and shuts down when the context gets cancelled.
func (s *Server) Start(ctx context.Context) error {
	startErrChan := make(chan error, 1)

	go func() {
		app.RecordInfoEvent(ctx, fmt.Sprintf("Starting HTTP server on %s", s.srv.Addr))
		startErrChan <- s.srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			return s.stop(ctx)
		case err := <-startErrChan:
			if err != http.ErrServerClosed {
				return fmt.Errorf("http server startup failed: %w", err)
			}
			return nil
		}
	}
}

func (s *Server) stop(ctx context.Context) error {
	cancelCtx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdownTimeout) // Cannot rely on root context as that might have been cancelled.
	defer cancel()

	app.RecordInfoEvent(ctx, "Attempting HTTP server graceful shutdown")
	if err := s.srv.Shutdown(cancelCtx); err != nil {
		app.RecordError(ctx, fmt.Errorf("httpserver:Server: graceful shutdown failed: %w", err))

		app.RecordInfoEvent(ctx, "Attempting HTTP server force shutdown")
		if err = s.srv.Close(); err != nil {
			app.RecordError(ctx, fmt.Errorf("httpserver:Server: force shutdown failed: %w", err))
			return err
		}
	}

	app.RecordInfoEvent(ctx, "HTTP server shutdown complete")

	return nil
}

// ServerOption customizes the Server
type ServerOption = func(*Server) error

// WithServerPort sets the server port to the given port
func WithServerPort(port int) ServerOption {
	return func(s *Server) error {
		s.srv.Addr = fmt.Sprintf(":%d", port)
		return nil
	}
}
