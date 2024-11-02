package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SetupContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()
	return ctx

}
