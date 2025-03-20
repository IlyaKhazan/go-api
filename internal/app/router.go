package app

import (
	"go-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func GetRouter(handler *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.GET("/flights", handler.GetAllFlights)
	router.GET("/flights/:flight_id", handler.GetFlight)
	router.POST("/flights", handler.InsertFlight)
	router.PUT("/flights/:flight_id", handler.UpdateFlight)
	router.DELETE("/flights/:flight_id", handler.DeleteFlight)

	return router
}
