package telemetry

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Logger = slog.Logger

func NewLogger() *Logger {
	return otelslog.NewLogger("")
}
