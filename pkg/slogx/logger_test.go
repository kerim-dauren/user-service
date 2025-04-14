package slogx

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		logger, err := NewLogger(nil)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("CustomConfig", func(t *testing.T) {
		cfg := &Config{
			Handler: "json",
			Level:   "INFO",
			Writer:  "stderr",
		}
		logger, err := NewLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("InvalidLevel", func(t *testing.T) {
		cfg := &Config{
			Handler: "text",
			Level:   "INVALID",
			Writer:  "stdout",
		}
		logger, err := NewLogger(cfg)
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("InvalidHandler", func(t *testing.T) {
		cfg := &Config{
			Handler: "invalid",
			Level:   "DEBUG",
			Writer:  "stdout",
		}
		logger, err := NewLogger(cfg)
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("InvalidWriter", func(t *testing.T) {
		cfg := &Config{
			Handler: "text",
			Level:   "DEBUG",
			Writer:  "/invalid/path/to/logfile.log",
		}
		logger, err := NewLogger(cfg)
		assert.Error(t, err)
		assert.Nil(t, logger)
	})
}

func TestGetHandler(t *testing.T) {
	t.Run("ValidTextHandler", func(t *testing.T) {
		cfg := &Config{
			Handler: "text",
			Level:   "INFO",
			Writer:  "stdout",
		}
		handler, err := getHandler(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, handler)
		assert.IsType(t, &slog.TextHandler{}, handler)
	})

	t.Run("ValidJSONHandler", func(t *testing.T) {
		cfg := &Config{
			Handler: "json",
			Level:   "DEBUG",
			Writer:  "stderr",
		}
		handler, err := getHandler(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, handler)
		assert.IsType(t, &slog.JSONHandler{}, handler)
	})

	t.Run("InvalidLogLevel", func(t *testing.T) {
		cfg := &Config{
			Handler: "text",
			Level:   "INVALID",
			Writer:  "stdout",
		}
		handler, err := getHandler(cfg)
		assert.Error(t, err)
		assert.Nil(t, handler)
	})

	t.Run("InvalidWriterPath", func(t *testing.T) {
		cfg := &Config{
			Handler: "text",
			Level:   "DEBUG",
			Writer:  "/invalid/path/to/logfile.log",
		}
		handler, err := getHandler(cfg)
		assert.Error(t, err)
		assert.Nil(t, handler)
	})

	t.Run("InvalidHandlerFormat", func(t *testing.T) {
		cfg := &Config{
			Handler: "invalid",
			Level:   "DEBUG",
			Writer:  "stdout",
		}
		handler, err := getHandler(cfg)
		assert.Error(t, err)
		assert.Nil(t, handler)
	})
}

func TestFileWriter(t *testing.T) {
	t.Run("ValidFileWriter", func(t *testing.T) {
		cfg := &Config{
			Handler: "text",
			Level:   "DEBUG",
			Writer:  "test.log",
		}
		handler, err := getHandler(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, handler)

		// Clean up
		_ = os.Remove("test.log")
	})
}
