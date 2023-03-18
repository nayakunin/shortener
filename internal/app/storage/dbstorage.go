package storage

import (
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	Storage
	Connection *pgx.Conn
}

func (s *DBStorage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	return "", true
}

func (s *DBStorage) Add(link string, userID string) (string, error) {
	s.Lock()
	defer s.Unlock()

	return "", nil
}

func (s *DBStorage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	links := make(map[string]string)

	return links, nil
}
