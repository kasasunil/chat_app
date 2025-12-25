package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger is a singleton logger instance
type Logger struct {
	level  LogLevel
	logger *log.Logger
	mu     sync.RWMutex
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			level:  INFO,
			logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		}
	})
	return instance
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(file *os.File) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetOutput(file)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level <= WARN {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.RLock()
	level := l.level
	l.mu.RUnlock()

	if level <= ERROR {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logger.Fatalf("[FATAL] "+format, v...)
}

// Panic logs a panic message and panics
func (l *Logger) Panic(format string, v ...interface{}) {
	l.logger.Panicf("[PANIC] "+format, v...)
}

// ParseLogLevel parses a string log level to LogLevel
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug", "DEBUG":
		return DEBUG
	case "info", "INFO":
		return INFO
	case "warn", "WARN", "warning", "WARNING":
		return WARN
	case "error", "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "INFO"
	}
}

// Convenience functions for global access
func Debug(format string, v ...interface{}) {
	GetLogger().Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	GetLogger().Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	GetLogger().Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	GetLogger().Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	GetLogger().Fatal(format, v...)
}

func Panic(format string, v ...interface{}) {
	GetLogger().Panic(format, v...)
}

// InitLogger initializes the logger with configuration
func InitLogger(level string, outputFile string) error {
	logger := GetLogger()

	// Set log level
	logger.SetLevel(ParseLogLevel(level))

	// Set output file if provided
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logger.SetOutput(file)
	}

	return nil
}
