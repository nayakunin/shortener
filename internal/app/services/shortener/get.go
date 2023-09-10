package shortener

import (
	"github.com/pkg/errors"
)

// Get returns link by key
func (s *Service) Get(key string) (string, error) {
	link, err := s.Storage.Get(key)
	if err != nil {
		return "", errors.Wrap(err, "unable to get link")
	}

	return link, nil
}
