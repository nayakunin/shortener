package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
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

	var wg sync.WaitGroup

	r, m := server.NewRouter(*cfg, storage, &wg)
	pprof.Register(r)

	srv := &http.Server{
		Addr:      cfg.ServerAddress,
		Handler:   r,
		TLSConfig: m.TLSConfig(),
	}

	// Setting up a channel to listen to OS signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-quit // blocking until a signal is received

		fmt.Println("Shutting down server...")

		// Context with a timeout to ensure all requests are processed
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		go func() {
			wg.Wait() // Wait for all handlers to complete
			cancel()  // Cancel the context when all handlers are done
		}()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("Server forced to shutdown: %v\n", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
