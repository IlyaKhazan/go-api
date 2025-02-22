package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	flightUC usecase.FlightProvider
}

func New(flightUC usecase.FlightProvider) *Handler {
	return &Handler{flightUC: flightUC}
}

func (h *Handler) GetAllFlightsHandler(c *gin.Context) {
	ctx := context.Background()

	flights, err := h.flightUC.GetAllFlights(ctx)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No flights found"})
			return
		}
		slog.Error("Failed to fetch flights", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}

	c.JSON(http.StatusOK, flights)
}

func (h *Handler) GetFlightHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	flight, err := h.flightUC.GetFlightByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, flight)
}

func (h *Handler) InsertFlightHandler(c *gin.Context) {
	ctx := context.Background()

	var flight usecase.Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	err := h.flightUC.InsertFlight(ctx, &flight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flight"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

func (h *Handler) UpdateFlightHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight usecase.Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	flight.FlightID = id
	err = h.flightUC.UpdateFlight(ctx, &flight)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": flight})
}

func (h *Handler) DeleteFlightHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	err = h.flightUC.DeleteFlight(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight successfully marked as deleted"})
}
