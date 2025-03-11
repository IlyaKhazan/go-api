package main

import (
	"context"
	"log/slog"
	"os"

	"go-api/config"
	"go-api/internal/app"
	"go-api/internal/handler"
	"go-api/internal/repository"
	"go-api/internal/usecase"
)

func main() {
	slog.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	dbConn, err := config.GetDBConnect(cfg)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	flightRepo := repository.NewFlightRepository(dbConn)
	flightUC := usecase.NewFlightUsecase(flightRepo)
	handle := handler.New(flightUC)

	router := app.GetRouter(handle)

	slog.Info("Server running on port", "address", cfg.Address)
	if err := router.Run(cfg.Address); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
