package configurator

import (
	"os"

	"smuggr.xyz/tukan-navigator/common/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Logger = logger.NewCustomLogger("conf")
var Config GlobalConfig

func loadEnv() *logger.MessageWrapper {
	err := godotenv.Load()
	if err != nil {
		return logger.ErrReadingResource.Format("", logger.ResourceEnv)
	}

	return logger.MsgResourceLoaded.Format(".env", logger.ResourceEnv)
}

func loadConfig(config *GlobalConfig) *logger.MessageWrapper {
	if err := viper.ReadInConfig(); err != nil {
		return logger.ErrReadingResource.Format("", logger.ResourceConfig)
	}

	err := viper.Unmarshal(config)
	if err != nil {
		return logger.ErrFormattingResource.Format(logger.ResourceConfig)
	}

	return logger.MsgResourceLoaded.Format(viper.ConfigFileUsed(), logger.ResourceConfig)
}

func Initialize() {
	Logger.Info(logger.MsgInitializing)
	Logger.Log(loadEnv())

	viper.AddConfigPath(os.Getenv("CONFIG_PATH"))
	viper.SetConfigType(os.Getenv("CONFIG_TYPE"))

	Logger.Log(loadConfig(&Config))
}
