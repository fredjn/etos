package logging

import (
	"io"
	"log"
	"os"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// LogLevelError represents error level logging
	LogLevelError LogLevel = iota
	// LogLevelWarning represents warning level logging
	LogLevelWarning
	// LogLevelInfo represents info level logging
	LogLevelInfo
	// LogLevelDebug represents debug level logging
	LogLevelDebug
)

var (
	// DefaultLogger is the default logger instance
	DefaultLogger *Logger
)

// Logger provides configurable logging
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger creates a new logger with the specified level and output
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level >= LogLevelError {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Warning logs a warning message
func (l *Logger) Warning(format string, v ...interface{}) {
	if l.level >= LogLevelWarning {
		l.logger.Printf("[WARNING] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level >= LogLevelInfo {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level >= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func init() {
	// Initialize the default logger
	DefaultLogger = NewLogger(LogLevelInfo, os.Stdout)
} 