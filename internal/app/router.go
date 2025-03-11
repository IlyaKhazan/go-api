package app

import (
	"go-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func GetRouter(handler *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.GET("/flights", handler.GetAllFlightsHandler)
	router.GET("/flights/:flight_id", handler.GetFlightHandler)
	router.POST("/flights", handler.InsertFlightHandler)
	router.PUT("/flights/:flight_id", handler.UpdateFlightHandler)
	router.DELETE("/flights/:flight_id", handler.DeleteFlightHandler)

	return router
}
