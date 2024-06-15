package handlers

import (
	"net/http"

	"smuggr.xyz/tukan-navigator/common/configurator"
	"smuggr.xyz/tukan-navigator/common/logger"
	"smuggr.xyz/tukan-navigator/core/grapher"

	"github.com/gin-gonic/gin"
)

var APIConfig = &configurator.Config.API

func bindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.BindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": logger.ErrInvalidRequestPayload})
		return false
	}
	return true
}

func GetShortestPathHandler(c *gin.Context) {
	var graph grapher.Graph
	if !bindJSON(c, &graph) {
		return
	}

	if path, distance := grapher.GetShortestPath(graph); path != nil{
		c.JSON(http.StatusOK, gin.H{"path": path, "distance": distance})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": logger.ErrNoPathFound})
}
