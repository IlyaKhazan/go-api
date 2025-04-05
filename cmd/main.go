package main

import (
	"context"
	"log/slog"
	"os"

	"go-api/config"
	"go-api/internal/app"
	"go-api/internal/cache"
	"go-api/internal/handler"
	"go-api/internal/repository"
	"go-api/internal/storage"
	"go-api/internal/usecase"
)

func main() {
	slog.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	dbConn, err := storage.GetDBConnect(cfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	flightsRepo := repository.NewFlightRepository(dbConn)
	cacheDecorator := cache.NewCacheDecorator(flightsRepo)
	flightUC := usecase.NewFlightUsecase(cacheDecorator)
	handle := handler.New(flightUC)
	router := app.GetRouter(handle)

	slog.Info("Server running on port", "address", cfg.Address)
	if err := router.Run(cfg.Address); err != nil {
		slog.Error("failed to start server", "error", err)
		return
	}
}
