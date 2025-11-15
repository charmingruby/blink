package main

import (
	"blink/apps/recall/config"
	"blink/apps/recall/internal/evaluate"
	"blink/lib/database"
	"blink/lib/env"
	"blink/lib/http/grpcx"
	"blink/lib/queue"
	"blink/lib/telemetry"
	"context"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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

	log.Info("postgres: connecting to postgres")

	db, err := database.NewPostgresClient(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("postgres: connection error", "error", err)

		return err
	}

	log.Info("postgres: connected to postgres")

	log.Info("rabbitmq: connecting to rabbitmq")

	pubsub, err := queue.NewRabbitMQPubSub(cfg.QueueURL)
	if err != nil {
		log.Error("rabbitmq: connection error", "error", err)

		return err
	}

	log.Info("rabbitmq: connected to rabbitmq")

	srv := grpcx.NewServer(cfg.ServerAddress)

	evaluate.Scaffold(srv.Conn, db, pubsub)

	log.Info("server: running", "address", cfg.ServerAddress)

	shutdownErrCh := make(chan error, 1)
	go gracefulShutdown(ctx, shutdownErrCh, db, srv)

	if err := srv.Start(); err != nil {
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

func gracefulShutdown(ctx context.Context, errCh chan error, db *sqlx.DB, srv *grpcx.Server) {
	<-ctx.Done()

	if err := db.Close(); err != nil {
		errCh <- err
		return
	}

	srv.Stop()

	errCh <- nil
}
