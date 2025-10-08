package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/logger"

	"github.com/gabrielaraujosouza/goflyway"
	_ "github.com/lib/pq"
)

// FlywayMigrator handles database migrations using goflyway
type FlywayMigrator struct {
	db       *sql.DB
	config   *FlywayConfig
	location string
}

// FlywayConfig holds configuration for Flyway migrations
type FlywayConfig struct {
	Provider      string
	Host          string
	Port          string
	User          string
	Password      string
	Database      string
	SSLMode       string
	ConnectionURL string
	Location      string
}

// NewFlywayMigrator creates a new Flyway migrator
func NewFlywayMigrator(config *FlywayConfig) (*FlywayMigrator, error) {
	logger.Info("Initializing Flyway migrator...")

	var dsn string

	// Use full connection URL if provided (typically for Supabase)
	if config.ConnectionURL != "" {
		dsn = config.ConnectionURL
		logger.Info("Using DATABASE_URL for Flyway migrations",
			"provider", config.Provider,
			"connectionURL", maskPassword(dsn),
			"hasSSLMode", containsParam(dsn, "sslmode"),
			"hasTimeout", containsParam(dsn, "connect_timeout"),
		)
	} else {
		// Build connection string from individual components
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)
		logger.Info("Using individual parameters for Flyway migrations",
			"provider", config.Provider,
			"host", config.Host,
			"port", config.Port,
			"database", config.Database,
			"user", config.User,
			"ssl_mode", config.SSLMode,
			"dsn_masked", maskPassword(dsn),
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool for migrations (conservative settings)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Flyway migrator connected to database successfully")

	// Set default location if not provided
	location := config.Location
	if location == "" {
		location = "./flyway_migrations"
	}

	// Check if migration directory exists
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return nil, fmt.Errorf("migration directory does not exist: %s", location)
	}

	return &FlywayMigrator{
		db:       db,
		config:   config,
		location: location,
	}, nil
}

// Migrate runs all pending migrations
func (fm *FlywayMigrator) Migrate() error {
	logger.Info("Starting Flyway migrations")

	// Configure goflyway
	conf := goflyway.GoFlywayConfig{
		Db:       fm.db,
		Driver:   goflyway.POSTGRES,
		Location: fm.location,
	}

	// Run migrations
	totalScriptsExecuted, err := goflyway.Migrate(conf)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("Flyway migrations completed successfully", "total_scripts", totalScriptsExecuted)
	return nil
}

