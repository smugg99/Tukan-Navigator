package routes

import (
	"net/http"
	"os"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func Initialize(defaultRouter *gin.Engine) {
	defaultLimiter := tollbooth.NewLimiter(1, nil)

    defaultRouter.StaticFS("/", http.Dir(os.Getenv("DIST_PATH")))

	rootGroup := defaultRouter.Group("/api/v1")
	rootGroup.Use(tollbooth_gin.LimitHandler(defaultLimiter))

	SetupGraphsRoutes(rootGroup)
}