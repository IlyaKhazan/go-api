package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

type Flight struct {
	FlightID        int    `json:"id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
}

var dbConn *pgx.Conn

func getConnect(urlString string) (*pgx.Conn, error) {
	fmt.Printf("Connecting to Postgres: %s\n", urlString)

	conn, err := pgx.Connect(context.Background(), urlString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}

func main() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	dbConn, err = getConnect(dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close(context.Background())

	log.Println("Connected to PostgreSQL successfully!")
	router := getRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func getRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/flights", getAllDBFlights)
	router.GET("/flights/:flight_id", getDBFlight)
	router.POST("/flights", insertDBFlight)
	router.PUT("/flights/:flight_id", updateDBFlight)
	router.DELETE("/flights/:flight_id", deleteDBFlight)

	return router
}

func getAllDBFlights(c *gin.Context) {
	rows, err := dbConn.Query(context.Background(), "SELECT id, destination_from, destination_to FROM public.flights")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}
	defer rows.Close()

	var flights []Flight
	for rows.Next() {
		var flight Flight
		if err := rows.Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning flights"})
			return
		}
		flights = append(flights, flight)
	}

	c.JSON(http.StatusOK, flights)
}

func getDBFlight(c *gin.Context) {
	idStr := c.Param("flight_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flight ID is invalid"})
		return
	}

	var flight Flight
	err = dbConn.QueryRow(context.Background(), "SELECT id, destination_from, destination_to FROM public.flights WHERE id=$1", id).
		Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, flight)
}

func insertDBFlight(c *gin.Context) {
	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dbConn.QueryRow(context.Background(),
		"INSERT INTO public.flights (destination_from, destination_to) VALUES ($1, $2) RETURNING id",
		flight.DestinationFrom, flight.DestinationTo).Scan(&flight.FlightID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flight"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flight successfully created", "flight": flight})
}

func updateDBFlight(c *gin.Context) {
	idStr := c.Param("flight_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flight ID is invalid"})
		return
	}

	var updatedFlight Flight
	if err := c.ShouldBindJSON(&updatedFlight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data"})
		return
	}

	result, err := dbConn.Exec(context.Background(),
		"UPDATE public.flights SET destination_from=$1, destination_to=$2 WHERE id=$3",
		updatedFlight.DestinationFrom, updatedFlight.DestinationTo, id)

	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	updatedFlight.FlightID = id
	c.JSON(http.StatusOK, gin.H{"message": "Flight updated successfully", "flight": updatedFlight})
}

func deleteDBFlight(c *gin.Context) {
	idStr := c.Param("flight_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flight ID is invalid"})
		return
	}

	result, err := dbConn.Exec(context.Background(), "DELETE FROM public.flights WHERE id=$1", id)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.Status(http.StatusNoContent) // ✅ Return 204 No Content on success
}
