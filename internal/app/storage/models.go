package storage

// BatchInput is a batch of input data.
type BatchInput struct {
	OriginalURL   string
	CorrelationID string
}

// BatchOutput is a batch of output data.
type BatchOutput struct {
	Key           string
	CorrelationID string
}

// Link is a link object.
type Link struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"-"`
	IsDeleted   bool   `json:"-"`
}
