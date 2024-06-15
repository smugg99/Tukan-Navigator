package routes

import (
	"smuggr.xyz/tukan-navigator/api/v1/handlers"
	// "smuggr.xyz/tukan-navigator/api/v1/middleware"

	"github.com/gin-gonic/gin"
)

func SetupGraphsRoutes(rootGroup *gin.RouterGroup) {
	graphsGroup := rootGroup.Group("/graphs")
	{
		graphsGroup.POST("/shortest-path", handlers.GetShortestPathHandler)
	}
}
