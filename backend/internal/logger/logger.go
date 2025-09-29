package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"stafind-backend/internal/constants"
	"strings"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Config holds the logger configuration
type Config struct {
	Level  LogLevel
	Format string // "json" or "text"
	Output string // "stdout", "stderr", or file path
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:  LevelInfo,
		Format: "json",
		Output: "stdout",
	}
}

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	config *Config
}

var (
	// Global logger instance
	global *Logger
)

// Init initializes the global logger with the given configuration
func Init(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Set up output
	var output io.Writer
	switch strings.ToLower(config.Output) {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constants.DefaultFilePermission)
		if err != nil {
			return err
		}
		output = file
	}

	// Set up level
	var level slog.Level
	switch strings.ToLower(string(config.Level)) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Set up handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	// Create handler based on format
	var handler slog.Handler
	switch strings.ToLower(config.Format) {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	case "text":
		handler = slog.NewTextHandler(output, opts)
	default:
		handler = slog.NewJSONHandler(output, opts)
	}

	// Create logger
	global = &Logger{
		Logger: slog.New(handler),
		config: config,
	}

	return nil
}

// Get returns the global logger instance
func Get() *Logger {
	if global == nil {
		// Initialize with default config if not already initialized
		Init(DefaultConfig())
	}
	return global
}

// WithContext returns a logger with the given context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Logger: l.Logger,
		config: l.config,
	}
}

// WithFields returns a logger with the given fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(args...),
		config: l.config,
	}
}

// WithComponent returns a logger with a component field
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.Logger.With("component", component),
		config: l.config,
	}
}

// Fatal logs a fatal error and exits the program
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Error(msg, args...)
	os.Exit(1)
}

// Fatalf logs a formatted fatal error and exits the program
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Errorf(format, args...)
	os.Exit(1)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(format, args...)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(format, args...)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(format, args...)
}

// Convenience functions for global logger
func Debug(msg string, args ...interface{}) {
	Get().Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	Get().Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Get().Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	Get().Fatal(msg, args...)
}

func Fatalf(format string, args ...interface{}) {
	Get().Fatalf(format, args...)
}

func WithFields(fields map[string]interface{}) *Logger {
	return Get().WithFields(fields)
}

func WithComponent(component string) *Logger {
	return Get().WithComponent(component)
}
