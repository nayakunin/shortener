package shortener

import "github.com/nayakunin/shortener/internal/app/interfaces"

// Stats returns stats
func (s *Service) Stats() (*interfaces.Stats, error) {
	stats, err := s.Storage.Stats()
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
