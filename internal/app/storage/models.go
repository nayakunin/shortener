package storage

type Link struct {
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"original_url"`
	UserId   string
}
