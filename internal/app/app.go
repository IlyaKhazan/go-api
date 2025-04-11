package app

import (
	"fmt"
	"log/slog"
	"time"

	"go-api/config"
	"go-api/database"
	"go-api/internal/cache"
	"go-api/internal/handler"
	"go-api/internal/repository"
	"go-api/internal/storage"
	"go-api/internal/usecase"
)

func Run() error {
	slog.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return fmt.Errorf("failed to load config: %w", err)
	}
	dbURL := cfg.GetDBConnString()
	if err = database.Migrate(dbURL); err != nil {
		slog.Error("database migration failed", "error", err)
		return fmt.Errorf("database migration failed: %w", err)
	}

	dbConn, err := storage.GetDBConnect(dbURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer func() {
		slog.Info("closing database connection")
		dbConn.Close()
	}()

	flightsRepo := repository.NewFlightRepository(dbConn)
	cacheDecorator := cache.NewDecorator(flightsRepo, 10*time.Second)
	cacheDecorator.StartCleanup(1 * time.Minute)
	flightUC := usecase.NewFlightUsecase(cacheDecorator)
	handle := handler.New(flightUC)
	router := GetRouter(handle)

	slog.Info("Server running on port", "address", cfg.Address)
	if err = router.Run(cfg.Address); err != nil {
		slog.Error("failed to start server", "error", err)
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
