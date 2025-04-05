package app

import (
	"go-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func GetRouter(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	flights := router.Group("/flights")
	flights.GET("", handler.GetAllFlights)
	flights.GET("/:flight_id", handler.GetFlight)
	flights.POST("", handler.InsertFlight)
	flights.PUT("/:flight_id", handler.UpdateFlight)
	flights.DELETE("/:flight_id", handler.DeleteFlight)

	return router
}
