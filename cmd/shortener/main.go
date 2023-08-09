package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/nayakunin/shortener/internal/app/server"
	"github.com/nayakunin/shortener/internal/app/server/config"
	storagePackage "github.com/nayakunin/shortener/internal/app/storage"
)

var (
	BuildVersion string
	BuildTime    string
	BuildCommit  string
)

const NA = "N/A"

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if BuildVersion == "" {
		BuildVersion = NA
	}

	if BuildTime == "" {
		BuildTime = NA
	}

	if BuildCommit == "" {
		BuildCommit = NA
	}

	fmt.Printf("Build Version: %s\n", BuildVersion)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Build Commit: %s\n", BuildCommit)

	storage, err := storagePackage.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := server.NewRouter(*cfg, storage)
	pprof.Register(r)

	log.Fatal(r.Run(cfg.ServerAddress))
}
