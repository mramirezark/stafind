package main

import (
	"fmt"
	"log"
	"os"
	"stafind-backend/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	config := database.NewFlywayConfigFromEnv()
	migrator, err := database.NewFlywayMigrator(config)
	if err != nil {
		log.Fatalf("Failed to initialize Flyway migrator: %v", err)
	}
	defer migrator.Close()

	command := os.Args[1]

	switch command {
	case "migrate":
		if err := migrator.Migrate(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✅ Migrations completed successfully!")

	case "info":
		info, err := migrator.GetMigrationInfo()
		if err != nil {
			log.Fatalf("Failed to get migration info: %v", err)
		}

		fmt.Println("📊 Migration Information")
		fmt.Println("========================")
		fmt.Printf("Applied migrations: %d\n", len(info.AppliedMigrations))
		fmt.Printf("Pending migrations: %d\n", len(info.PendingMigrations))

		if len(info.AppliedMigrations) > 0 {
			fmt.Println("\n✅ Applied Migrations:")
			for _, migration := range info.AppliedMigrations {
				status := "✅"
				if !migration.Success {
					status = "❌"
				}
				fmt.Printf("  %s %s - %s (%s)\n", status, migration.Version, migration.Description, migration.InstalledOn)
			}
		}

		if len(info.PendingMigrations) > 0 {
			fmt.Println("\n⏳ Pending Migrations:")
			for _, migration := range info.PendingMigrations {
				fmt.Printf("  ⏳ %s\n", migration)
			}
		} else {
			fmt.Println("\n✅ No pending migrations")
		}

	case "validate":
		if err := migrator.ValidateMigrations(); err != nil {
			log.Fatalf("Migration validation failed: %v", err)
		}
		fmt.Println("✅ All migrations are valid!")

	case "pending":
		info, err := migrator.GetMigrationInfo()
		if err != nil {
			log.Fatalf("Failed to get migration info: %v", err)
		}

		if len(info.PendingMigrations) == 0 {
			fmt.Println("✅ No pending migrations")
		} else {
			fmt.Printf("⏳ %d pending migrations:\n", len(info.PendingMigrations))
			for _, migration := range info.PendingMigrations {
				fmt.Printf("  - %s\n", migration)
			}
		}

	case "applied":
		info, err := migrator.GetMigrationInfo()
		if err != nil {
			log.Fatalf("Failed to get migration info: %v", err)
		}

		if len(info.AppliedMigrations) == 0 {
			fmt.Println("📝 No applied migrations")
		} else {
			fmt.Printf("✅ %d applied migrations:\n", len(info.AppliedMigrations))
			for _, migration := range info.AppliedMigrations {
				status := "✅"
				if !migration.Success {
					status = "❌"
				}
				fmt.Printf("  %s %s - %s\n", status, migration.Version, migration.Description)
			}
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Flyway CLI for StaffFind")
	fmt.Println("========================")
	fmt.Println("Usage: go run cmd/flyway-cli/main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  migrate    - Run all pending migrations")
	fmt.Println("  info       - Show detailed migration information")
	fmt.Println("  validate   - Validate all migration files")
	fmt.Println("  pending    - List pending migrations")
	fmt.Println("  applied    - List applied migrations")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DB_HOST         - Database host (default: localhost)")
	fmt.Println("  DB_PORT         - Database port (default: 5432)")
	fmt.Println("  DB_USER         - Database user (default: postgres)")
	fmt.Println("  DB_PASSWORD     - Database password (default: password)")
	fmt.Println("  DB_NAME         - Database name (default: stafind)")
	fmt.Println("  DB_SSLMODE      - SSL mode (default: disable)")
	fmt.Println("  FLYWAY_LOCATIONS - Migration directory (default: ./flyway_migrations)")
}
