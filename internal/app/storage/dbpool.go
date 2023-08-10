package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

type InitDBFunc func(conn *pgxpool.Conn) error
