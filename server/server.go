package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	gin.DisableConsoleColor()
	// gin.ForceConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/port-registrations", PortRegistrationsHandler)
	// r.POST("/old-port-registrations", HandlerOld)

	return r
}

func Start() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
