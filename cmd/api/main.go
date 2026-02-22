package main

import (
	"context"
	"os"
	"os/signal"
	"starter-boilerplate/internal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := internal.InitializeApp(ctx)

	if err := app.Run(ctx); err != nil {
		zap.L().Fatal("server error", zap.Error(err))
	}
}
