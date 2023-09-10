package shortener

import (
	"fmt"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/pkg/errors"
)

func (s *Service) ShortenBatch(userID string, input []interfaces.BatchInput) ([]interfaces.BatchOutput, error) {
	output, err := s.Storage.AddBatch(input, userID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to shorten batch")
	}

	response := make([]interfaces.BatchOutput, len(output))
	for i, v := range output {
		response[i] = interfaces.BatchOutput{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", s.Cfg.BaseURL, v.Key),
		}
	}

	return response, nil
}
