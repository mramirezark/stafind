# Docker Setup for StaffFind

This guide will help you set up and run StaffFind using Docker containers.

## 🚀 Quick Start

### Option 1: Using the Makefile (Recommended)

```bash
# Start development environment (PostgreSQL + pgAdmin)
make dev

# Start production environment (full stack)
make prod

# Stop all services
make stop

# Show help
make help
```

### Option 2: Using the startup script

```bash
# Make the script executable (first time only)
chmod +x start-docker.sh

# Start development environment
./start-docker.sh dev

# Start production environment
./start-docker.sh prod

# Stop all services
./start-docker.sh stop
```

### Option 3: Using Docker Compose directly

```bash
# Development environment
cd backend
docker-compose -f docker-compose.dev.yml up -d

# Production environment
docker-compose up -d
```

## 📁 Docker Structure

```
stafind/
├── docker-compose.yml              # Production environment
├── backend/
│   ├── Dockerfile                  # Backend container
│   ├── docker-compose.dev.yml      # Development environment
│   ├── docker.env                  # Docker environment variables
│   └── init-db.sql                 # Database initialization
├── frontend/
│   └── Dockerfile                  # Frontend container
├── Makefile                        # Docker management commands
└── start-docker.sh                 # Startup script
```

## 🛠️ Development Environment

### Start Development Services

The development environment includes:
- **PostgreSQL**: Database server on port 5432
- **pgAdmin**: Database management interface on port 5050

```bash
make dev
```

### Access Services

- **Database**: `localhost:5432`
  - Database: `stafind`
  - Username: `postgres`
  - Password: `password`

- **pgAdmin**: `http://localhost:5050`
  - Email: `admin@stafind.com`
  - Password: `admin`

### Run Backend Locally

After starting the development database:

```bash
# With Air for live reloading (recommended)
cd backend
air

# Or simple mode without live reloading
cd backend
go run cmd/server/main.go
```

The backend will automatically run Flyway migrations on startup.

## 🏭 Production Environment

### Start Production Services

The production environment includes:
- **PostgreSQL**: Database server
- **Backend**: Go API server on port 8080
- **Frontend**: React application on port 3000

```bash
make prod
```

### Access Services

- **Backend API**: `http://localhost:8080`
- **Frontend**: `http://localhost:3000`
- **Database**: `localhost:5432`

## 📊 Database Management

### Run Migrations

```bash
make migrate
```

### Connect to Database

```bash
# Using Docker
make db-shell

# Or directly with psql
docker-compose exec postgres psql -U postgres -d stafind
```

### Reset Database

⚠️ **Warning**: This will delete all data!

```bash
make db-reset
```

## 🔍 Monitoring and Logs

### View Logs

```bash
# All services
make logs

# Specific service
make logs-backend
make logs-db
```

### Check Status

```bash
make status
```

### Health Checks

```bash
make health
```

## 🛠️ Available Commands

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start development environment |
| `make prod` | Start production environment |
| `make stop` | Stop all services |
| `make logs` | Show logs for all services |
| `make logs-backend` | Show backend logs |
| `make logs-db` | Show database logs |
| `make migrate` | Run database migrations |
| `make status` | Show status of all services |
| `make build-backend` | Build backend Docker image |
| `make build-frontend` | Build frontend Docker image |
| `make build-all` | Build all Docker images |
| `make cleanup` | Remove all Docker resources |
| `make db-shell` | Connect to PostgreSQL shell |
| `make db-reset` | Reset database (⚠️ destructive) |
| `make health` | Check health of all services |
| `make quick-start` | Quick start for development (with Air) |
| `make backend-air` | Run backend with Air (live reloading) |
| `make backend-dev` | Run backend in development mode |
| `make quick-start-simple` | Quick start without Air |

### Startup Script Commands

```bash
./start-docker.sh dev      # Start development environment
./start-docker.sh dev-air  # Start development with Air (live reloading)
./start-docker.sh prod     # Start production environment
./start-docker.sh stop     # Stop all services
./start-docker.sh logs     # Show logs
./start-docker.sh migrate  # Run migrations
./start-docker.sh status   # Show status
./start-docker.sh cleanup  # Clean up Docker resources
```

## 🔧 Configuration

### Environment Variables

The Docker setup uses the following environment variables:

#### Database Configuration
```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=stafind
DB_SSLMODE=disable
```

#### Server Configuration
```bash
PORT=8080
HOST=0.0.0.0
```

#### Flyway Configuration
```bash
FLYWAY_LOCATIONS=./flyway_migrations
```

### Custom Configuration

To use custom environment variables:

1. **Copy the environment file**:
   ```bash
   cp backend/docker.env backend/.env
   ```

2. **Modify the values** in `backend/.env`

3. **Update docker-compose.yml** to use the custom env file:
   ```yaml
   backend:
     env_file:
       - ./backend/.env
   ```

## 🐛 Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
Error: bind: address already in use
```
**Solution**: Stop the service using the port or change the port in docker-compose.yml

#### 2. Database Connection Refused
```bash
Error: dial tcp [::1]:5432: connect: connection refused
```
**Solution**: Ensure PostgreSQL container is running and healthy:
```bash
make status
docker-compose logs postgres
```

#### 3. Migration Failures
```bash
Error: migration failed
```
**Solution**: Check database connection and migration files:
```bash
make migrate
make logs-backend
```

#### 4. Container Won't Start
```bash
Error: container failed to start
```
**Solution**: Check container logs:
```bash
docker-compose logs backend
```

### Debugging Commands

```bash
# Check container status
docker ps -a

# Check logs
docker-compose logs -f

# Check resource usage
docker stats

# Remove all containers and volumes
make cleanup
```

### Reset Everything

If you encounter persistent issues:

```bash
# Stop all services
make stop

# Remove all containers, volumes, and images
make cleanup

# Start fresh
make dev
```

## 📈 Performance Tips

### Development

1. **Use development environment** for local development
2. **Run backend locally** for faster iteration
3. **Use pgAdmin** for database inspection

### Production

1. **Use production environment** for testing
2. **Monitor resource usage** with `docker stats`
3. **Use health checks** to monitor service status

## 🔒 Security Considerations

### Development

- Default passwords are used for convenience
- Services are exposed on localhost only
- pgAdmin is included for database management

### Production

- Change default passwords
- Use environment variables for secrets
- Configure proper CORS settings
- Use HTTPS in production
- Consider using Docker secrets for sensitive data

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [pgAdmin Documentation](https://www.pgadmin.org/docs/)
- [Flyway Migrations Guide](backend/FLYWAY_MIGRATIONS.md)

## 🆘 Getting Help

If you encounter issues:

1. Check the logs: `make logs`
2. Verify service status: `make status`
3. Check health: `make health`
4. Review this documentation
5. Check the [Flyway Migrations Guide](backend/FLYWAY_MIGRATIONS.md)
