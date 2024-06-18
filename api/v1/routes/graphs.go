package routes

import (
	"smuggr.xyz/tukan-navigator/api/v1/handlers"

	"github.com/gin-gonic/gin"
)

func SetupGraphsRoutes(router *gin.Engine, rootGroup *gin.RouterGroup) {
	graphGroup := rootGroup.Group("/graph")
	{
		graphGroup.GET("/:hash", handlers.RetrieveGraphHandler)
		graphGroup.POST("/store", handlers.StoreGraphHandler)
		graphGroup.POST("/shortest-path", handlers.GetShortestPathHandler)
	}
}
