package main

import (
	"blink/apps/persister/config"
	"blink/apps/persister/internal/blink"
	"blink/lib/database"
	"blink/lib/env"
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

	db, err := database.NewPostgresClient(ctx, cfg.PostgresURL)
	if err != nil {
		log.Error("postgres: connection error", "error", err)

		return err
	}

	log.Info("postgres: connected to postgres")

	log.Info("rabbitmq: connecting to rabbitmq")

	pubsub, err := queue.NewRabbitMQPubSub(cfg.RabbitMQURL)
	if err != nil {
		log.Error("rabbitmq: connection error", "error", err)

		return err
	}

	log.Info("rabbitmq: connected to rabbitmq")

	shutdownErrCh := make(chan error, 1)
	go gracefulShutdown(ctx, shutdownErrCh, db)

	log.Info("subscriber: listening", "queue", cfg.QueueName)

	if err := blink.Register(log, db, pubsub, cfg.QueueName); err != nil {
		log.Error("subscriber: subscription error", "error", err)

		return err
	}

	err = <-shutdownErrCh
	if err != nil {
		log.Error("shutdown: shutdown error", "error", err)

		return err
	}

	log.Info("shutdown: gracefully shutdown")

	return nil
}

func gracefulShutdown(ctx context.Context, errCh chan error, db *sqlx.DB) {
	<-ctx.Done()

	if err := db.Close(); err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}
