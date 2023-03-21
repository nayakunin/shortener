package storage

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nayakunin/shortener/internal/app/utils"
)

type DBStorage struct {
	sync.Mutex
	Connection *pgx.Conn
}

func newDBStorage(databaseURL string) (*DBStorage, error) {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(context.Background(), `SELECT 1 FROM links LIMIT 1`)
	if err != nil {
		_, err := conn.Exec(context.Background(), `CREATE TABLE links (
			id SERIAL PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			original_url VARCHAR(255) UNIQUE NOT NULL,
			user_id VARCHAR(255) NOT NULL
		)`)
		if err != nil {
			return nil, err
		}
	}

	return &DBStorage{
		Connection: conn,
	}, nil
}

func (s *DBStorage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result string
	err := s.Connection.QueryRow(ctx, "SELECT original_url FROM links WHERE key = $1", key).Scan(&result)
	if err != nil {
		return "", false
	}

	return result, true
}

func (s *DBStorage) Add(link string, userID string) (string, error) {
	s.Lock()
	defer s.Unlock()

	key := utils.Encode(link)

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	ctx := context.Background()

	_, err := s.Connection.Exec(ctx, "INSERT INTO links (key, original_url, user_id) VALUES ($1, $2, $3)", key, link, userID)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == pgerrcode.UniqueViolation {
				var prevKey string
				err := s.Connection.QueryRow(ctx, "SELECT key FROM links WHERE original_url = $1", link).Scan(&prevKey)
				if err != nil {
					return "", err
				}

				return prevKey, ErrKeyExists
			}
		}
		return "", err
	}

	return key, nil
}

func (s *DBStorage) AddBatch(batches []BatchInput, userID string) ([]BatchOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.Connection.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	stmt, err := tx.Prepare(ctx, "insert", "INSERT INTO links (key, original_url, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}

	output := make([]BatchOutput, len(batches))
	for i, linkObject := range batches {
		if _, err := url.ParseRequestURI(linkObject.OriginalURL); err != nil {
			return nil, ErrBatchInvalidURL
		}

		key := utils.Encode(linkObject.OriginalURL)

		_, err = tx.Exec(ctx, stmt.Name, key, linkObject.OriginalURL, userID)
		if err != nil {
			return nil, err
		}

		output[i] = BatchOutput{
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

func (s *DBStorage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	links := make(map[string]string)
	err := s.Connection.QueryRow(ctx, "SELECT key, original_url FROM links WHERE user_id = $1", id).Scan(links)
	if err != nil {
		return nil, err
	}

	return links, nil
}
