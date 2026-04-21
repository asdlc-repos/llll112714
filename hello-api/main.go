package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hello-api/internal/handlers"
	"hello-api/internal/middleware"
)

const (
	defaultAddr     = ":9090"
	shutdownTimeout = 10 * time.Second
)

func main() {
	addr := defaultAddr
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handlers.Hello)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/", handlers.NotFound)

	handler := middleware.Recovery(middleware.Logging(mux))

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("hello-api listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Printf("received signal %s, shutting down", sig)
	case err := <-errCh:
		log.Printf("server error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		os.Exit(1)
	}
	log.Println("server stopped")
}
