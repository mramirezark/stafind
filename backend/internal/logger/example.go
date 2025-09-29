package logger

import (
	"context"
	"fmt"
	"time"
)

// Example demonstrates various logging capabilities
func Example() {
	// Initialize logger
	Init(&Config{
		Level:  LevelDebug,
		Format: "json",
		Output: "stdout",
	})

	log := Get()

	// Basic logging
	log.Info("Application started")
	log.Debug("Debug information", "user_id", 123, "action", "login")
	log.Warn("Warning message", "retry_count", 3)
	log.Error("Error occurred", "error", fmt.Errorf("database connection failed"))

	// Structured logging with fields
	log.WithFields(map[string]interface{}{
		"user_id":    123,
		"session_id": "abc-123",
		"ip_address": "192.168.1.1",
	}).Info("User logged in")

	// Component-specific logging
	dbLogger := log.WithComponent("database")
	dbLogger.Info("Database connection established", "host", "localhost", "port", 5432)
	dbLogger.Error("Query failed", "query", "SELECT * FROM users", "error", "timeout")

	// Context-aware logging
	ctx := context.WithValue(context.Background(), "request_id", "req-456")
	ctxLogger := log.WithContext(ctx)
	ctxLogger.Info("Processing request", "endpoint", "/api/users")

	// Performance logging
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	log.Info("Operation completed", "duration_ms", time.Since(start).Milliseconds())

	// Fatal logging (will exit the program)
	// log.Fatal("Critical error", "error", "system failure")
}

// ExampleLogbackStyle demonstrates Logback-style logging patterns
func ExampleLogbackStyle() {
	log := Get()

	// Similar to Logback's MDC (Mapped Diagnostic Context)
	log.WithFields(map[string]interface{}{
		"trace_id": "trace-789",
		"span_id":  "span-456",
	}).Info("Processing request")

	// Similar to Logback's structured logging
	log.Info("User action",
		"user_id", 123,
		"action", "create_user",
		"resource", "users",
		"timestamp", time.Now().Unix(),
	)

	// Error with stack trace context
	log.Error("Database operation failed",
		"operation", "insert",
		"table", "users",
		"error", "constraint violation",
		"retry_attempt", 1,
	)
}
