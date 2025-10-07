# StaffFind Makefile for Docker Management

.PHONY: help dev prod stop logs migrate status cleanup build-backend build-frontend backend-air backend-dev backend-dev-simple frontend-dev frontend-build frontend-install full-dev quick-start quick-start-simple

# Default target
help: ## Show this help message
	@echo "StaffFind Docker Management"
	@echo "=========================="
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development environment
dev: ## Start development environment (PostgreSQL + pgAdmin)
	@echo "Starting development environment..."
	cd backend && docker compose -f docker-compose.dev.yml up -d
	@echo "Development environment ready!"
	@echo "Database: localhost:5432"
	@echo "pgAdmin: http://localhost:5050 (admin@stafind.com / admin)"

# Production environment
prod: ## Start production environment (full stack)
	@echo "Starting production environment..."
	docker compose up -d
	@echo "Production environment ready!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

# Stop all services
stop: ## Stop all Docker services
	@echo "Stopping all services..."
	docker compose down
	cd backend && docker compose -f docker-compose.dev.yml down
	@echo "All services stopped."

# Show logs
logs: ## Show logs for all services
	docker compose logs -f

logs-backend: ## Show backend logs
	docker compose logs -f backend

logs-db: ## Show database logs
	docker compose logs -f postgres

# Run migrations
migrate: ## Run database migrations
	@echo "Running database migrations..."
	cd backend && go run cmd/flyway-cli/main.go migrate

# Show status
status: ## Show status of all services
	@echo "Production services:"
	docker compose ps
	@echo ""
	@echo "Development services:"
	cd backend && docker compose -f docker-compose.dev.yml ps

# Build services
build-backend: ## Build backend Docker image
	@echo "Building backend image..."
	docker compose build backend

build-frontend: ## Build frontend Docker image
	@echo "Building frontend image..."
	docker compose build frontend

build-all: ## Build all Docker images
	@echo "Building all images..."
	docker compose build

# Clean up
cleanup: ## Remove all containers, volumes, and images
	@echo "This will remove all Docker resources. Are you sure? (y/N)"
	@read -r response && [ "$$response" = "y" ]
	docker compose down -v --remove-orphans
	cd backend && docker compose -f docker-compose.dev.yml down -v --remove-orphans
	docker system prune -f
	@echo "Cleanup completed!"

# Development helpers
dev-logs: ## Show development environment logs
	cd backend && docker compose -f docker-compose.dev.yml logs -f

dev-stop: ## Stop development environment
	cd backend && docker compose -f docker-compose.dev.yml down

# Backend development
backend-dev: ## Run backend in development mode (requires dev environment)
	@echo "Starting backend in development mode..."
	cd backend && go run cmd/server/main.go

backend-air: ## Run backend with Air for live reloading (requires dev environment)
	@echo "Starting backend with Air for live reloading..."
	cd backend && air -c .air.toml

backend-dev-simple: ## Run backend in development mode (simple, no Air)
	@echo "Starting backend in simple development mode..."
	cd backend && go run cmd/server/main.go

# Database management
db-shell: ## Connect to PostgreSQL shell
	docker compose exec postgres psql -U postgres -d stafind

db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "WARNING: This will delete all data. Are you sure? (y/N)"
	@read -r response && [ "$$response" = "y" ]
	docker compose down postgres
	docker volume rm stafind_postgres_data || true
	docker compose up -d postgres
	@echo "Database reset completed!"

db-clean: ## Clean all database objects (local PostgreSQL only)
	@echo "Cleaning all database objects..."
	cd backend && docker compose -f docker-compose.dev.yml exec postgres psql -U postgres -d stafind -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	cd backend && docker compose -f docker-compose.dev.yml exec postgres psql -U postgres -d stafind -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
	@echo "Running migrations..."
	cd backend && go run cmd/flyway-cli/main.go migrate
	@echo "Database cleaned and migrated successfully!"

db-clean-go: ## Clean database (PostgreSQL or Supabase - pure Go, no psql/Docker required) ⭐
	@echo "Running database clean utility (pure Go)..."
	cd backend && go run cmd/db-clean/main.go

# Database provider switching
db-use-postgres: ## Switch to local PostgreSQL
	@echo "Switching to local PostgreSQL..."
	@if [ -f backend/.env.local ]; then \
		cp backend/.env.local backend/.env; \
		echo "✅ Switched to local PostgreSQL (from .env.local)"; \
	elif [ -f backend/config.env.local ]; then \
		cp backend/config.env.local backend/.env; \
		echo "✅ Switched to local PostgreSQL (from config.env.local - legacy)"; \
	else \
		echo "⚠️  .env.local not found. Creating from template..."; \
		cp backend/config.env.example backend/.env; \
		echo "DB_PROVIDER=postgres" >> backend/.env; \
		echo "✅ Created .env with PostgreSQL settings"; \
	fi
	@echo "Don't forget to restart your backend!"

