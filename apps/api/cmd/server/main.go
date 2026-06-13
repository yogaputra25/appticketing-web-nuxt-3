package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ticketing/api/internal/config"
	"github.com/ticketing/api/internal/database"
	"github.com/ticketing/api/internal/router"
	"github.com/ticketing/api/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "database error: %v\n", err)
		os.Exit(1)
	}

	rdb, err := database.NewRedis(cfg.RedisURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "redis error: %v\n", err)
		os.Exit(1)
	}

	// Background workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.StartExpiryJobs(ctx, db, rdb, &worker.Config{
		QueueSessionTTLMinutes: cfg.QueueSessionTTLMinutes,
		WarRateLimitPerMin:     cfg.WarRateLimitPerMin,
	})

	r := router.New(cfg, db, rdb)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("🚀 API running on :%d (env=%s)\n", cfg.AppPort, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nshutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}
