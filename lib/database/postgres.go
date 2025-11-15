package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const postgresDriver = "postgres"

var ErrQueryPreparation = errors.New("query preparation error")

type Querier interface {
	sqlx.Queryer
	sqlx.Execer
	sqlx.Preparer
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

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

type PostgresTxManager struct {
	db *sqlx.DB
}

func NewPostgresTxManager(db *sqlx.DB) *PostgresTxManager {
	return &PostgresTxManager{db: db}
}

func (tm *PostgresTxManager) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := tm.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
