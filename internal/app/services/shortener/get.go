package shortener

import (
	"github.com/pkg/errors"
)

func (s *Service) Get(key string) (string, error) {
	link, err := s.Storage.Get(key)
	if err != nil {
		return "", errors.Wrap(err, "unable to get link")
	}

	return link, nil
}
