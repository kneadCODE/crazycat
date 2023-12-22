package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/kneadCODE/crazycat/apps/golib/app"
)

// New returns a new instance of Server.
func New(ctx context.Context, router Router, options ...ServerOption) (*Server, error) {
	s := &Server{
		srv: &http.Server{
			Addr:         ":9000",
			Handler:      router.Handler(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			BaseContext: func(net.Listener) context.Context {
				// TODO: Figure out why app pkg does not have a way to create a new instance
				// return app.NewContext(context.Background(), app.FromContext(ctx))
				return context.Background()
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
		app.TrackInfoEvent(ctx, fmt.Sprintf("Starting HTTP server on %s", s.srv.Addr))
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
	ctx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdownTimeout) // Cannot rely on root context as that might have been cancelled.
	defer cancel()

	app.TrackInfoEvent(ctx, "Attempting HTTP server graceful shutdown")
	if err := s.srv.Shutdown(ctx); err != nil {
		app.TrackErrorEvent(ctx, fmt.Errorf("httpserver:Server: graceful shutdown failed: %w", err))

		app.TrackInfoEvent(ctx, "Attempting HTTP server force shutdown")
		if err = s.srv.Close(); err != nil {
			app.TrackErrorEvent(ctx, fmt.Errorf("httpserver:Server: force shutdown failed: %w", err))
			return err
		}
	}

	app.TrackInfoEvent(ctx, "HTTP server shutdown complete")

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
