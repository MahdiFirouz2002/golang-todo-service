package httpserver

import (
	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

// Dependencies holds the HTTP-layer collaborators required to build the router.
type Dependencies struct {
	Health *handler.HealthHandler
	Tasks  *handler.TaskHandler
}

// NewRouter constructs the Gin engine and registers routes.
func NewRouter(deps Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	registerRoutes(router, deps)

	return router
}

func registerRoutes(router *gin.Engine, deps Dependencies) {
	health := router.Group("/health")
	{
		health.GET("/live", deps.Health.Live)
		health.GET("/ready", deps.Health.Ready)
	}

	if deps.Tasks != nil {
		tasks := router.Group("/api/v1/tasks")
		{
			tasks.POST("", deps.Tasks.Create)
			tasks.GET("", deps.Tasks.List)
			tasks.GET("/:id", deps.Tasks.Get)
			tasks.PUT("/:id", deps.Tasks.Update)
			tasks.DELETE("/:id", deps.Tasks.Delete)
		}
	}
}
