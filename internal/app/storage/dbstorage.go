package storage

import (
	"context"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/utils"
)

// Timeout is a timeout for all db operations
const Timeout = 5 * time.Second

// DBStorage is a storage based on PostgreSQL
type DBStorage struct {
	Pool          DBPool
	requestBuffer *RequestBuffer
}

func initDB(conn *pgxpool.Conn) error {
	_, err := conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		key VARCHAR(255) NOT NULL,
		original_url VARCHAR(255) UNIQUE NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE
	)`)
	if err != nil {
		return err
	}

	return nil
}

func newDBStorage(pool DBPool, initDB InitDBFunc, requestBuffer *RequestBuffer) (*DBStorage, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	err = initDB(conn)
	if err != nil {
		return nil, err
	}

	db := DBStorage{
		Pool:          pool,
		requestBuffer: requestBuffer,
	}

	go db.requestBufferWorker(context.Background())

	return &db, nil
}

func (s *DBStorage) requestBufferWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.requestBuffer.ticker.C:
			s.processDeleteRequests()
		case <-s.requestBuffer.isBufferFullCh:
			s.requestBuffer.ticker.Reset(DeleteRequestsTimeout)
			s.processDeleteRequests()
		}
	}
}

// Get returns original URL by key
func (s *DBStorage) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	var originalURL string
	var isDeleted bool
	err = conn.QueryRow(ctx, "SELECT original_url, is_deleted FROM links WHERE key = $1", key).Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", err
	}

	if isDeleted {
		return "", ErrKeyDeleted
	}

	return originalURL, nil
}

// Add adds new link to storage
func (s *DBStorage) Add(link string, userID string) (string, error) {
	key := utils.Encode(link)

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, "INSERT INTO links (key, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING", key, link, userID)
	if err != nil {
		return "", err
	}

	if res.RowsAffected() == 0 {
		var prevKey string
		err := conn.QueryRow(ctx, "SELECT key FROM links WHERE original_url = $1", link).Scan(&prevKey)
		if err != nil {
			return "", err
		}

		return prevKey, ErrKeyExists
	}

	return key, nil
}

// AddBatch adds new links to storage
func (s *DBStorage) AddBatch(batches []interfaces.BatchInput, userID string) ([]interfaces.BatchOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	stmt, err := tx.Prepare(ctx, "insert", "INSERT INTO links (key, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING")
	if err != nil {
		return nil, err
	}

	output := make([]interfaces.BatchOutput, len(batches))
	for i, linkObject := range batches {
		if _, err := url.ParseRequestURI(linkObject.OriginalURL); err != nil {
			return nil, ErrBatchInvalidURL
		}

		key := utils.Encode(linkObject.OriginalURL)

		_, err = tx.Exec(ctx, stmt.Name, key, linkObject.OriginalURL, userID)
		if err != nil {
			return nil, err
		}

		output[i] = interfaces.BatchOutput{
			Key:           key,
			CorrelationID: linkObject.CorrelationID,
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetUrlsByUser returns all user's URLs that are not deleted
func (s *DBStorage) GetUrlsByUser(id string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	links := make(map[string]string)
	err = conn.QueryRow(ctx, "SELECT key, original_url FROM links WHERE user_id = $1", id).Scan(links)
	if err != nil {
		return nil, err
	}

	return links, nil
}

// DeleteUserUrls deletes all user's URLs
func (s *DBStorage) DeleteUserUrls(userID string, keys []string) error {
	s.requestBuffer.AddRequest(userID, keys)
	return nil
}

func (s *DBStorage) processDeleteRequests() {
	if len(s.requestBuffer.buffer) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	requests := s.requestBuffer.GetRequests()

	for _, batch := range requests {
		_, err := conn.Exec(ctx, "UPDATE links SET is_deleted = TRUE WHERE user_id = $1 AND key = ANY($2)", batch.UserID, batch.Keys)
		if err != nil {
			return
		}
	}
}
