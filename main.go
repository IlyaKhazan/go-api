package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"go-api/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Flight struct {
	FlightID        int    `json:"id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
}

var dbConn *pgx.Conn

func GetAllFlightsFromDB() ([]Flight, error) {
	rows, err := dbConn.Query(context.Background(), "SELECT id, destination_from, destination_to FROM public.flights WHERE deleted_at IS NULL")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch flights from database")
	}
	defer rows.Close()

	var flights []Flight
	for rows.Next() {
		var flight Flight
		if err := rows.Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo); err != nil {
			return nil, errors.Wrap(err, "error scanning flights")
		}
		flights = append(flights, flight)
	}

	if len(flights) == 0 {
		return nil, pgx.ErrNoRows // ✅ Return `pgx.ErrNoRows` if no active flights exist
	}

	return flights, nil
}

func GetFlightByIDFromDB(id int) (*Flight, error) {
	var flight Flight
	err := dbConn.QueryRow(context.Background(),
		"SELECT id, destination_from, destination_to FROM public.flights WHERE id=$1 AND deleted_at IS NULL", id).
		Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows // ✅ Return `pgx.ErrNoRows` for deleted or non-existing flights
		}
		return nil, errors.Wrap(err, "failed to fetch flight")
	}

	return &flight, nil
}

// InsertFlightToDB inserts a new flight into the database
func InsertFlightToDB(flight *Flight) error {
	err := dbConn.QueryRow(context.Background(),
		"INSERT INTO flights (destination_from, destination_to) VALUES ($1, $2) RETURNING id",
		flight.DestinationFrom, flight.DestinationTo).Scan(&flight.FlightID)

	if err != nil {
		return errors.Wrap(err, "failed to insert flight")
	}

	return nil
}

// UpdateFlightInDB updates a flight in the database
func UpdateFlightInDB(flight *Flight) error {
	result, err := dbConn.Exec(context.Background(),
		"UPDATE flights SET destination_from=$1, destination_to=$2 WHERE id=$3",
		flight.DestinationFrom, flight.DestinationTo, flight.FlightID)

	if err != nil {
		return errors.Wrap(err, "failed to update flight")
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // ✅ Handle "no rows updated" correctly
	}

	return nil
}

// DeleteFlightFromDB removes a flight by ID
func DeleteFlightFromDB(id int) error {
	result, err := dbConn.Exec(context.Background(), "UPDATE public.flights SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)

	if err != nil {
		return errors.Wrap(err, "failed to soft delete flight")
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // ✅ Return `pgx.ErrNoRows` if flight is already deleted or doesn't exist
	}

	return nil
}

// Handlers

func getAllFlightsHandler(c *gin.Context) {
	flights, err := GetAllFlightsFromDB()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No flights found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}

	c.JSON(http.StatusOK, flights)
}

func getFlightHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"}) // ✅ 400 Bad Request
		return
	}

	flight, err := GetFlightByIDFromDB(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"}) // ✅ 404 Not Found (for both missing & deleted flights)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve flight"}) // ✅ 500 Internal Server Error
		return
	}

	c.JSON(http.StatusOK, flight) // ✅ 200 OK (only returns non-deleted flights)
}

func insertFlightHandler(c *gin.Context) {
	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	if err := InsertFlightToDB(&flight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flight"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

func updateFlightHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	flight.FlightID = id
	err = UpdateFlightInDB(&flight)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"}) // ✅ 404 Not Found
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update flight"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": flight})
}

func deleteFlightHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("flight_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	err = DeleteFlightFromDB(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found or already deleted"}) // ✅ Return 404 if no action was taken
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete flight"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight successfully marked as deleted"}) // ✅ Return 200 instead of 204
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn, err = config.GetDBConnect(cfg)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	slog.Info("Application started successfully!")
	router := getRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func getRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/flights", getAllFlightsHandler)
	router.GET("/flights/:flight_id", getFlightHandler)
	router.POST("/flights", insertFlightHandler)
	router.PUT("/flights/:flight_id", updateFlightHandler)
	router.DELETE("/flights/:flight_id", deleteFlightHandler)

	return router
}
