# StaffFind - Employee Matching System

StaffFind is a comprehensive application that allows employees to submit job descriptions or tech stack requests and receive a curated list of employees within the company who best match the requested skills.

## Features

- **Job Request Management**: Create and manage job descriptions with required and preferred skills
- **Employee Profiles**: Maintain employee profiles with skills, experience levels, and proficiency ratings
- **Smart Matching Algorithm**: Advanced algorithm that matches employees based on:
  - Required and preferred skills
  - Skill proficiency levels
  - Years of experience
  - Department alignment
  - Experience level matching
  - Location preferences
- **Search & Discovery**: Powerful search interface to find employees based on various criteria
- **Modern UI**: Built with React and Material-UI for a responsive, intuitive interface

## Tech Stack

### Backend
- **Go 1.21+** - Latest Go version for high performance
- **Fiber** - Express-inspired HTTP web framework for Go
- **PostgreSQL** - Database for data persistence
- **goflyway** - Flyway-style database migrations
- **Layered Architecture** - Handler-Service-Repository pattern
- **CORS** - Cross-origin resource sharing support

### Frontend
- **Next.js 14** - Full-stack React framework with App Router
- **React 18** - Latest React with TypeScript
- **Material-UI (MUI)** - Modern component library
- **TypeScript** - Full type safety
- **Axios** - HTTP client for API communication

## Project Structure

```
stafind/
├── backend/
│   ├── cmd/server/          # Main server application
│   ├── internal/
│   │   ├── models/         # Data models and DTOs
│   │   ├── handlers/       # HTTP request handlers (Fiber)
│   │   ├── services/       # Business logic layer
│   │   ├── repositories/   # Data access layer
│   │   ├── database/       # Database connection and migrations
│   │   └── matching/       # Matching algorithm implementation
│   ├── flyway_migrations/  # Flyway-style migration files
│   └── config.env.example  # Environment configuration template
├── frontend/
│   ├── app/                # Next.js App Router
│   ├── components/         # React components
│   ├── lib/                # Utilities and API client
│   ├── next.config.js      # Next.js configuration
│   └── package.json        # Frontend dependencies
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.24 or later
- Node.js 18 or later
- PostgreSQL 15 or later
- npm or yarn
- Docker (optional)

### Quick Start with Docker

```bash
# Start development environment
make dev

# Start backend with Air (live reloading)
make backend-air

# Start frontend development
make frontend-dev

# Or start full stack
make full-dev
```

### Manual Setup

#### Backend Setup

1. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables:**
   ```bash
   cp config.env.example .env
   ```
   
   Edit `.env` with your database configuration:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=stafind
   PORT=8080
   GIN_MODE=debug
   ```

4. **Set up PostgreSQL database:**
   ```bash
   createdb stafind
   ```

5. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

   The API will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

   The application will be available at `http://localhost:3000`

4. **Build for production:**
   ```bash
   npm run build
   npm start
   ```

## API Endpoints

### Employees
- `GET /api/v1/employees` - Get all employees
- `GET /api/v1/employees/:id` - Get employee by ID
- `POST /api/v1/employees` - Create new employee
- `PUT /api/v1/employees/:id` - Update employee
- `DELETE /api/v1/employees/:id` - Delete employee

### Job Requests
- `GET /api/v1/job-requests` - Get all job requests
- `GET /api/v1/job-requests/:id` - Get job request by ID
- `POST /api/v1/job-requests` - Create new job request
- `GET /api/v1/job-requests/:id/matches` - Get matches for job request

### Search
- `POST /api/v1/search` - Search employees by criteria

### Skills
- `GET /api/v1/skills` - Get all available skills
- `POST /api/v1/skills` - Create new skill

## Usage

### Creating a Job Request

1. Navigate to the "Job Request" tab
2. Fill in the job title and description
3. Select required and preferred skills
4. Set experience level, department, and location preferences
5. Submit the request

### Managing Employees

1. Go to the "Employees" tab
2. Add new employees with their skills and proficiency levels
3. Edit existing employee profiles
4. Search and filter employees by various criteria

### Finding Matches

1. Use the "Search" tab to find employees based on specific criteria
2. Adjust the minimum match score slider to filter results
3. View detailed match information including matching skills
4. Contact matched employees for opportunities

## Architecture

### Backend Architecture
The backend follows a clean layered architecture pattern:

- **Handlers Layer**: HTTP request/response handling using Fiber framework
- **Services Layer**: Business logic and validation
- **Repositories Layer**: Data access and database operations
- **Models Layer**: Data transfer objects (DTOs) and domain models

This separation ensures:
- **Testability**: Each layer can be tested independently
- **Maintainability**: Clear separation of concerns
- **Scalability**: Easy to modify or extend individual layers
- **Reusability**: Services can be reused across different handlers

### Dependency Injection
The application uses dependency injection to wire up the layers:
```go
// Repositories (Data Access)
employeeRepo := repositories.NewEmployeeRepository(db.DB)
jobRequestRepo := repositories.NewJobRequestRepository(db.DB)

// Services (Business Logic)
employeeService := services.NewEmployeeService(employeeRepo)
jobRequestService := services.NewJobRequestService(jobRequestRepo)

// Handlers (HTTP Layer)
handlers := handlers.NewHandlers(employeeService, jobRequestService, ...)
```

## Matching Algorithm

The matching algorithm considers multiple factors:

1. **Required Skills** (Weight: 3.0)
   - Base score for having the skill
   - Proficiency level bonus (1-5 scale)
   - Years of experience bonus

2. **Preferred Skills** (Weight: 1.0)
   - Similar scoring to required skills but with lower weight

3. **Department Match** (Weight: 2.0)
   - Bonus for matching department

4. **Experience Level** (Weight: 1.5)
   - Bonus for meeting or exceeding required level

5. **Location Match** (Weight: 1.0)
   - Bonus for matching location

The algorithm also applies a coverage multiplier to encourage higher skill coverage matches.

## Development

### Running Tests

Backend tests:
```bash
cd backend
go test ./...
```

Frontend tests:
```bash
cd frontend
npm test
```

### Database Migrations

The project uses **goflyway** for Flyway-style database migrations:

#### Using Flyway CLI

```bash
# Run all pending migrations
go run cmd/flyway-cli/main.go migrate

# Show migration information
go run cmd/flyway-cli/main.go info

# Validate migration files
go run cmd/flyway-cli/main.go validate

# List pending migrations
go run cmd/flyway-cli/main.go pending

# List applied migrations
go run cmd/flyway-cli/main.go applied
```

#### Migration Files

Migrations are located in `flyway_migrations/` and follow Flyway naming conventions:
- `V1__Create_users_and_roles.sql` - User authentication and role system
- `V2__Create_initial_schema.sql` - Core application tables
- `V3__Insert_sample_skills.sql` - Sample skill data
- `V4__Add_sample_employees.sql` - Sample employee profiles
- `V5__Add_sample_job_requests.sql` - Sample job requests

#### Creating New Migrations

1. Create a new file: `V5__Your_description.sql`
2. Write your SQL migration
3. Validate: `go run cmd/flyway-cli/main.go validate`
4. Apply: `go run cmd/flyway-cli/main.go migrate`

For detailed information, see [FLYWAY_MIGRATIONS.md](backend/FLYWAY_MIGRATIONS.md).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
