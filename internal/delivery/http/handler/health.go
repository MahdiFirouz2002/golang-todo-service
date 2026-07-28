package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler exposes liveness and readiness probes.
type HealthHandler struct {
	startedAt time.Time
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startedAt: time.Now().UTC()}
}

// Live responds when the process is running.
// GET /health/live
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// Ready responds when the service is ready to accept traffic.
// Persistence checks will be wired in a later stage.
// GET /health/ready
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"uptime":  time.Since(h.startedAt).Round(time.Second).String(),
	})
}
