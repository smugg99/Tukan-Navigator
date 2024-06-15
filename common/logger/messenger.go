package logger

import (
	"fmt"
	"regexp"
)

type MessageWrapper struct {
	Key      string    `json:"key"`
	Message  string    `json:"message"`
	LogLevel LogLevel  `json:"log_level"`
}

func removeFormattingCharacters(message string) string {
	re := regexp.MustCompile(`%[a-zA-Z]`)
	return re.ReplaceAllString(message, "")
}

func NewMessageWrapper(key string, message string, logLevel LogLevel) *MessageWrapper {
	return &MessageWrapper{
		Key:      key,
		Message:  message,
		LogLevel: logLevel,
	}
}

func (e *MessageWrapper) String() string {
	return removeFormattingCharacters(e.Message)
}

func (e *MessageWrapper) Msg() string {
	return e.String()
}

func (e *MessageWrapper) Error() string {
	return e.Msg()
}

func (e *MessageWrapper) FormatMsg(vars ...interface{}) string {
	if vars == nil {
		return e.Msg()
    }

	return fmt.Sprintf(e.Message, vars...)
}

func (e *MessageWrapper) Format(vars ...interface{}) *MessageWrapper {
	copy := *e
	copy.Message = e.FormatMsg(vars...)
	return &copy
}

func (e *MessageWrapper) IsError() bool {
	if e.LogLevel == ErrorLevel || e.LogLevel == FatalLevel{
		return true
	}

	return false
}