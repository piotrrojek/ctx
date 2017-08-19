package ctx

import (
	"context"
	"os"
	"os/signal"
)

// WithSignal gets parent context and signals on which context should be cancelled.
// It returns context and a cancel function.
func WithSignal(parent context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sig...)

	ctx, cancel := context.WithCancel(parent)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
