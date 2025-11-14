package evaluate

import (
	"blink/lib/core"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

func findTracerByIP(ctx context.Context, db *sqlx.DB, ip string) (core.Tracer, error) {
	ctx, stop := context.WithTimeout(ctx, 5*time.Second)
	defer stop()

	query := "SELECT * FROM tracers WHERE ip = $1"

	var tracer core.Tracer
	if err := db.GetContext(ctx, &tracer, query, ip); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Tracer{}, nil
		}

		return core.Tracer{}, err
	}

	return tracer, nil
}

func storeTracer(ctx context.Context, db *sqlx.DB, tr core.Tracer) error {
	ctx, stop := context.WithTimeout(ctx, 10*time.Second)
	defer stop()

	query := "INSERT INTO tracers (id, ip, total_blinks, last_blink_at, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := db.ExecContext(
		ctx,
		query,
		tr.ID,
		tr.IP,
		tr.TotalBlinks,
		tr.LastBlinkAt,
		tr.UpdatedAt,
		tr.CreatedAt,
	)

	return err
}
