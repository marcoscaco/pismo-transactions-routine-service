package httpapi

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const serviceName = "pismo-transactions-routine-service"

// Handlers groups the HTTP handlers wired into the router.
type Handlers struct {
	Health *HealthHandler
}

// NewRouter registers the application routes. A nil metrics value disables
// HTTP metrics collection.
func NewRouter(handlers Handlers, metrics *Metrics) *gin.Engine {
	router := gin.New()

	// registrar os middlewares
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(serviceName))
	if metrics != nil {
		router.Use(metrics.Middleware())
	}

	// registrar as rotas
	router.GET("/health", handlers.Health.Handle)

	return router
}
