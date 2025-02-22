package main

import (
	"context"
	"log/slog"
	"os"

	"go-api/config"
	"go-api/internal/handler"
	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
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

	flightUC := usecase.New(dbConn)

	h := handler.New(flightUC)

	router := gin.Default()
	router.GET("/flights", h.GetAllFlightsHandler)
	router.GET("/flights/:flight_id", h.GetFlightHandler)
	router.POST("/flights", h.InsertFlightHandler)
	router.PUT("/flights/:flight_id", h.UpdateFlightHandler)
	router.DELETE("/flights/:flight_id", h.DeleteFlightHandler)

	slog.Info("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
