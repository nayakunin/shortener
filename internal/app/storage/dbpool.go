package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool is an interface for pgxpool.Pool
type DBPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

// InitDBFunc is a function that is used to initialize database
type InitDBFunc func(conn *pgxpool.Conn) error
