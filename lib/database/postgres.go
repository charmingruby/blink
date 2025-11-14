package database

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const postgresDriver = "postgres"

var ErrQueryPreparation = errors.New("query preparation error")

func NewPostgresClient(ctx context.Context, url string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, postgresDriver, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
