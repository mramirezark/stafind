#!/bin/bash

echo "Starting StaffFind Application..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Please start Docker first."
    exit 1
fi

# Start the application using Docker Compose
echo "Starting services with Docker Compose..."
docker-compose up --build

echo "Application started!"
echo "Frontend: http://localhost:3000"
echo "Backend API: http://localhost:8080"
echo "PostgreSQL: localhost:5432"
