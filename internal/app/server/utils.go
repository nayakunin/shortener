package server

import "fmt"

func (s Server) formatURL(key string) string {
	return fmt.Sprintf("%s/%s", s.Cfg.BaseURL, key)
}
