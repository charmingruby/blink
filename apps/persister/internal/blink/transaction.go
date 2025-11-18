package blink

import (
	"blink/lib/database"
	"blink/lib/telemetry"
	"context"

	"github.com/jmoiron/sqlx"
)

type tracerBlinkTransactionManager struct {
	txManager  *database.PostgresTxManager
	tracerRepo *tracerRepository
	blinkRepo  *blinkRepository
}

func newTracerBlinkTransactionManager(db *sqlx.DB) *tracerBlinkTransactionManager {
	return &tracerBlinkTransactionManager{
		txManager:  database.NewPostgresTxManager(db),
		tracerRepo: newTracerRepository(db),
		blinkRepo:  newBlinkRepository(db),
	}
}

func (tm *tracerBlinkTransactionManager) executeInTransaction(
	ctx context.Context,
	fn func(*tracerRepository, *blinkRepository) error,
) error {
	ctx, span := telemetry.StartSpan(ctx, "blink.tracerBlinkTransactionManager.executeInTransaction")
	defer span.End()

	return tm.txManager.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		tracerRepoWithTx := tm.tracerRepo.withTx(tx)
		blinkRepoWithTx := tm.blinkRepo.withTx(tx)

		return fn(tracerRepoWithTx, blinkRepoWithTx)
	})
}
