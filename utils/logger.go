package utils

import (
	"fmt"
	"time"
)

// Logger 日志
type Logger struct {
	logType string
}

// Message is a function to return debug message.
func (ll *Logger) Message(msg string) {
	fmt.Printf("[%s] %s %s\n", ll.logType, time.Now().Format("2006-01-02 15:04:05.00000 -0700"), msg)
}

// Log 返回日志对象
func Log(logType string) *Logger {

	return &Logger{logType: logType}
}
