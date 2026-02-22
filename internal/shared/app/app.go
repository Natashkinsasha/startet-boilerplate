package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"starter-boilerplate/internal/shared/config"

	gohuma "github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	gogrpc "google.golang.org/grpc"
)

type App struct {
	HTTPServer *http.Server
	GRPCServer *gogrpc.Server
	Config     *config.Config
	Api        gohuma.API
	ready      chan struct{}
	startErr   chan error
}

func New(httpSrv *http.Server, cfg *config.Config, grpcSrv *gogrpc.Server, api gohuma.API) *App {
	return &App{
		HTTPServer: httpSrv,
		GRPCServer: grpcSrv,
		Config:     cfg,
		Api:        api,
		ready:      make(chan struct{}),
		startErr:   make(chan error, 1),
	}
}

// Ready returns a channel that is closed once both HTTP and gRPC servers are listening.
func (a *App) Ready() <-chan struct{} {
	return a.ready
}

// StartErr returns a channel that receives an error if the server fails to start.
func (a *App) StartErr() <-chan error {
	return a.startErr
}

// BaseURL returns the HTTP base URL of the application.
func (a *App) BaseURL() string {
	return fmt.Sprintf("http://localhost:%d", a.Config.App.Port)
}

func (a *App) Run(ctx context.Context) error {
	httpLn, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.App.Port))
	if err != nil {
		err = fmt.Errorf("http listen: %w", err)
		a.startErr <- err
		return err
	}

	grpcLn, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.GRPC.Port))
	if err != nil {
		httpLn.Close()
		err = fmt.Errorf("grpc listen: %w", err)
		a.startErr <- err
		return err
	}

	close(a.ready)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		zap.L().Info("http server started", zap.Int("port", a.Config.App.Port))
		if err := a.HTTPServer.Serve(httpLn); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		zap.L().Info("grpc server started", zap.Int("port", a.Config.GRPC.Port))
		if err := a.GRPCServer.Serve(grpcLn); err != nil {
			return fmt.Errorf("grpc server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		return a.shutdown()
	})

	return g.Wait()
}

func (a *App) shutdown() error {
	zap.L().Info("shutting down servers...")

	timeout := a.Config.App.ShutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var g errgroup.Group

	g.Go(func() error {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		zap.L().Info("http server stopped")
		return nil
	})

	g.Go(func() error {
		stopped := make(chan struct{})
		go func() {
			a.GRPCServer.GracefulStop()
			close(stopped)
		}()

		t := time.NewTimer(timeout)
		defer t.Stop()

		select {
		case <-stopped:
			zap.L().Info("grpc server stopped gracefully")
		case <-t.C:
			zap.L().Warn("grpc server shutdown timed out, forcing stop")
			a.GRPCServer.Stop()
		}
		return nil
	})

	return g.Wait()
}
