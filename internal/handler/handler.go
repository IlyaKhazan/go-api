package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"go-api/internal/apperr"
	"go-api/internal/mapper"
	"go-api/internal/model"
	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	flightUC usecase.FlightProvider
}

func New(flightUC usecase.FlightProvider) *Handler {
	return &Handler{flightUC: flightUC}
}

func (h *Handler) GetAllFlights(c *gin.Context) {
	flightsDTO, err := h.flightUC.GetAllFlights(c)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		slog.Error("failed to fetch flights", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch flights"})
		return
	}

	flightsResp := make([]model.FlightResponse, len(flightsDTO))
	for i, flight := range flightsDTO {
		flightsResp[i] = mapper.ToFlightResponse(&flight)
	}

	c.JSON(http.StatusOK, flightsResp)
}

func (h *Handler) GetFlight(c *gin.Context) {
	id, err := uuid.FromString(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight ID"})
		return
	}
	flightDTO, err := h.flightUC.GetFlightByID(c, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
			return
		}
		slog.Error("failed to fetch flight", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create flight"})
		return
	}

	c.JSON(http.StatusOK, mapper.ToFlightResponse(flightDTO))
}

func (h *Handler) InsertFlight(c *gin.Context) {
	var req model.FlightRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight data"})
		return
	}

	flightDTO := mapper.ToFlightDTO(req)
	err := h.flightUC.InsertFlight(c, &flightDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create flight"})
		return
	}

	resp := mapper.ToFlightResponse(&flightDTO)
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateFlight(c *gin.Context) {
	var req model.FlightRequest

	id, err := uuid.FromString(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight data"})
		return
	}

	flightDTO := mapper.ToFlightDTOWithID(req, id)

	if err := h.flightUC.UpdateFlight(c, &flightDTO); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
			return
		}
		slog.Error("failed to update flight", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update flight"})
		return
	}

	resp := mapper.ToFlightResponse(&flightDTO)
	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": resp})
}

func (h *Handler) DeleteFlight(c *gin.Context) {
	id, err := uuid.FromString(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fnvalid flight ID"})
		return
	}

	err = h.flightUC.DeleteFlight(c, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		slog.Error("failed to delete flight", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete flight"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "flight successfully deleted"})
}
