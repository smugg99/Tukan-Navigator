package routes

import (
	"os"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Initialize(defaultRouter *gin.Engine) {
	defaultLimiter := tollbooth.NewLimiter(1, nil)

	defaultRouter.Use(static.Serve("/", static.LocalFile(os.Getenv("DIST_PATH"), false)))

	rootGroup := defaultRouter.Group("/api/v1")
	rootGroup.Use(tollbooth_gin.LimitHandler(defaultLimiter))

	SetupGraphsRoutes(defaultRouter, rootGroup)

	defaultRouter.NoRoute(func(c *gin.Context) {
		c.File(os.Getenv("DIST_PATH") + "/index.html")
	})
}
