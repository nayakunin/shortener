package server

import (
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

type Server struct {
	cfg     config.Config
	storage storage.Storager
}

func NewServer(cfg config.Config, s storage.Storager) Server {
	return Server{
		cfg,
		s,
	}
}
