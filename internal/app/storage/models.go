package storage

type BatchInput struct {
	OriginalURL   string
	CorrelationID string
}

type BatchOutput struct {
	Key           string
	CorrelationID string
}

type Link struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"-"`
}
