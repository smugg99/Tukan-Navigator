package handlers

import (
	"net/http"

	"smuggr.xyz/tukan-navigator/common/configurator"
	"smuggr.xyz/tukan-navigator/common/logger"
	"smuggr.xyz/tukan-navigator/core/grapher"
	"smuggr.xyz/tukan-navigator/core/datastorer"

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

func StoreGraphHandler(c *gin.Context) {
	var graph grapher.GraphProject
	if !bindJSON(c, &graph) {
		return
	}

	// TO-DO: Validate graph
	if sortedGraph, graphHash, err := datastorer.SortAndHashGraph(graph); err == nil {
		Logger.Debug("sortedGraph", sortedGraph, graphHash)
		if err := datastorer.StoreGraph(sortedGraph, graphHash); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"hash": graphHash,
				"error": logger.ErrResourceAlreadyExists.Format(graphHash, logger.ResourceGraph),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"hash": graphHash})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": logger.ErrHashingResource.Format("", logger.ResourceGraph)})
		return
	}
}

func RetrieveGraphHandler(c *gin.Context) {
	hash := c.Param("hash")

	if graph, err := datastorer.GetGraphByHash(hash); err == nil {
		c.JSON(http.StatusOK, graph)
		c.Redirect(http.StatusFound, "/?graph="+hash)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": logger.ErrResourceNotFound.Format(hash, logger.ResourceGraph)})
	}
}