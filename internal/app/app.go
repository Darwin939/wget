package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"wget/internal/server"
)

func Start() {
	app := server.NewApp()
	app.Initialize()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt
		cancel()
	}()
	app.Run(ctx)
}
