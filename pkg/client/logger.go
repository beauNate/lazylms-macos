package client

import "fmt"

// LogLevel represents the severity of a log message
type LogLevel string

const (
	LogLevelError LogLevel = "ERROR"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelDebug LogLevel = "DEBUG"
)

// Logger handles logging to a channel
type Logger struct {
	Channel chan string
}

func NewLogger(ch chan string) *Logger {
	return &Logger{Channel: ch}
}

// logToChannel sends a formatted log message to the logger's channel
func (l *Logger) logToChannel(level LogLevel, message string, args ...interface{}) {
	if l.Channel == nil {
		return
	}

	var prefix string
	if level == LogLevelInfo {
		prefix = "INFO:"
	} else {
		prefix = fmt.Sprintf("[%s]", level)
	}

	formattedMsg := fmt.Sprintf("%s "+message, append([]interface{}{prefix}, args...)...)

	select {
	case l.Channel <- formattedMsg:
	default:
		// Channel is full or closed, drop the message to prevent blocking
	}
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	l.logToChannel(LogLevelError, message, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logToChannel(LogLevelWarn, message, args...)
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	l.logToChannel(LogLevelInfo, message, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logToChannel(LogLevelDebug, message, args...)
}
