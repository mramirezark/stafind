package database

import (
	"database/sql"
	"fmt"
	"os"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/logger"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// DBConfig holds database configuration
type DBConfig struct {
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

// NewConnection creates a new database connection
func NewConnection() (*DB, error) {
	logger.Info("Initializing database connection...")

	config := loadDBConfig()
	logDBConfig(config)

	var dsn string

	// Use full connection URL if provided (typically for Supabase)
	if config.ConnectionURL != "" {
		dsn = config.ConnectionURL
		logger.Info("Using DATABASE_URL connection string",
			"provider", config.Provider,
			"connectionURL", maskPassword(dsn),
			"hasSSLMode", containsParam(dsn, "sslmode"),
			"hasTimeout", containsParam(dsn, "connect_timeout"),
			"hasPooler", containsParam(dsn, "pool_"),
		)
	} else {
		// Build DSN from individual components
		dsn = buildDSN(config)
		logger.Info("Built DSN from individual components",
			"provider", config.Provider,
			"host", config.Host,
			"port", config.Port,
			"database", config.DBName,
			"user", config.User,
			"sslmode", config.SSLMode,
			"supabasePooler", config.SupabasePooler,
		)
	}

	logger.Info("Attempting to open database connection...",
		"dsn_masked", maskPassword(dsn),
		"connection_pool_config", map[string]interface{}{
			"max_open_conns":     config.MaxOpenConns,
			"max_idle_conns":     config.MaxIdleConns,
			"conn_max_lifetime":  config.ConnMaxLifetime.String(),
			"conn_max_idle_time": config.ConnMaxIdleTime.String(),
		},
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open database connection", "error", err.Error())
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	logger.Info("Configuring connection pool...")
	configureConnectionPool(db, config)

	logger.Info("Testing database connection with ping...")
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database",
			"error", err.Error(),
			"provider", config.Provider,
			"host", config.Host,
			"port", config.Port,
		)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("✅ Database connection established successfully",
		"provider", config.Provider,
		"host", config.Host,
		"port", config.Port,
		"database", config.DBName,
		"connection_pool", map[string]interface{}{
			"max_open_conns": config.MaxOpenConns,
			"max_idle_conns": config.MaxIdleConns,
		},
	)

	return &DB{db}, nil
}

// loadDBConfig loads database configuration from environment
func loadDBConfig() *DBConfig {
	provider := getEnv(constants.EnvDBProvider, constants.DefaultDBProvider)

	// Default SSL mode depends on provider
	defaultSSLMode := constants.DefaultSSLMode
	if provider == constants.DBProviderSupabase {
		defaultSSLMode = constants.SupabaseSSLMode
	}

	config := &DBConfig{
		Provider:        provider,
		Host:            getEnv(constants.EnvDBHost, constants.DefaultDBHost),
		Port:            getEnv(constants.EnvDBPort, constants.DefaultDBPort),
		User:            getEnv(constants.EnvDBUser, constants.DefaultDBUser),
		Password:        getEnv(constants.EnvDBPassword, constants.DefaultDBPassword),
		DBName:          getEnv(constants.EnvDBName, constants.DefaultDBName),
		SSLMode:         getEnv(constants.EnvDBSSLMode, defaultSSLMode),
		ConnectionURL:   getEnv(constants.EnvDBConnectionURL, ""),
		MaxOpenConns:    getEnvAsInt(constants.EnvDBMaxOpenConns, constants.SupabaseDefaultMaxOpenConns),
		MaxIdleConns:    getEnvAsInt(constants.EnvDBMaxIdleConns, constants.SupabaseDefaultMaxIdleConns),
		ConnMaxLifetime: time.Duration(getEnvAsInt(constants.EnvDBConnMaxLifetime, constants.SupabaseDefaultConnMaxLife)) * time.Second,
		ConnMaxIdleTime: time.Duration(getEnvAsInt(constants.EnvDBConnMaxIdleTime, constants.SupabaseDefaultConnMaxIdle)) * time.Second,
		SupabasePooler:  getEnv(constants.EnvSupabasePooler, constants.SupabaseDefaultPoolerMode),
	}

	return config
}

// buildDSN builds a PostgreSQL connection string from config
func buildDSN(config *DBConfig) string {
	// Base DSN
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

// configureConnectionPool sets up connection pool parameters
func configureConnectionPool(db *sql.DB, config *DBConfig) {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
}

// RunMigrations runs database migrations
func (db *DB) RunMigrations() error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
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

// logDBConfig logs the database configuration (with sensitive data masked)
func logDBConfig(config *DBConfig) {
	logger.Info("Database configuration loaded",
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
