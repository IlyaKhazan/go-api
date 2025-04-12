package app

import (
	"fmt"
	"time"

	"go-api/internal/handler"
	"go-api/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetRouter(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Next()

		metrics.HTTPRequests.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			fmt.Sprintf("%d", c.Writer.Status()),
		).Inc()
	})

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()

		status := fmt.Sprintf("%d", c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		metrics.HTTPDuration.WithLabelValues(method, path, status).Observe(duration)
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	flights := router.Group("/flights")
	flights.GET("", handler.GetAllFlights)
	flights.GET("/:flight_id", handler.GetFlight)
	flights.POST("", handler.InsertFlight)
	flights.PUT("/:flight_id", handler.UpdateFlight)
	flights.DELETE("/:flight_id", handler.DeleteFlight)

	return router
}
