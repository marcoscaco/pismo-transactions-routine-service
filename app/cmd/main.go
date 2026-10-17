package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with health check endpoint.
// This function is designed to be testable by separating router setup from main execution.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

func main() {
	router := SetupRouter()
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
