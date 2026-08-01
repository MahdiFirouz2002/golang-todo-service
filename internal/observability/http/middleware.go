package obshttp

import (
	"github.com/MahdiFirouz2002/golang-todo-service/internal/observability/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Register wires tracing and Prometheus middleware into the router.
func Register(router *gin.Engine, serviceName string) {
	router.Use(otelgin.Middleware(serviceName))
	router.Use(metrics.Middleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
