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
		mockInitDb := func(conn *pgxpool.Conn) error {
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
		res, err := newDBStorage(mockPool, mockInitDb, mockRequestBuffer)
		assert.NoError(t, err)
		assert.Equal(t, &DBStorage{
			Pool:          mockPool,
			requestBuffer: mockRequestBuffer,
		}, res)
	})

	//t.Run("Add and Get", func(t *testing.T) {
	//	storage := newDBStorage()
	//	key, err := storage.Add("https://www.example.com", "user1")
	//	assert.NoError(t, err)
	//
	//	retrievedLink, err := storage.Get(key)
	//	assert.NoError(t, err)
	//	assert.Equal(t, "https://www.example.com", retrievedLink)
	//})
	//
	//t.Run("AddBatch", func(t *testing.T) {
	//	storage := newDBStorage()
	//	batches := []interfaces.BatchInput{
	//		{OriginalURL: "https://www.example1.com"},
	//		{OriginalURL: "https://www.example2.com"},
	//	}
	//
	//	output, err := storage.AddBatch(batches, "user1")
	//	assert.NoError(t, err)
	//	assert.Len(t, output, 2)
	//})
	//
	//t.Run("DeleteUserUrls", func(t *testing.T) {
	//	storage := newDBStorage()
	//
	//	key, err := storage.Add("https://www.example.com", "user1")
	//	assert.NoError(t, err)
	//
	//	err = storage.DeleteUserUrls("user1", []string{key})
	//	assert.NoError(t, err)
	//
	//	_, err = storage.Get(key)
	//	assert.Error(t, err)
	//})
}
