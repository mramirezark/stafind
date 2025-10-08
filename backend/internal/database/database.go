package database

import (
	"database/sql"
	"fmt"
	"stafind-backend/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// DBConfig holds database configuration (now uses SharedDBConfig)
type DBConfig struct {
	*SharedDBConfig
}

// NewConnection creates a new database connection
func NewConnection() (*DB, error) {
	logger.Info("Initializing database connection...")

	// Use shared configuration
	sharedConfig := LoadSharedDBConfig()
	config := &DBConfig{SharedDBConfig: sharedConfig}

	dsn := config.BuildDSN()

	logger.Info("Using shared database configuration",
		"provider", config.Provider,
		"connectionURL", config.GetMaskedDSN(),
		"hasSSLMode", config.HasSSLMode(),
		"hasTimeout", config.HasTimeout(),
		"hasPooler", config.HasPooler(),
	)

	logger.Info("Attempting to open database connection...",
		"dsn_masked", config.GetMaskedDSN(),
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

// loadDBConfig is deprecated - use LoadSharedDBConfig() instead
func loadDBConfig() *DBConfig {
	sharedConfig := LoadSharedDBConfig()
	return &DBConfig{SharedDBConfig: sharedConfig}
}

// buildDSN is deprecated - use config.BuildDSN() instead
func buildDSN(config *DBConfig) string {
	return config.BuildDSN()
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

// Helper functions moved to config.go for shared use
