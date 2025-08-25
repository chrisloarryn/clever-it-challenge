package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testErrorMessage = "error message"
	testError        = "test error"
	testMessage      = "test message"
)

func TestNewStructuredLogger(t *testing.T) {
	logger := NewStructuredLogger("debug", "json")
	assert.NotNil(t, logger)
}

func TestLevels(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	logger := NewStructuredLogger("info", "json").(*StructuredLogger)
	logger.logger.SetOutput(&buf)

	ctx := context.Background()

	// Should not be logged
	logger.Debug(ctx, "debug message", nil)
	assert.Equal(t, "", buf.String())

	// Should be logged
	logger.Info(ctx, "info message", nil)
	assert.True(t, strings.Contains(buf.String(), "info message"))
	buf.Reset()

	logger.Warn(ctx, "warn message", nil)
	assert.True(t, strings.Contains(buf.String(), "warn message"))
	buf.Reset()

	logger.Error(ctx, testErrorMessage, errors.New(testError), nil)
	assert.True(t, strings.Contains(buf.String(), testErrorMessage))
	assert.True(t, strings.Contains(buf.String(), testError))
	buf.Reset()
}

func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	logger := NewStructuredLogger("debug", "json").(*StructuredLogger)
	logger.logger.SetOutput(&buf)

	ctx := context.Background()
	fields := map[string]interface{}{"key": "value"}
	err := errors.New(testError)

	logger.Info(ctx, testMessage, fields)

	var entry LogEntry
	json.Unmarshal(buf.Bytes(), &entry)

	assert.Equal(t, "info", entry.Level)
	assert.Equal(t, testMessage, entry.Message)
	assert.Equal(t, "value", entry.Fields["key"])

	buf.Reset()

	logger.Error(ctx, testErrorMessage, err, fields)
	json.Unmarshal(buf.Bytes(), &entry)

	assert.Equal(t, "error", entry.Level)
	assert.Equal(t, testErrorMessage, entry.Message)
	assert.Equal(t, err.Error(), entry.Error)
}

func TestTextFormat(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	logger := NewStructuredLogger("debug", "text").(*StructuredLogger)
	logger.logger.SetOutput(&buf)

	ctx := context.Background()
	fields := map[string]interface{}{"key": "value"}
	err := errors.New(testError)

	logger.Info(ctx, testMessage, fields)
	logOutput := buf.String()

	assert.Contains(t, logOutput, "[info]")
	assert.Contains(t, logOutput, testMessage)
	assert.Contains(t, logOutput, `"key":"value"`)

	buf.Reset()

	logger.Error(ctx, testErrorMessage, err, fields)
	logOutput = buf.String()

	assert.Contains(t, logOutput, "[error]")
	assert.Contains(t, logOutput, testErrorMessage)
	assert.Contains(t, logOutput, "error="+testError)
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()
	// These should not panic
	logger.Info(context.Background(), "test", nil)
	logger.Error(context.Background(), "test", nil, nil)
	logger.Debug(context.Background(), "test", nil)
	logger.Warn(context.Background(), "test", nil)
}
