package main

import (
	"os"
	"os/signal"
	"syscall"

	"smuggr.xyz/tukan-navigator/common/configurator"
	"smuggr.xyz/tukan-navigator/common/logger"
	"smuggr.xyz/tukan-navigator/core/datastorer"

	api "smuggr.xyz/tukan-navigator/api/v1"
)

var Logger = logger.DefaultLogger

func WaitForTermination() {
	callChan := make(chan os.Signal, 1)
	signal.Notify(callChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	Logger.Debug("waiting for termination signal...")
	<-callChan
	Logger.Debug("termination signal received")
}

func Initialize() {
	logger.Initialize()
	configurator.Initialize()

	datastorer.Logger.Log(datastorer.Initialize())

	apiMsg, apiChan := api.Initialize()
	api.Logger.Log(apiMsg)

	if msg := <-apiChan; msg != nil {
		api.Logger.Log(msg)
	}
}

func Cleanup() {
	if err := recover(); err != nil {
		Logger.Warn("panic", err)
	}

	Logger.Log(logger.MsgCleaningUp)

	if err := api.Cleanup(); err != nil {
		api.Logger.Error(err.Error())
	}

	if err := datastorer.Cleanup(); err != nil {
		datastorer.Logger.Error(err.Error())
	}
}

func main() {
	Initialize()

	defer Cleanup()
	defer WaitForTermination()
}
