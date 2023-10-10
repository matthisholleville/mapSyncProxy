package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/matthisholleville/mapsyncproxy/api"
)

var (
	sigs = make(chan os.Signal, 1)
)

func main() {
	ctx := context.Background()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	api.Server(ctx, sigs)
}
