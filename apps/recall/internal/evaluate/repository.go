package evaluate

import (
	"blink/lib/core"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type tracerRepo struct {
	db *sqlx.DB
}

func newTracerRepo(db *sqlx.DB) *tracerRepo {
	return &tracerRepo{db: db}
}

func (r *tracerRepo) findTracerByIP(ctx context.Context, ip string) (core.Tracer, error) {
	ctx, stop := context.WithTimeout(ctx, 5*time.Second)
	defer stop()

	query := "SELECT * FROM tracers WHERE ip = $1"

	var tracer core.Tracer
	if err := r.db.GetContext(ctx, &tracer, query, ip); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Tracer{}, nil
		}

		return core.Tracer{}, err
	}

	return tracer, nil
}
