package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/autotls"
	"github.com/nayakunin/shortener/internal/app/server"
	"github.com/nayakunin/shortener/internal/app/server/config"

	storagePackage "github.com/nayakunin/shortener/internal/app/storage"
)

// Build information. Populated at build-time.
var (
	// BuildVersion is a version of the build
	BuildVersion string
	// BuildTime is a time of the build
	BuildTime string
	// BuildCommit is a git commit of the build
	BuildCommit string
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if BuildVersion == "" {
		BuildVersion = "N/A"
	}

	if BuildTime == "" {
		BuildTime = "N/A"
	}

	if BuildCommit == "" {
		BuildCommit = "N/A"
	}

	fmt.Printf("Build Version: %s\n", BuildVersion)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Build Commit: %s\n", BuildCommit)

	storage, err := storagePackage.NewStorage(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	r, manager := server.NewRouter(*cfg, storage)
	pprof.Register(r)

	log.Fatal(autotls.RunWithManager(r, &manager))
}
