package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"starter-boilerplate/internal"
	"starter-boilerplate/internal/shared/huma"
	"syscall"
)

func main() {
	appCtx, appStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer appStop()
	app := internal.InitializeApp(appCtx)
	slog.Info("app initialized")
	huma.GenerateSpecFile(app.Api)
	slog.Info("swagger spec generated")
}
