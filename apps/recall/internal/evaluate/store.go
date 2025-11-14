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
