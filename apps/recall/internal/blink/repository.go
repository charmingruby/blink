package blink

import (
	"blink/lib/core"
	"blink/lib/telemetry"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type tracerRepository struct {
	db *sqlx.DB
}

func newTracerRepository(db *sqlx.DB) *tracerRepository {
	return &tracerRepository{db: db}
}

func (r *tracerRepository) findByNickname(ctx context.Context, nickname string) (core.Tracer, error) {
	ctx, span := telemetry.StartSpan(ctx, "blink.tracerRepository.findByNickname")
	defer span.End()

	ctx, stop := context.WithTimeout(ctx, 5*time.Second)
	defer stop()

	query := "SELECT * FROM tracers WHERE nickname = $1"

	var tracer core.Tracer
	if err := r.db.GetContext(ctx, &tracer, query, nickname); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Tracer{}, nil
		}

		telemetry.RecordError(ctx, err)

		return core.Tracer{}, err
	}

	return tracer, nil
}