// GetMigrationInfo returns information about migration status
func (fm *FlywayMigrator) GetMigrationInfo() (*MigrationInfo, error) {
	// Get migration history from goflyway's schema_version table
	query := `
		SELECT version, description, type, script, checksum, installed_by, installed_on, execution_time, success
		FROM schema_version 
		ORDER BY installed_rank
	`

	rows, err := fm.db.Query(query)
	if err != nil {
		// If schema_version table doesn't exist, return empty info
		if err.Error() == "relation \"schema_version\" does not exist" {
			return &MigrationInfo{
				AppliedMigrations: []AppliedMigration{},
				PendingMigrations: fm.getPendingMigrations(),
			}, nil
		}
		return nil, fmt.Errorf("failed to query migration history: %w", err)
	}
	defer rows.Close()

	var appliedMigrations []AppliedMigration
	for rows.Next() {
		var migration AppliedMigration
		err := rows.Scan(
			&migration.Version,
			&migration.Description,
			&migration.Type,
			&migration.Script,
			&migration.Checksum,
			&migration.InstalledBy,
			&migration.InstalledOn,
			&migration.ExecutionTime,
			&migration.Success,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		appliedMigrations = append(appliedMigrations, migration)
	}

	return &MigrationInfo{
		AppliedMigrations: appliedMigrations,
		PendingMigrations: fm.getPendingMigrations(),
	}, nil
}

// getPendingMigrations returns a list of pending migration files
func (fm *FlywayMigrator) getPendingMigrations() []string {
	var pending []string

	// Read migration directory
	files, err := filepath.Glob(filepath.Join(fm.location, "V*.sql"))
	if err != nil {
		logger.Error("Error reading migration directory", "error", err)
		return pending
	}

	// Get applied migrations to determine which are pending
	info, err := fm.GetMigrationInfo()
	if err != nil {
		logger.Error("Error getting migration info", "error", err)
		return pending
	}

	appliedVersions := make(map[string]bool)
	for _, migration := range info.AppliedMigrations {
		appliedVersions[migration.Version] = true
	}

	// Check each file to see if it's pending
	for _, file := range files {
		filename := filepath.Base(file)
		// Extract version from filename (e.g., V1__Description.sql -> V1)
		if len(filename) > 0 && filename[0] == 'V' {
			versionEnd := 1
			for i := 1; i < len(filename); i++ {
				if filename[i] == '_' || filename[i] == '.' {
					break
				}
				versionEnd++
			}
			version := filename[:versionEnd]

			if !appliedVersions[version] {
				pending = append(pending, filename)
			}
		}
	}

	return pending
}

// ValidateMigrations validates that all migration files are properly formatted
func (fm *FlywayMigrator) ValidateMigrations() error {
	files, err := filepath.Glob(filepath.Join(fm.location, "V*.sql"))
	if err != nil {
		return fmt.Errorf("error reading migration directory: %w", err)
	}

	var errors []string
	for _, file := range files {
		filename := filepath.Base(file)

		// Validate filename format: V<version>__<description>.sql
		if !fm.IsValidMigrationFilename(filename) {
			errors = append(errors, fmt.Sprintf("Invalid migration filename: %s (should be V<version>__<description>.sql)", filename))
		}

		// Validate file content is not empty
		content, err := os.ReadFile(file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error reading migration file %s: %v", filename, err))
			continue
		}

		if len(content) == 0 {
			errors = append(errors, fmt.Sprintf("Migration file %s is empty", filename))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("migration validation failed:\n%s", fmt.Sprintf("- %s\n", errors))
	}

	logger.Info("Migration validation passed", "valid_files", len(files))
	return nil
}

// IsValidMigrationFilename validates the migration filename format
func (fm *FlywayMigrator) IsValidMigrationFilename(filename string) bool {
	// Should start with V followed by version number, then __, then description, then .sql
	if len(filename) < 6 { // Minimum: V1__.sql
		return false
	}

	if filename[0] != 'V' {
		return false
	}

	// Find the first underscore
	firstUnderscore := -1
	for i := 1; i < len(filename); i++ {
		if filename[i] == '_' {
			firstUnderscore = i
			break
		}
		if filename[i] < '0' || filename[i] > '9' {
			return false // Version should only contain numbers
		}
	}

	if firstUnderscore == -1 || firstUnderscore >= len(filename)-5 {
		return false
	}

	// Check for double underscore
	if firstUnderscore+1 >= len(filename) || filename[firstUnderscore+1] != '_' {
		return false
	}

	// Check for .sql extension
	if filepath.Ext(filename) != ".sql" {
		return false
	}

	return true
}

// Close closes the database connection
func (fm *FlywayMigrator) Close() error {
	return fm.db.Close()
}

// MigrationInfo holds information about migration status
type MigrationInfo struct {
	AppliedMigrations []AppliedMigration
	PendingMigrations []string
}

// AppliedMigration represents a successfully applied migration
type AppliedMigration struct {
	Version       string
	Description   string
	Type          string
	Script        string
	Checksum      int64
	InstalledBy   string
	InstalledOn   string
	ExecutionTime int64
	Success       bool
}

// NewFlywayConfigFromEnv creates a FlywayConfig from environment variables
func NewFlywayConfigFromEnv() *FlywayConfig {
	provider := getEnv(constants.EnvDBProvider, constants.DefaultDBProvider)

	// Default SSL mode depends on provider
	defaultSSLMode := constants.DefaultSSLMode
	if provider == constants.DBProviderSupabase {
		defaultSSLMode = constants.SupabaseSSLMode
	}

	return &FlywayConfig{
		Provider:      provider,
		Host:          getEnv(constants.EnvDBHost, constants.DefaultDBHost),
		Port:          getEnv(constants.EnvDBPort, constants.DefaultDBPort),
		User:          getEnv(constants.EnvDBUser, constants.DefaultDBUser),
		Password:      getEnv(constants.EnvDBPassword, constants.DefaultDBPassword),
		Database:      getEnv(constants.EnvDBName, constants.DefaultDBName),
		SSLMode:       getEnv(constants.EnvDBSSLMode, defaultSSLMode),
		ConnectionURL: getEnv(constants.EnvDBConnectionURL, ""),
		Location:      getEnv(constants.EnvFlywayLocations, constants.DefaultFlywayLocations),
	}
}
