# Structured Logging with slog

This package provides a structured logging solution using Go's `slog` package, similar to Logback in Java.

## Features

- **Structured JSON logging** with timestamps and levels
- **Context-aware logging** with request tracing
- **Component-specific loggers** for different parts of the application
- **Configurable log levels** (DEBUG, INFO, WARN, ERROR)
- **Multiple output formats** (JSON, text)
- **Performance logging** with duration tracking
- **Fatal logging** with automatic program termination

## Quick Start

```go
import "stafind-backend/internal/logger"

// Initialize with default config
logger.Init(nil)

// Get the global logger
log := logger.Get()

// Basic logging
log.Info("Application started")
log.Error("Something went wrong", "error", err)

// Structured logging
log.WithFields(map[string]interface{}{
    "user_id": 123,
    "action": "login",
}).Info("User logged in")

// Component-specific logging
dbLogger := log.WithComponent("database")
dbLogger.Info("Connection established", "host", "localhost")
```

## Configuration

```go
config := &logger.Config{
    Level:  logger.LevelInfo,
    Format: "json",        // "json" or "text"
    Output: "stdout",      // "stdout", "stderr", or file path
}

logger.Init(config)
```

## Log Levels

- `DEBUG`: Detailed information for debugging
- `INFO`: General information about program execution
- `WARN`: Warning messages for potential issues
- `ERROR`: Error messages for failures

## Structured Fields

All log messages support structured fields:

```go
log.Info("User action",
    "user_id", 123,
    "action", "create_user",
    "resource", "users",
    "timestamp", time.Now().Unix(),
)
```

## Component Loggers

Create component-specific loggers for better organization:

```go
dbLogger := log.WithComponent("database")
apiLogger := log.WithComponent("api")
authLogger := log.WithComponent("auth")

dbLogger.Info("Query executed", "query", "SELECT * FROM users")
apiLogger.Info("Request processed", "endpoint", "/api/users")
authLogger.Warn("Invalid token", "token", "abc123")
```

## Context-Aware Logging

Use context for request tracing:

```go
ctx := context.WithValue(context.Background(), "request_id", "req-456")
ctxLogger := log.WithContext(ctx)
ctxLogger.Info("Processing request", "endpoint", "/api/users")
```

## Performance Logging

Track operation duration:

```go
start := time.Now()
// ... perform operation ...
log.Info("Operation completed", "duration_ms", time.Since(start).Milliseconds())
```

## Fatal Logging

For critical errors that should terminate the program:

```go
log.Fatal("Critical system failure", "error", "database unavailable")
// Program will exit with code 1
```

## Output Examples

### JSON Format
```json
{
  "time": "2025-09-28T23:41:58.521412712-07:00",
  "level": "INFO",
  "source": {
    "function": "main.main",
    "file": "/workspace/stafind/backend/cmd/server/main.go",
    "line": 27
  },
  "msg": "Server starting",
  "port": "8080"
}
```

### Text Format
```
2025/09/28 23:41:58 INFO Server starting port=8080
```

## Migration from Standard log Package

### Before (standard log)
```go
import "log"

log.Println("User logged in")
log.Printf("User %d logged in", userID)
log.Fatal("Database connection failed:", err)
```

### After (structured logging)
```go
import "stafind-backend/internal/logger"

log := logger.Get()
log.Info("User logged in", "user_id", userID)
log.Fatal("Database connection failed", "error", err)
```

## Best Practices

1. **Use structured fields** instead of string formatting
2. **Include relevant context** in log messages
3. **Use appropriate log levels** for different types of messages
4. **Create component loggers** for different parts of your application
5. **Include error details** when logging errors
6. **Use consistent field names** across your application

## Similar to Logback

This implementation provides similar functionality to Logback:

- **Structured logging** like Logback's JSON encoder
- **MDC-like functionality** with `WithFields()`
- **Component logging** like Logback's logger hierarchy
- **Configurable levels** and output formats
- **Context propagation** for request tracing
