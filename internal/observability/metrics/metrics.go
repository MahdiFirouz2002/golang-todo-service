package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of HTTP requests processed",
	}, []string{"method", "path", "status"})

	requestLatencyHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_latency_histogram",
		Help:    "HTTP request latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	tasksCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_count",
		Help: "Current number of tasks stored in the database",
	})
)

// Middleware records request count and latency metrics.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		status := strconv.Itoa(c.Writer.Status())
		requestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		requestLatencyHistogram.WithLabelValues(c.Request.Method, path).Observe(time.Since(start).Seconds())
	}
}

// SetTasksCount updates the tasks_count gauge.
func SetTasksCount(count float64) {
	tasksCount.Set(count)
}
