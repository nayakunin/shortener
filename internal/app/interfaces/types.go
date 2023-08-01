package interfaces

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

// Storage is a storage interface that implements all methods for working with storage
type Storage interface {
	Get(key string) (string, error)
	Add(link string, userID string) (string, error)
	AddBatch(batch []BatchInput, userID string) ([]BatchOutput, error)
	GetUrlsByUser(userID string) (map[string]string, error)
	DeleteUserUrls(userID string, keys []string) error
}
