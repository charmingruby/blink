package main

import (
	"blink/apps/edgeapi/config"
	"blink/lib/database"
	"blink/lib/env"
	"blink/lib/http/rest"
	"blink/lib/telemetry"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
)

const TIMEOUT = 30 * time.Second

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	log := telemetry.NewLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	log.Info("config: loading")

	cfg, err := env.Load[config.Config]()
	if err != nil {
		log.Error("config: loading error", "error", err)

		return err
	}

	log.Info("config: loaded")

	log.Info("tracer: loading")

	if err := telemetry.NewTracer(telemetry.TracerConfig{
		ServiceName: cfg.ServiceName,
		Endpoint:    cfg.OTLPExporterEndpoint,
	}); err != nil {
		log.Info("tracer: setup error", "error", err)

		return err
	}

	log.Info("tracer: ready")

	log.Info("postgres: connecting")

	db, err := database.NewPostgresClient(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("postgres: connection error")

		return err
	}

	log.Info("postgres: connected")

	srv := rest.NewServer(cfg.Port)

	log.Info("server: running", "port", cfg.Port)

	shutdownErrCh := make(chan error, 1)
	go gracefulShutdown(ctx, shutdownErrCh, &srv.Server, db)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server: starting error", "error", err)

		return err
	}

	log.Info("shutdown: shutting down application")

	err = <-shutdownErrCh
	if err != nil {
		log.Error("shutdown: shutdown error", "error", err)

		return err
	}

	log.Info("shutdown: gracefully shutdown")

	return nil
}

func gracefulShutdown(ctx context.Context, errCh chan error, srv *http.Server, db *sqlx.DB) {
	<-ctx.Done()

	shutdownCtx, stop := context.WithTimeout(context.Background(), TIMEOUT)
	defer stop()

	if err := db.Close(); err != nil {
		errCh <- err
		return
	}

	err := srv.Shutdown(shutdownCtx)
	switch err {
	case nil:
		errCh <- nil
		return
	case context.DeadlineExceeded:
		errCh <- fmt.Errorf("shutdown: forcing closing the server, %w", err)
	default:
		errCh <- fmt.Errorf("shutdown: forcing closing the server, %w", err)
	}
}
