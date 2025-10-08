package database

import (
	"fmt"
	"os"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SharedDBConfig holds database configuration that can be used by both main connection and Flyway
type SharedDBConfig struct {
	Provider        string
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	ConnectionURL   string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	SupabasePooler  string
}

var (
	// Global configuration instance (singleton)
	globalConfig *SharedDBConfig
	configOnce   sync.Once
)

// LoadSharedDBConfig loads database configuration from environment variables
// This is the single source of truth for database configuration
// Uses singleton pattern to ensure configuration is loaded only once
func LoadSharedDBConfig() *SharedDBConfig {
	configOnce.Do(func() {
		logger.Info("Loading database configuration (first time only)...")
		globalConfig = loadDBConfigFromEnv()
	})
	return globalConfig
}

// loadDBConfigFromEnv loads database configuration from environment variables
// This is the actual implementation that does the work
func loadDBConfigFromEnv() *SharedDBConfig {
	provider := getEnv(constants.EnvDBProvider, constants.DefaultDBProvider)
	connectionURL := getEnv(constants.EnvDBConnectionURL, "")

	// Default SSL mode depends on provider
	defaultSSLMode := constants.DefaultSSLMode
	if provider == constants.DBProviderSupabase {
		defaultSSLMode = constants.SupabaseSSLMode
	}

	config := &SharedDBConfig{
		Provider:        provider,
		ConnectionURL:   connectionURL,
		MaxOpenConns:    getEnvAsInt(constants.EnvDBMaxOpenConns, constants.SupabaseDefaultMaxOpenConns),
		MaxIdleConns:    getEnvAsInt(constants.EnvDBMaxIdleConns, constants.SupabaseDefaultMaxIdleConns),
		ConnMaxLifetime: time.Duration(getEnvAsInt(constants.EnvDBConnMaxLifetime, constants.SupabaseDefaultConnMaxLife)) * time.Second,
		ConnMaxIdleTime: time.Duration(getEnvAsInt(constants.EnvDBConnMaxIdleTime, constants.SupabaseDefaultConnMaxIdle)) * time.Second,
		SupabasePooler:  getEnv(constants.EnvSupabasePooler, constants.SupabaseDefaultPoolerMode),
	}

	// Only load individual parameters if DATABASE_URL is not provided
	if connectionURL == "" {
		config.Host = getEnv(constants.EnvDBHost, constants.DefaultDBHost)
		config.Port = getEnv(constants.EnvDBPort, constants.DefaultDBPort)
		config.User = getEnv(constants.EnvDBUser, constants.DefaultDBUser)
		config.Password = getEnv(constants.EnvDBPassword, constants.DefaultDBPassword)
		config.DBName = getEnv(constants.EnvDBName, constants.DefaultDBName)
		config.SSLMode = getEnv(constants.EnvDBSSLMode, defaultSSLMode)
		logger.Info("Using individual database parameters (no DATABASE_URL provided)")
	} else {
		// Set defaults for individual parameters when using DATABASE_URL (for logging purposes)
		config.Host = "from_database_url"
		config.Port = "from_database_url"
		config.User = "from_database_url"
		config.Password = "from_database_url"
		config.DBName = "from_database_url"
		config.SSLMode = "from_database_url"
		logger.Info("Using DATABASE_URL connection string (individual parameters ignored)")
	}

	// Log configuration details
	logSharedDBConfig(config)

	return config
}

// GetSharedDBConfig returns the cached database configuration
// Returns nil if configuration hasn't been loaded yet
func GetSharedDBConfig() *SharedDBConfig {
	return globalConfig
}

// ResetSharedDBConfig resets the cached configuration (useful for testing)
func ResetSharedDBConfig() {
	configOnce = sync.Once{}
	globalConfig = nil
}

// logSharedDBConfig logs the shared database configuration (with sensitive data masked)
func logSharedDBConfig(config *SharedDBConfig) {
	logger.Info("Shared database configuration loaded",
		"provider", config.Provider,
		"host", config.Host,
		"port", config.Port,
		"database", config.DBName,
		"user", config.User,
		"password_set", config.Password != "",
		"sslmode", config.SSLMode,
		"connection_url_set", config.ConnectionURL != "",
		"supabase_pooler", config.SupabasePooler,
		"connection_pool", map[string]interface{}{
			"max_open_conns":     config.MaxOpenConns,
			"max_idle_conns":     config.MaxIdleConns,
			"conn_max_lifetime":  config.ConnMaxLifetime.String(),
			"conn_max_idle_time": config.ConnMaxIdleTime.String(),
		},
	)

	// Log environment variables status
	logger.Info("Environment variables status",
		"DB_PROVIDER", getEnvStatus(constants.EnvDBProvider),
		"DATABASE_URL", getEnvStatus(constants.EnvDBConnectionURL),
		"DB_HOST", getEnvStatus(constants.EnvDBHost),
		"DB_PORT", getEnvStatus(constants.EnvDBPort),
		"DB_USER", getEnvStatus(constants.EnvDBUser),
		"DB_PASSWORD", getEnvStatus(constants.EnvDBPassword),
		"DB_NAME", getEnvStatus(constants.EnvDBName),
		"DB_SSL_MODE", getEnvStatus(constants.EnvDBSSLMode),
	)
}

// BuildDSN builds a PostgreSQL connection string from shared config
func (config *SharedDBConfig) BuildDSN() string {
	// Use full connection URL if provided (typically for Supabase)
	if config.ConnectionURL != "" {
		return config.ConnectionURL
	}

	// Build DSN from individual components
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Add Supabase-specific parameters if using Supabase
	if config.Provider == constants.DBProviderSupabase {
		// Add pooler mode parameter for Supabase
		// This tells Supabase which pooler mode to use (transaction, session, or statement)
		if config.SupabasePooler != "" {
			dsn += " options='-c search_path=public'"
		}
	}

	return dsn
}

// GetMaskedDSN returns the DSN with sensitive information masked for logging
func (config *SharedDBConfig) GetMaskedDSN() string {
	dsn := config.BuildDSN()
	return maskPassword(dsn)
}

// HasSSLMode checks if the DSN contains SSL mode parameter
func (config *SharedDBConfig) HasSSLMode() bool {
	return containsParam(config.BuildDSN(), "sslmode")
}

// HasTimeout checks if the DSN contains timeout parameter
func (config *SharedDBConfig) HasTimeout() bool {
	return containsParam(config.BuildDSN(), "connect_timeout")
}

// HasPooler checks if the DSN contains pooler parameters
func (config *SharedDBConfig) HasPooler() bool {
	dsn := config.BuildDSN()
	return containsParam(dsn, "pool_")
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a fallback default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return intValue
		}
		logger.Warn("Invalid integer value for environment variable, using default",
			"key", key,
			"value", value,
			"default", defaultValue,
		)
	}
	return defaultValue
}

// maskPassword masks sensitive information in connection strings
func maskPassword(dsn string) string {
	// Simple password masking - replace password part with ***
	if strings.Contains(dsn, ":") && strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		if len(parts) >= 2 {
			userPart := parts[0]
			if strings.Contains(userPart, ":") {
				userParts := strings.Split(userPart, ":")
				if len(userParts) >= 3 {
					// postgresql://user:password@host
					return userParts[0] + ":***@" + strings.Join(parts[1:], "@")
				}
			}
		}
	}
	return dsn
}

// containsParam checks if a URL contains a specific parameter
func containsParam(url, param string) bool {
	return strings.Contains(url, param+"=")
}

// getEnvStatus returns the status of an environment variable
func getEnvStatus(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return "not_set"
	}
	if key == constants.EnvDBPassword || key == constants.EnvDBConnectionURL {
		return "set_masked"
	}
	return "set"
}
