package storage

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nayakunin/shortener/internal/app/utils"
)

const Timeout = 5 * time.Second
const MaxRequests = 100
const DeleteRequestsTimeout = 3 * time.Second

type RequestBatch struct {
	UserID string
	Keys   []string
}

type RequestBuffer struct {
	buffer      chan RequestBatch
	maxRequests int
	timer       *time.Timer
	startChan   chan struct{}
}

type DBStorage struct {
	Connection    *pgx.Conn
	requestBuffer *RequestBuffer
}

func newRequestBuffer(maxRequests int) *RequestBuffer {
	return &RequestBuffer{
		buffer:      make(chan RequestBatch, maxRequests),
		maxRequests: maxRequests,
		timer:       time.NewTimer(DeleteRequestsTimeout),
		startChan:   make(chan struct{}),
	}
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

	db := DBStorage{
		Connection:    conn,
		requestBuffer: newRequestBuffer(MaxRequests),
	}

	go func() {
		for {
			fmt.Println("Waiting for requests")
			select {
			case <-db.requestBuffer.timer.C:
				fmt.Println("Timeout")
				db.requestBuffer.timer.Reset(DeleteRequestsTimeout)
				db.processDeleteRequests()
			case <-db.requestBuffer.startChan:
				db.requestBuffer.timer.Reset(DeleteRequestsTimeout)
				db.processDeleteRequests()
			}
		}
	}()

	return &db, nil

}

func (s *DBStorage) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	links := make(map[string]string)
	err := s.Connection.QueryRow(ctx, "SELECT key, original_url FROM links WHERE user_id = $1", id).Scan(links)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (s *DBStorage) DeleteUserUrls(userID string, keys []string) error {
	s.requestBuffer.AddRequest(userID, keys)
	return nil
}

func (s *DBStorage) processDeleteRequests() {
	fmt.Println("Processing requests", len(s.requestBuffer.buffer))
	if len(s.requestBuffer.buffer) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	requests := s.requestBuffer.GetRequests()
	fmt.Println("Requests", requests)

	for _, batch := range requests {
		_, err := s.Connection.Exec(ctx, "UPDATE links SET is_deleted = TRUE WHERE user_id = $1 AND key = ANY($2)", batch.UserID, batch.Keys)
		if err != nil {
			return
		}
	}
}

func (rb *RequestBuffer) AddRequest(userID string, keys []string) {
	rb.buffer <- RequestBatch{
		UserID: userID,
		Keys:   keys,
	}

	if len(rb.buffer) == rb.maxRequests-1 {
		rb.timer.Reset(DeleteRequestsTimeout)
		rb.startChan <- struct{}{}
	}
}

func (rb *RequestBuffer) GetRequests() []RequestBatch {
	requests := make([]RequestBatch, 0, rb.maxRequests)

	for i := 0; i < rb.maxRequests; i++ {
		select {
		case request := <-rb.buffer:
			requests = append(requests, request)
		default:
			return requests
		}
	}

	return requests
}
