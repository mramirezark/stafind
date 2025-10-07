package main

import (
	"database/sql"
	"fmt"
	"os"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/database"
	"stafind-backend/internal/logger"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger
	if err := logger.Init(nil); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log := logger.Get()

	// Load environment variables
	// Try .env first (standard), then fall back to config.env (legacy)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("config.env"); err != nil {
			log.Warn("No .env or config.env file found, using environment variables")
		} else {
			log.Info("Loaded environment from config.env (consider renaming to .env)")
		}
	} else {
		log.Info("Loaded environment from .env")
	}

	// Get provider
	provider := os.Getenv(constants.EnvDBProvider)
	if provider == "" {
		provider = constants.DefaultDBProvider
	}

	log.Info("Database Clean Utility", "provider", provider)

	// Confirm action
	fmt.Println("⚠️  WARNING: This will DELETE ALL DATA from your database!")
	fmt.Printf("Provider: %s\n", provider)
	fmt.Print("Are you sure you want to continue? (type 'yes' to confirm): ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "yes" {
		fmt.Println("❌ Cancelled.")
		os.Exit(0)
	}

	// Connect to database
	log.Info("Connecting to database...")
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	log.Info("Connected successfully")

	// Execute cleanup
	log.Info("Dropping public schema...")
	if err := executeSQL(db.DB, "DROP SCHEMA public CASCADE;"); err != nil {
		log.Fatal("Failed to drop schema", "error", err)
	}

	log.Info("Creating public schema...")
	if err := executeSQL(db.DB, "CREATE SCHEMA public;"); err != nil {
		log.Fatal("Failed to create schema", "error", err)
	}

	log.Info("Creating extensions...")
	if err := executeSQL(db.DB, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		log.Warn("Failed to create uuid-ossp extension (may already exist)", "error", err)
	}

	log.Info("Database cleaned successfully!")

	// Run migrations
	log.Info("Running migrations...")
	flywayConfig := database.NewFlywayConfigFromEnv()
	flywayMigrator, err := database.NewFlywayMigrator(flywayConfig)
	if err != nil {
		log.Fatal("Failed to initialize Flyway migrator", "error", err)
	}
	defer flywayMigrator.Close()

	// Validate migrations
	if err := flywayMigrator.ValidateMigrations(); err != nil {
		log.Fatal("Migration validation failed", "error", err)
	}

	// Run migrations
	if err := flywayMigrator.Migrate(); err != nil {
		log.Fatal("Failed to run migrations", "error", err)
	}

	log.Info("✅ Database cleaned and migrated successfully!")
}

func executeSQL(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}