db-use-supabase: ## Switch to Supabase
	@echo "Switching to Supabase..."
	@if [ -f backend/.env.supabase ]; then \
		cp backend/.env.supabase backend/.env; \
		echo "✅ Switched to Supabase (from .env.supabase)"; \
	elif [ -f backend/config.env.supabase ]; then \
		cp backend/config.env.supabase backend/.env; \
		echo "✅ Switched to Supabase (from config.env.supabase - legacy)"; \
	elif [ -f backend/supabase.env.example ]; then \
		echo "⚠️  .env.supabase not found."; \
		echo "Please copy backend/supabase.env.example to backend/.env.supabase"; \
		echo "and update it with your Supabase credentials."; \
		exit 1; \
	else \
		echo "❌ supabase.env.example not found!"; \
		exit 1; \
	fi
	@echo "Don't forget to restart your backend!"

db-test-connection: ## Test database connection
	@echo "Testing database connection..."
	@cd backend && go run -exec 'env' cmd/server/main.go 2>&1 | grep -i "connected to database" || echo "❌ Failed to connect"

db-show-config: ## Show current database configuration
	@echo "Current Database Configuration:"
	@echo "==============================="
	@if [ -f backend/.env ]; then \
		echo "Source: backend/.env"; \
		echo "Provider: $$(grep '^DB_PROVIDER=' backend/.env | cut -d'=' -f2 || echo 'postgres (default)')"; \
		echo "Host: $$(grep '^DB_HOST=' backend/.env | cut -d'=' -f2 || echo 'not set')"; \
		echo "Port: $$(grep '^DB_PORT=' backend/.env | cut -d'=' -f2 || echo 'not set')"; \
		echo "Database: $$(grep '^DB_NAME=' backend/.env | cut -d'=' -f2 || echo 'not set')"; \
		echo "SSL Mode: $$(grep '^DB_SSLMODE=' backend/.env | cut -d'=' -f2 || echo 'not set')"; \
		if grep -q '^DATABASE_URL=' backend/.env; then \
			echo "Using DATABASE_URL: Yes"; \
		else \
			echo "Using DATABASE_URL: No"; \
		fi; \
	elif [ -f backend/config.env ]; then \
		echo "Source: backend/config.env (legacy)"; \
		echo "Provider: $$(grep '^DB_PROVIDER=' backend/config.env | cut -d'=' -f2 || echo 'postgres (default)')"; \
		echo "Host: $$(grep '^DB_HOST=' backend/config.env | cut -d'=' -f2 || echo 'not set')"; \
		echo "Port: $$(grep '^DB_PORT=' backend/config.env | cut -d'=' -f2 || echo 'not set')"; \
		echo "Database: $$(grep '^DB_NAME=' backend/config.env | cut -d'=' -f2 || echo 'not set')"; \
		echo "SSL Mode: $$(grep '^DB_SSLMODE=' backend/config.env | cut -d'=' -f2 || echo 'not set')"; \
		if grep -q '^DATABASE_URL=' backend/config.env; then \
			echo "Using DATABASE_URL: Yes"; \
		else \
			echo "Using DATABASE_URL: No"; \
		fi; \
	else \
		echo "❌ backend/.env not found!"; \
		echo "Create one from the template: cp backend/config.env.example backend/.env"; \
	fi

# Health checks
health: ## Check health of all services
	@echo "Checking service health..."
	@curl -f http://localhost:8080/health && echo "Backend: ✅" || echo "Backend: ❌"
	@curl -f http://localhost:3000 && echo "Frontend: ✅" || echo "Frontend: ❌"

# Frontend development
frontend-dev: ## Run frontend in development mode (requires dev environment)
	@echo "Starting frontend in development mode..."
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	@echo "Building frontend for production..."
	cd frontend && npm run build

frontend-install: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Full stack development
full-dev: ## Start full development environment (database + backend + frontend)
	@echo "Starting full development environment..."
	cd backend && docker compose -f docker-compose.dev.yml up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 10
	@echo "Starting backend with Air..."
	cd backend && air &
	@sleep 5
	@echo "Starting frontend..."
	cd frontend && npm run dev

# Quick start for development
quick-start: ## Quick start for development (database + backend with Air)
	@echo "Quick starting development environment..."
	cd backend && docker compose -f docker-compose.dev.yml up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 10
	cd backend && air

quick-start-simple: ## Quick start for development (database + backend simple)
	@echo "Quick starting development environment..."
	cd backend && docker compose -f docker-compose.dev.yml up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 10
	cd backend && go run cmd/server/main.go
