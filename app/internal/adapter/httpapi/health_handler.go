package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler serves the service health check.
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
