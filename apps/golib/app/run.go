package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Run runs the various services and listens to exit signals to terminate all the services.
func Run(ctx context.Context, services ...service) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(services))

	RecordInfoEvent(ctx, "Starting all services")

	for i := range services {
		svc := services[i]
		go func() {
			defer wg.Done()
			if err := svc(ctx); err != nil {
				RecordError(ctx, fmt.Errorf("svc err: %w", err))

				cancel() // cancel ctx for other services to terminate.
			}
		}()
	}

	select {
	case sig := <-exitSignalStub():
		RecordInfoEvent(ctx, fmt.Sprintf(
			"Exit signal: [%s] received. Terminating all services",
			sig.String()),
		)

		cancel()
	case <-ctx.Done():
		RecordInfoEvent(ctx, "Context cancelled. Terminating all services")
	}

	wg.Wait()

	RecordInfoEvent(ctx, "All services shut down. Exiting app.")
}

// service represents an executable that is context aware and will return an error if encountered.
type service func(ctx context.Context) error

// will be used for stubbing in tests
func exitSignal() <-chan os.Signal {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGTERM, os.Interrupt, os.Kill)
	return exitChan
}
