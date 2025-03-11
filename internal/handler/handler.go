package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"go-api/internal/mapper"
	"go-api/internal/model"
	"go-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	flightUC *usecase.FlightUsecase
}

var errNotFound = errors.New("not found")

func New(flightUC *usecase.FlightUsecase) *Handler {
	return &Handler{flightUC: flightUC}
}

func (h *Handler) GetAllFlightsHandler(c *gin.Context) {
	flightsDTO, err := h.flightUC.GetAllFlights(c)
	if err != nil {
		if errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No flights found"})
			return
		}
		slog.Error("Failed to fetch flights", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}

	flightsResp := make([]model.FlightResponse, len(flightsDTO))
	for i, flight := range flightsDTO {
		flightsResp[i] = mapper.ToFlightResponse(flight)
	}

	c.JSON(http.StatusOK, flightsResp)
}

func (h *Handler) GetFlightHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flight ID must be greater than 0"})
		return
	}
	flightDTO, err := h.flightUC.GetFlightByID(c, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
			return
		}
		slog.Error("Failed to fetch flight", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flight"})
		return
	}

	flightResp := mapper.ToFlightResponse(*flightDTO)

	c.JSON(http.StatusOK, flightResp)
}

func (h *Handler) InsertFlightHandler(c *gin.Context) {
	var req model.FlightRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	flightDTO := mapper.ToFlightDTO(req)
	err := h.flightUC.InsertFlight(c, &flightDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flight"})
		return
	}

	resp := mapper.ToFlightResponse(flightDTO)
	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": resp})
}

func (h *Handler) UpdateFlightHandler(c *gin.Context) {
	var req model.FlightRequest

	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil || id <= 0 { // ✅ Объединяем проверки ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	flightDTO := mapper.ToFlightDTOWithID(req, id)

	err = h.flightUC.UpdateFlight(c, &flightDTO)
	if err != nil {
		if errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
			return
		}
		slog.Error("Failed to update flight", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update flight"})
		return
	}

	resp := mapper.ToFlightResponse(flightDTO)
	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": resp})
}

func (h *Handler) DeleteFlightHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	err = h.flightUC.DeleteFlight(c, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
			return
		}
		slog.Error("Failed to delete flight", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete flight"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight successfully deleted"})
}
