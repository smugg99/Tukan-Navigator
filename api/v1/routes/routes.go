package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
)

func Initialize(defaultRouter *gin.Engine) {
	defaultLimiter := tollbooth.NewLimiter(1, nil)

	rootGroup := defaultRouter.Group("/api/v1")
	rootGroup.Use(tollbooth_gin.LimitHandler(defaultLimiter))

	SetupGraphsRoutes(rootGroup)
}