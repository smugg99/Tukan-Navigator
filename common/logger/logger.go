package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type CustomLogger struct {
	*log.Logger
}

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

func NewCustomLogger(packageName string) *CustomLogger {
	logger := log.New(os.Stdout)

	logger.SetLevel(log.DebugLevel)
	logger.SetStyles(GetPackageLoggerStyle(packageName))

	return &CustomLogger{Logger: logger}
}

var DefaultLogger = NewCustomLogger("main")

func (c *CustomLogger) Log(message *MessageWrapper) {
	switch message.LogLevel {
	case DebugLevel:
		c.Debug(message.Message)
	case InfoLevel:
		c.Info(message.Message)
	case WarnLevel:
		c.Warn(message.Message)
	case ErrorLevel:
		c.Error(message.Message)
	case FatalLevel:
		c.Fatal(message.Message)
	default:
		c.Info(message.Message)
	}
}

func Initialize() {
	DefaultLogger.Log(MsgInitializing)

	log.SetLevel(log.DebugLevel)
	log.SetStyles(GetDefaultLoggerStyle())
}
