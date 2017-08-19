package ctx

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWithSignal(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Error("cannot get the root process")
	}

	ctx := context.Background()
	ctx, _ = WithSignal(ctx, syscall.SIGXCPU)
	ch := make(chan struct{}, 1)
	go functionToBeCancelled(ctx, ch)
	time.Sleep(50 * time.Millisecond)

	err = proc.Signal(syscall.SIGXCPU)
	if err != nil {
		t.Error("cannot send interrupt signal")
	}

	select {
	case <-ch:
		fmt.Println("the function exited")
	case <-time.After(50 * time.Millisecond):
		t.Error("function didn't exit")
	}
}

func functionToBeCancelled(ctx context.Context, ch chan struct{}) {
	select {
	case <-ctx.Done():
		close(ch)
		return
	default:
	}
}
