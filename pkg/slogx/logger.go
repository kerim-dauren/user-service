package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	DefaultHandler = "text"   // Default log handler format
	DefaultWriter  = "stdout" // Default log writer
	DefaultLevel   = "DEBUG"  // Default log level
)

// Config holds the logger configuration.
type Config struct {
	Handler string // "text" or "json"
	Level   string // "DEBUG", "INFO", "WARN", "ERROR"
	Writer  string // "stdout", "stderr", or a file path
}

// NewLogger configures and sets the default logger.
// If cfg is nil, default values are used.
func NewLogger(cfg *Config) (*slog.Logger, error) {
	if cfg == nil {
		cfg = &Config{
			Handler: DefaultHandler,
			Level:   DefaultLevel,
			Writer:  DefaultWriter,
		}
	}

	// Set default values if fields are empty.
	if cfg.Handler == "" {
		cfg.Handler = DefaultHandler
	}
	if cfg.Level == "" {
		cfg.Level = DefaultLevel
	}
	if cfg.Writer == "" {
		cfg.Writer = DefaultWriter
	}

	handler, err := getHandler(cfg)
	if err != nil {
		return nil, err
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger, nil
}

// getHandler returns a log handler based on the configuration settings.
func getHandler(cfg *Config) (slog.Handler, error) {
	levelVar := &slog.LevelVar{}
	switch strings.ToUpper(cfg.Level) {
	case "INFO":
		levelVar.Set(slog.LevelInfo)
	case "WARN":
		levelVar.Set(slog.LevelWarn)
	case "ERROR":
		levelVar.Set(slog.LevelError)
	case "DEBUG":
		levelVar.Set(slog.LevelDebug)
	default:
		return nil, fmt.Errorf("unrecognized log level: %s", cfg.Level)
	}

	opts := &slog.HandlerOptions{Level: levelVar}

	// Select the writer based on configuration.
	var w io.Writer
	switch strings.ToLower(cfg.Writer) {
	case "stdout":
		w = os.Stdout
	case "stderr":
		w = os.Stderr
	default:
		// If a file path is provided, attempt to open the file for writing.
		file, err := os.OpenFile(cfg.Writer, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", cfg.Writer, err)
		}
		w = file
		// Note: Ensure to close the file properly when the application shuts down.
	}

	// Select the log handler format.
	switch strings.ToLower(cfg.Handler) {
	case "text":
		return slog.NewTextHandler(w, opts), nil
	case "json":
		return slog.NewJSONHandler(w, opts), nil
	default:
		return nil, fmt.Errorf("unexpected log format: %s", cfg.Handler)
	}
}
