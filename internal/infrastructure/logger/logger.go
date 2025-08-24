package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"beers-challenge/internal/core/ports/secondary"
)

// LogLevel represents the log level
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// StructuredLogger implements the secondary.Logger interface
type StructuredLogger struct {
	level  LogLevel
	format string
	logger *log.Logger
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(level, format string) secondary.Logger {
	logLevel := LevelInfo
	switch level {
	case "debug":
		logLevel = LevelDebug
	case "info":
		logLevel = LevelInfo
	case "warn":
		logLevel = LevelWarn
	case "error":
		logLevel = LevelError
	}

	return &StructuredLogger{
		level:  logLevel,
		format: format,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Info logs an info message
func (l *StructuredLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.shouldLog(LevelInfo) {
		l.log(LevelInfo, msg, fields, nil)
	}
}

// Error logs an error message
func (l *StructuredLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	if l.shouldLog(LevelError) {
		l.log(LevelError, msg, fields, err)
	}
}

// Debug logs a debug message
func (l *StructuredLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.shouldLog(LevelDebug) {
		l.log(LevelDebug, msg, fields, nil)
	}
}

// Warn logs a warning message
func (l *StructuredLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	if l.shouldLog(LevelWarn) {
		l.log(LevelWarn, msg, fields, nil)
	}
}

// shouldLog determines if a message should be logged based on level
func (l *StructuredLogger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}

	return levels[level] >= levels[l.level]
}

// log writes the log entry
func (l *StructuredLogger) log(level LogLevel, msg string, fields map[string]interface{}, err error) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     string(level),
		Message:   msg,
		Fields:    fields,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	var output string
	if l.format == "json" {
		jsonBytes, jsonErr := json.Marshal(entry)
		if jsonErr != nil {
			output = fmt.Sprintf(`{"timestamp":"%s","level":"error","message":"failed to marshal log entry","error":"%s"}`,
				time.Now().UTC().Format(time.RFC3339), jsonErr.Error())
		} else {
			output = string(jsonBytes)
		}
	} else {
		// Plain text format
		output = fmt.Sprintf("[%s] %s: %s", entry.Timestamp, entry.Level, entry.Message)
		if len(entry.Fields) > 0 {
			fieldsJSON, _ := json.Marshal(entry.Fields)
			output += fmt.Sprintf(" fields=%s", string(fieldsJSON))
		}
		if entry.Error != "" {
			output += fmt.Sprintf(" error=%s", entry.Error)
		}
	}

	l.logger.Println(output)
}

// NoOpLogger is a logger that does nothing (useful for testing)
type NoOpLogger struct{}

// NewNoOpLogger creates a new no-op logger
func NewNoOpLogger() secondary.Logger {
	return &NoOpLogger{}
}

// Info does nothing - this is intentional for testing purposes
func (l *NoOpLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	// No-op implementation for testing
}

// Error does nothing - this is intentional for testing purposes
func (l *NoOpLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	// No-op implementation for testing
}

// Debug does nothing - this is intentional for testing purposes
func (l *NoOpLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	// No-op implementation for testing
}

// Warn does nothing - this is intentional for testing purposes
func (l *NoOpLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	// No-op implementation for testing
}
