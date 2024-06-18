package routes

import (
	"smuggr.xyz/tukan-navigator/api/v1/handlers"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func SetupGraphsRoutes(router *gin.Engine, rootGroup *gin.RouterGroup) {
	graphsLimiter := tollbooth.NewLimiter(1.5, nil)

	graphGroup := rootGroup.Group("/graph")
	graphGroup.Use(tollbooth_gin.LimitHandler(graphsLimiter))
	{
		graphGroup.GET("/:hash", handlers.RetrieveGraphHandler)
		graphGroup.POST("/store", handlers.StoreGraphHandler)
		graphGroup.POST("/shortest-path", handlers.GetShortestPathHandler)
	}
}
