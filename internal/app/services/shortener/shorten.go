package shortener

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

var (
	ErrInvalidURL    = errors.New("invalid url")
	ErrShortenFailed = errors.New("unable to shorten url")
)

func (s *Service) Shorten(userID string, urlString string) (string, error) {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return "", ErrInvalidURL
	}

	key, err := s.Storage.Add(urlString, userID)
	shortURL := fmt.Sprintf("%s/%s", s.Cfg.BaseURL, key)
	if err != nil {
		return shortURL, errors.Wrap(err, ErrShortenFailed.Error())
	}

	return shortURL, nil
}
