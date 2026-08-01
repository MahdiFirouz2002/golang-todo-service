package httpserver

import (
	"io/fs"
	"net/http"

	appdocs "github.com/MahdiFirouz2002/golang-todo-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerSwagger(router *gin.Engine) {
	specFS, err := fs.Sub(appdocs.OpenAPI, ".")
	if err != nil {
		panic(err)
	}

	router.StaticFS("/openapi", http.FS(specFS))
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.FileFromFS("openapi.yaml", http.FS(specFS))
	})

	url := ginSwagger.URL("/openapi.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
