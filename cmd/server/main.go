package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tinygo/internal/config"
	"tinygo/internal/database"
	"tinygo/internal/logger"
	"tinygo/internal/shortener"
	"tinygo/internal/storage"
	httphandler "tinygo/internal/transport/http"
)

func main() {
	cfg, err := config.LoadWithViper()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	logger.Init(cfg.LogLevel, cfg.LogFormat)
	logger.Log.Info("starting tinygo server",
		"addr", cfg.Addr,
		"base_url", cfg.BaseURL,
		"database", cfg.Database.Driver,
	)

	// Initialize database
	if err := database.Init(cfg.Database); err != nil {
		logger.Log.Fatalf("init database: %v", err)
	}
	defer database.Close()

	// Create store
	store := storage.NewGormStore()

	svc := shortener.NewService(store, cfg.BaseURL, cfg.CodeLength)
	router := httphandler.NewMux(svc)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Log.Info("listening on server", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Errorf("shutdown error: %v", err)
	}
}
