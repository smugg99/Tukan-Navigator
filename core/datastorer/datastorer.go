package datastorer

import (
	"os"

	"smuggr.xyz/jungle-path-finder/common/configurator"
	"smuggr.xyz/jungle-path-finder/common/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Config *configurator.DatastorerConfig
var Logger = logger.NewCustomLogger("data")

var DB *gorm.DB

func GetDSN() string {
	return "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable TimeZone=UTC"
}

func Initialize() *logger.MessageWrapper {
	Logger.Log(logger.MsgInitializing)

	Config = &configurator.Config.Datastorer

	db, err := gorm.Open(postgres.Open(GetDSN()))
	if err != nil {
		return logger.ErrInitializing.Format(err.Error())
	}

	DB = db

	if err := DB.AutoMigrate(&User{}, &Device{}); err != nil {
		return logger.ErrInitializing.Format(err.Error())
	}

	return logger.MsgInitialized
}

func Cleanup() *logger.MessageWrapper {
	Logger.Info(logger.MsgCleaningUp)

	sqlDB, err := DB.DB()
	if err != nil {
		return logger.ErrCleaningUp.Format(err.Error())
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			return logger.ErrCleaningUp.Format(err.Error())
		}
	}

	return logger.MsgCleanedUp
}
