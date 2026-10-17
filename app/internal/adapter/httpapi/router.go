package httpapi

import "github.com/gin-gonic/gin"

// Handlers groups the HTTP handlers wired into the router.
type Handlers struct {
	Health *HealthHandler
}

// NewRouter registers the application routes.
func NewRouter(handlers Handlers) *gin.Engine {
	router := gin.Default()
	router.GET("/health", handlers.Health.Handle)

	return router
}
