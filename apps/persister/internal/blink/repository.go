package blink

import (
	"blink/lib/core"
	"blink/lib/database"
	"context"

	"github.com/jmoiron/sqlx"
)

type tracerRepository struct {
	db database.Querier
}

func newTracerRepository(db *sqlx.DB) *tracerRepository {
	return &tracerRepository{db: db}
}

func (r *tracerRepository) withTx(tx *sqlx.Tx) *tracerRepository {
	return &tracerRepository{db: tx}
}

func (r *tracerRepository) create(ctx context.Context, tr core.Tracer) error {
	query := "INSERT INTO tracers (id, nickname, total_blinks, created_at, last_blink_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.ExecContext(ctx, query, tr.ID, tr.Nickname, tr.TotalBlinks, tr.CreatedAt, tr.LastBlinkAt, tr.UpdatedAt)

	return err
}

func (r *tracerRepository) update(ctx context.Context, tr core.Tracer) error {
	query := "UPDATE tracers SET total_blinks = $1, last_blink_at = $2, updated_at = $3 WHERE nickname = $4"

	_, err := r.db.ExecContext(ctx, query, tr.TotalBlinks, tr.LastBlinkAt, tr.UpdatedAt, tr.Nickname)

	return err
}

type blinkRepository struct {
	db database.Querier
}

func newBlinkRepository(db *sqlx.DB) *blinkRepository {
	return &blinkRepository{db: db}
}

func (r *blinkRepository) withTx(tx *sqlx.Tx) *blinkRepository {
	return &blinkRepository{db: tx}
}

func (r *blinkRepository) create(ctx context.Context, bl core.Blink) error {
	query := "INSERT INTO blinks (id, tracer_id, created_at) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, query, bl.ID, bl.TracerID, bl.CreatedAt)

	return err
}
