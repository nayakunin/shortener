package shortener

import "github.com/pkg/errors"

func (s *Service) DeleteUserUrls(userID string, keys []string) error {
	if err := s.Storage.DeleteUserUrls(userID, keys); err != nil {
		return errors.Wrap(err, "unable to delete user urls")
	}

	return nil
}
