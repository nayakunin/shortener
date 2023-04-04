package storage

import (
	"context"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nayakunin/shortener/internal/app/utils"
)

const TIMEOUT = 5 * time.Second

type DBStorage struct {
	Connection *pgx.Conn
}

func initDB(conn *pgx.Conn) error {
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

func newDBStorage(databaseURL string) (*DBStorage, error) {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	err = initDB(conn)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		Connection: conn,
	}, nil
}

func (s *DBStorage) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	type Link struct {
		OriginalURL string
		IsDeleted   bool
	}

	var result Link
	err := s.Connection.QueryRow(ctx, "SELECT original_url, is_deleted FROM links WHERE key = $1", key).Scan(&result)
	if err != nil {
		return "", err
	}

	if result.IsDeleted {
		return "", ErrKeyDeleted
	}

	return result.OriginalURL, nil
}

func (s *DBStorage) Add(link string, userID string) (string, error) {
	key := utils.Encode(link)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	res, err := s.Connection.Exec(ctx, "INSERT INTO links (key, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING", key, link, userID)
	if err != nil {
		return "", err
	}

	if res.RowsAffected() == 0 {
		var prevKey string
		err := s.Connection.QueryRow(ctx, "SELECT key FROM links WHERE original_url = $1", link).Scan(&prevKey)
		if err != nil {
			return "", err
		}

		return prevKey, ErrKeyExists
	}

	return key, nil
}

func (s *DBStorage) AddBatch(batches []BatchInput, userID string) ([]BatchOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
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

	stmt, err := tx.Prepare(ctx, "insert", "INSERT INTO links (key, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING")
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
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	links := make(map[string]string)
	err := s.Connection.QueryRow(ctx, "SELECT key, original_url FROM links WHERE user_id = $1", id).Scan(links)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (s *DBStorage) DeleteUserUrls(userID string, keys []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	_, err := s.Connection.Exec(ctx, "UPDATE links SET is_deleted = TRUE WHERE user_id = $1 AND key = ANY($2)", userID, keys)
	if err != nil {
		return err
	}

	return nil
}
