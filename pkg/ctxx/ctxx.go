package ctxx

import (
	"context"
	"os"
	"os/signal"
)

// InterruptContext returns a context that is canceled when the process is interrupted.
func InterruptContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
}
