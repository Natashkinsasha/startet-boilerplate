package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"starter-boilerplate/internal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := internal.InitializeApp(ctx)

	if err := app.Run(ctx); err != nil {
		slog.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}
