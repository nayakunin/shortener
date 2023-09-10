package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/nayakunin/shortener/internal/app/config"
	"github.com/nayakunin/shortener/internal/app/grpc"
	"github.com/nayakunin/shortener/internal/app/rest"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	pb "github.com/nayakunin/shortener/proto"
	grpcCore "google.golang.org/grpc"

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

	shortenerService := shortener.NewShortenerService(*cfg, storage)

	r, m := rest.NewRouter(shortenerService, &wg, cfg.AuthSecret, cfg.TrustedSubnet)
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

		fmt.Println("Shutting down rest...")

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

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpcCore.NewServer()
		pb.RegisterShortenerServer(grpcServer, grpc.NewServer(shortenerService))

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	if err := srv.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
