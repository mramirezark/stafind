#!/bin/bash

# StaffFind Docker Startup Script
# This script helps you manage the Docker environment for StaffFind

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  StaffFind Docker Management${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Function to start development environment
start_dev() {
    print_status "Starting development environment..."
    print_status "Starting PostgreSQL database..."
    
    cd backend
    docker-compose -f docker-compose.dev.yml up -d
    
    print_status "Waiting for database to be ready..."
    sleep 10
    
    print_status "Development environment is ready!"
    print_status "Database: localhost:5432"
    print_status "pgAdmin: http://localhost:5050 (admin@stafind.com / admin)"
    print_status ""
    print_status "You can now run the backend locally:"
    print_status "cd backend && air                    # With live reloading"
    print_status "cd backend && go run cmd/server/main.go  # Simple mode"
}

# Function to start production environment
start_prod() {
    print_status "Starting production environment..."
    
    docker-compose up -d
    
    print_status "Waiting for services to be ready..."
    sleep 15
    
    print_status "Production environment is ready!"
    print_status "Backend API: http://localhost:8080"
    print_status "Frontend: http://localhost:3000"
    print_status "Database: localhost:5432"
}

# Function to stop all services
stop_all() {
    print_status "Stopping all services..."
    
    docker-compose down
    cd backend
    docker-compose -f docker-compose.dev.yml down
    
    print_status "All services stopped."
}

# Function to show logs
show_logs() {
    local service=${1:-""}
    
    if [ -z "$service" ]; then
        print_status "Showing logs for all services..."
        docker-compose logs -f
    else
        print_status "Showing logs for $service..."
        docker-compose logs -f "$service"
    fi
}

# Function to run migrations
run_migrations() {
    print_status "Running database migrations..."
    
    cd backend
    docker-compose -f docker-compose.dev.yml up -d postgres
    
    print_status "Waiting for database to be ready..."
    sleep 10
    
    # Run migrations using the CLI
    go run cmd/flyway-cli/main.go migrate
    
    print_status "Migrations completed!"
}

# Function to clean up
cleanup() {
    print_warning "This will remove all containers, volumes, and images. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_status "Cleaning up Docker resources..."
        
        docker-compose down -v --remove-orphans
        cd backend
        docker-compose -f docker-compose.dev.yml down -v --remove-orphans
        
        docker system prune -f
        
        print_status "Cleanup completed!"
    else
        print_status "Cleanup cancelled."
    fi
}

# Function to show status
show_status() {
    print_status "Docker service status:"
    docker-compose ps
    
    echo ""
    print_status "Development services status:"
    cd backend
    docker-compose -f docker-compose.dev.yml ps
}

# Function to start development with Air
start_dev_air() {
    print_status "Starting development environment with Air..."
    print_status "Starting PostgreSQL database..."
    
    cd backend
    docker-compose -f docker-compose.dev.yml up -d
    
    print_status "Waiting for database to be ready..."
    sleep 10
    
    print_status "Starting backend with Air (live reloading)..."
    air
}

# Function to show help
show_help() {
    print_header
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  dev          Start development environment (PostgreSQL + pgAdmin)"
    echo "  dev-air      Start development environment with Air (live reloading)"
    echo "  prod         Start production environment (full stack)"
    echo "  stop         Stop all services"
    echo "  logs [service] Show logs for all services or specific service"
    echo "  migrate      Run database migrations"
    echo "  status       Show status of all services"
    echo "  cleanup      Remove all containers, volumes, and images"
    echo "  help         Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev       # Start development environment"
    echo "  $0 dev-air   # Start development with Air (live reloading)"
    echo "  $0 prod      # Start production environment"
    echo "  $0 logs backend  # Show backend logs"
    echo "  $0 migrate   # Run database migrations"
}

# Main script logic
main() {
    check_docker
    
    case "${1:-help}" in
        "dev")
            start_dev
            ;;
        "dev-air")
            start_dev_air
            ;;
        "prod")
            start_prod
            ;;
        "stop")
            stop_all
            ;;
        "logs")
            show_logs "$2"
            ;;
        "migrate")
            run_migrations
            ;;
        "status")
            show_status
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Run main function with all arguments
main "$@"
