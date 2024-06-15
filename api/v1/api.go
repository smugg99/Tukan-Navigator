package v1

import (
	"os"
	"strconv"

	"smuggr.xyz/tukan-navigator/api/v1/routes"
	"smuggr.xyz/tukan-navigator/common/configurator"
	"smuggr.xyz/tukan-navigator/common/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Config = &configurator.Config.API
var Logger = logger.NewCustomLogger("api/v1")

var DefaultRouter *gin.Engine

func SetupCors() {
	DefaultRouter.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			Logger.Debug(origin)
			return true
		},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	}))

	DefaultRouter.SetTrustedProxies([]string{})
}

func Initialize() (*logger.MessageWrapper, chan *logger.MessageWrapper) {
	Logger.Log(logger.MsgInitializing)

	Config = &configurator.Config.API
	gin.SetMode(os.Getenv("GIN_MODE"))

	DefaultRouter = gin.Default()
	SetupCors()

	routes.Initialize(DefaultRouter)

	errCh := make(chan *logger.MessageWrapper)
	go func() {
		err := DefaultRouter.Run(":" + strconv.Itoa(int(Config.Port)))
		errCh <- logger.ErrUnexpected.Format(err.Error())
	}()

	return logger.MsgInitialized, errCh
}

func Cleanup() *logger.MessageWrapper {
	return nil
}
