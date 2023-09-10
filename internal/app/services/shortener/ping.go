package shortener

import (
	"context"

	"github.com/nayakunin/shortener/internal/app/storage"
	"github.com/pkg/errors"
)

var (
	ErrDBPingFailed = errors.New("unable to ping database")
)

func (s *Service) Ping() error {
	dbStorage, ok := s.Storage.(*storage.DBStorage)
	if !ok {
		return nil
	}

	if err := dbStorage.Pool.Ping(context.Background()); err != nil {
		return errors.Wrap(err, ErrDBPingFailed.Error())
	}

	return nil
}
