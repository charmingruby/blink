package main

import (
	"blink/apps/recall/config"
	"blink/lib/env"
	"blink/lib/http/grpcx"
	"blink/lib/telemetry"
	"os"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	log := telemetry.NewLogger()

	log.Info("config: loading")

	cfg, err := env.Load[config.Config]()
	if err != nil {
		log.Error("config: loading error", "error", err)

		return err
	}

	log.Info("config: loaded")

	srv := grpcx.NewServer(cfg.SeverAddress)

	log.Info("server: running", "address", cfg.SeverAddress)

	if err := srv.Start(); err != nil {
		log.Error("server: starting error", "error", err)

		return err
	}

	return nil
}
