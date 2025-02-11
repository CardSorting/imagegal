package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	// DEBUG level for detailed information
	DEBUG Level = iota
	// INFO level for general operational information
	INFO
	// ERROR level for errors that need attention
	ERROR
)

// Logger represents the application logger
type Logger struct {
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
}

// New creates a new Logger instance
func New() *Logger {
	flags := log.Ldate | log.Ltime | log.LUTC

	return &Logger{
		debug: log.New(os.Stdout, "DEBUG: ", flags),
		info:  log.New(os.Stdout, "INFO: ", flags),
		error: log.New(os.Stderr, "ERROR: ", flags),
	}
}

// formatMessage formats a log message with fields
func formatMessage(msg string, fields ...interface{}) string {
	if len(fields) == 0 {
		return msg
	}

	// Add timestamp
	message := fmt.Sprintf("[%s] %s", time.Now().UTC().Format(time.RFC3339), msg)

	// Add fields
	if len(fields) > 0 {
		message += " |"
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				message += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}

	return message
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.debug.Println(formatMessage(msg, fields...))
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.info.Println(formatMessage(msg, fields...))
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	// Add error to fields
	errorFields := append([]interface{}{"error", err.Error()}, fields...)
	l.error.Println(formatMessage(msg, errorFields...))
}

// WithFields creates a new log entry with fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return l
}

// WithError adds an error to the log entry
func (l *Logger) WithError(err error) *Logger {
	return l
}

// Default logger instance
var defaultLogger = New()

// GetDefault returns the default logger instance
func GetDefault() *Logger {
	return defaultLogger
}
