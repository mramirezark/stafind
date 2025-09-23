# Flyway-Style Migrations for StaffFind

This project uses **goflyway**, a Go library that provides Flyway-like database migration capabilities for PostgreSQL.

## 🚀 Features

- **Flyway-Compatible**: Follows standard Flyway naming conventions and behavior
- **Version Control**: Tracks applied migrations in `schema_version` table
- **Validation**: Validates migration files before execution
- **CLI Tools**: Built-in command-line tools for migration management
- **Rollback Support**: Can handle rollback scenarios
- **Sample Data**: Includes sample engineers, skills, and job requests

## 📁 Migration Structure

```
flyway_migrations/
├── V1__Create_initial_schema.sql      # Database schema
├── V2__Insert_sample_skills.sql       # Sample skills data
├── V3__Add_sample_engineers.sql       # Sample engineers and skills
└── V4__Add_sample_job_requests.sql    # Sample job requests
```

### Naming Convention

Migrations follow the Flyway naming pattern:
- **Format**: `V<version>__<description>.sql`
- **Version**: Sequential number (V1, V2, V3, ...)
- **Separator**: Double underscore (`__`)
- **Description**: Human-readable description
- **Extension**: `.sql`

## 🛠️ CLI Tools

### Migration Management

```bash
# Run all pending migrations
go run cmd/flyway-cli/main.go migrate

# Show detailed migration information
go run cmd/flyway-cli/main.go info

# Validate migration files
go run cmd/flyway-cli/main.go validate

# List pending migrations
go run cmd/flyway-cli/main.go pending

# List applied migrations
go run cmd/flyway-cli/main.go applied
```

### Migration Testing

```bash
# Test migration file validation
go run cmd/flyway-test/main.go
```

## ⚙️ Configuration

### Environment Variables

```bash
# Database connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=stafind
DB_SSLMODE=disable

# Migration settings
FLYWAY_LOCATIONS=./flyway_migrations
```

### Configuration File

The `flyway.conf` file provides additional configuration options:

```properties
# Database connection
flyway.url=jdbc:postgresql://localhost:5432/stafind
flyway.user=postgres
flyway.password=password

# Migration settings
flyway.schemas=public
flyway.locations=filesystem:./flyway_migrations
flyway.baselineOnMigrate=true
flyway.validateOnMigrate=true

# Naming conventions
flyway.sqlMigrationPrefix=V
flyway.sqlMigrationSeparator=__
flyway.sqlMigrationSuffixes=.sql
```

## 📊 Migration History

### V1: Initial Schema
- Creates all core tables (engineers, skills, job_requests, matches)
- Sets up foreign key relationships
- Creates performance indexes
- Establishes database constraints

### V2: Sample Skills
- Inserts 20 common technical skills
- Covers programming languages, frameworks, and tools
- Provides foundation for engineer skill matching

### V3: Sample Engineers
- Creates 5 sample engineers with different skill levels
- Associates engineers with relevant skills
- Includes proficiency levels and experience years
- Covers various departments and locations

### V4: Sample Job Requests
- Creates 5 sample job requests
- Demonstrates different skill requirements
- Shows various priority levels and departments
- Includes realistic job descriptions

## 🔧 Usage in Code

### Automatic Migration on Startup

The server automatically runs migrations on startup:

```go
// In cmd/server/main.go
flywayConfig := database.NewFlywayConfigFromEnv()
flywayMigrator, err := database.NewFlywayMigrator(flywayConfig)
if err != nil {
    log.Fatal("Failed to initialize Flyway migrator:", err)
}

// Validate migrations first
if err := flywayMigrator.ValidateMigrations(); err != nil {
    log.Fatal("Migration validation failed:", err)
}

// Run migrations
if err := flywayMigrator.Migrate(); err != nil {
    log.Fatal("Failed to run Flyway migrations:", err)
}
```

### Manual Migration Management

```go
// Create migrator
config := database.NewFlywayConfigFromEnv()
migrator, err := database.NewFlywayMigrator(config)
if err != nil {
    log.Fatal(err)
}
defer migrator.Close()

// Get migration info
info, err := migrator.GetMigrationInfo()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Applied: %d, Pending: %d\n", 
    len(info.AppliedMigrations), 
    len(info.PendingMigrations))

// Run migrations
if err := migrator.Migrate(); err != nil {
    log.Fatal(err)
}
```

## 📋 Creating New Migrations

### 1. Create Migration File

Create a new file in `flyway_migrations/` following the naming convention:

```bash
# Example: V5__Add_user_authentication.sql
touch flyway_migrations/V5__Add_user_authentication.sql
```

### 2. Write Migration SQL

```sql
-- V5__Add_user_authentication.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

### 3. Validate Migration

```bash
go run cmd/flyway-cli/main.go validate
```

### 4. Apply Migration

```bash
go run cmd/flyway-cli/main.go migrate
```

## 🔍 Migration Status

### Check Applied Migrations

```bash
go run cmd/flyway-cli/main.go applied
```

Output:
```
✅ 4 applied migrations:
  ✅ V1 - Create initial schema
  ✅ V2 - Insert sample skills
  ✅ V3 - Add sample engineers
  ✅ V4 - Add sample job requests
```

### Check Pending Migrations

```bash
go run cmd/flyway-cli/main.go pending
```

### Detailed Information

```bash
go run cmd/flyway-cli/main.go info
```

## 🛡️ Best Practices

### 1. Migration Design
- **Idempotent**: Migrations should be safe to run multiple times
- **Atomic**: Each migration should be a complete, logical change
- **Sequential**: Use sequential version numbers
- **Descriptive**: Use clear, descriptive names

### 2. File Organization
- Keep migrations in the `flyway_migrations/` directory
- Use consistent naming conventions
- Group related changes in single migrations
- Separate schema changes from data changes when possible

### 3. Testing
- Validate migrations before applying
- Test migrations on development databases first
- Use the CLI tools to check status
- Verify data integrity after migrations

### 4. Rollback Strategy
- Design migrations with rollback in mind
- Keep backup of important data
- Test rollback procedures
- Document rollback steps

## 🔧 Troubleshooting

### Common Issues

1. **Connection Refused**
   ```
   Error: dial tcp [::1]:5432: connect: connection refused
   ```
   **Solution**: Ensure PostgreSQL is running and accessible

2. **Invalid Filename Format**
   ```
   Error: Invalid migration filename: V1_Description.sql
   ```
   **Solution**: Use double underscore: `V1__Description.sql`

3. **Migration Already Applied**
   ```
   Error: Migration V1 already applied
   ```
   **Solution**: Check migration history and use next version number

### Debugging

```bash
# Validate all migration files
go run cmd/flyway-cli/main.go validate

# Check migration status
go run cmd/flyway-cli/main.go info

# Test filename validation
go run cmd/flyway-test/main.go
```

## 📚 References

- [goflyway GitHub](https://github.com/gabrielaraujosouza/goflyway)
- [Flyway Documentation](https://flywaydb.org/documentation/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🎯 Migration Checklist

- [ ] Migration file follows naming convention (`V<version>__<description>.sql`)
- [ ] SQL syntax is valid PostgreSQL
- [ ] Migration is idempotent
- [ ] File is not empty
- [ ] Changes are logical and atomic
- [ ] Indexes are created for performance
- [ ] Constraints are properly defined
- [ ] Migration has been tested locally
- [ ] Documentation is updated if needed
