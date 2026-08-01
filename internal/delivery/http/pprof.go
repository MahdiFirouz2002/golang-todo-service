package httpserver

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func registerProfiling(router *gin.Engine) {
	pprof.Register(router, "/debug/pprof")
}
