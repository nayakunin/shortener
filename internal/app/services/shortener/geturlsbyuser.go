package shortener

import (
	"fmt"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/pkg/errors"
)

var (
	ErrNoUrlsFound = errors.New("no urls found")
)

func (s *Service) GetUrlsByUser(userID string) ([]interfaces.Link, error) {
	urls, err := s.Storage.GetUrlsByUser(userID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get urls by user")
	}

	if len(urls) == 0 {
		return nil, ErrNoUrlsFound
	}

	response := make([]interfaces.Link, 0, len(urls))
	for key, originalURL := range urls {
		response = append(response, interfaces.Link{
			ShortURL:    fmt.Sprintf("%s/%s", s.Cfg.BaseURL, key),
			OriginalURL: originalURL,
		})
	}

	return response, nil
}
