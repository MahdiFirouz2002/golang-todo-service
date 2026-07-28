package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler exposes liveness and readiness probes.
type HealthHandler struct {
	startedAt time.Time
	ping      func(ctx context.Context) error
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler(ping func(ctx context.Context) error) *HealthHandler {
	return &HealthHandler{
		startedAt: time.Now().UTC(),
		ping:      ping,
	}
}

// Live responds when the process is running.
// GET /health/live
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// Ready responds when the service is ready to accept traffic.
// GET /health/ready
func (h *HealthHandler) Ready(c *gin.Context) {
	if h.ping != nil {
		if err := h.ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  "database unavailable",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"uptime": time.Since(h.startedAt).Round(time.Second).String(),
	})
}
