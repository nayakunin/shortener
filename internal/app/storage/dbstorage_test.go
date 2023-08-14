package storage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

type MockPool struct {
	BeginFunc   func(ctx context.Context) (pgx.Tx, error)
	AcquireFunc func(ctx context.Context) (*pgxpool.Conn, error)
	PingFunc    func(ctx context.Context) error
}

func (m *MockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.BeginFunc(ctx)
}

func (m *MockPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return m.AcquireFunc(ctx)
}

func (m *MockPool) Ping(ctx context.Context) error {
	return nil
}

func TestDBStorage(t *testing.T) {
	t.Run("newDBStorage", func(t *testing.T) {
		mockInitDB := func(conn *pgxpool.Conn) error {
			return nil
		}
		mockPool := &MockPool{
			BeginFunc: func(ctx context.Context) (pgx.Tx, error) {
				return nil, nil
			},
			AcquireFunc: func(ctx context.Context) (*pgxpool.Conn, error) {
				return nil, nil
			},
			PingFunc: func(ctx context.Context) error {
				return nil
			},
		}
		mockRequestBuffer := newRequestBuffer(MaxRequests)
		res, err := newDBStorage(mockPool, mockInitDB, mockRequestBuffer)
		assert.NoError(t, err)
		assert.Equal(t, &DBStorage{
			Pool:          mockPool,
			requestBuffer: mockRequestBuffer,
		}, res)
	})
}
