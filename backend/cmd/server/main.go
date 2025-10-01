package main

import (
	"os"
	"stafind-backend/cmd/server/routes"
	"stafind-backend/internal/database"
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/logger"
	"stafind-backend/internal/repositories"
	"stafind-backend/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize structured logging
	if err := logger.Init(nil); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	log := logger.Get()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Run Flyway migrations
	flywayConfig := database.NewFlywayConfigFromEnv()
	flywayMigrator, err := database.NewFlywayMigrator(flywayConfig)
	if err != nil {
		log.Fatal("Failed to initialize Flyway migrator", "error", err)
	}
	defer flywayMigrator.Close()

	// Validate migrations first
	if err := flywayMigrator.ValidateMigrations(); err != nil {
		log.Fatal("Migration validation failed", "error", err)
	}

	// Run migrations
	if err := flywayMigrator.Migrate(); err != nil {
		log.Fatal("Failed to run Flyway migrations", "error", err)
	}

	// Initialize repositories
	employeeRepo, err := repositories.NewEmployeeRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize employee repository", "error", err)
	}

	skillRepo, err := repositories.NewSkillRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize skill repository", "error", err)
	}

	userRepo, err := repositories.NewUserRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize user repository", "error", err)
	}

	roleRepo, err := repositories.NewRoleRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize role repository", "error", err)
	}

	uploadedFileRepo, err := repositories.NewUploadedFileRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize uploaded file repository", "error", err)
	}
	aiAgentRepo, err := repositories.NewAIAgentRepository(db)
	if err != nil {
		log.Fatal("Failed to initialize AI agent repository", "error", err)
	}
	apiKeyRepo, err := repositories.NewAPIKeyRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize API key repository", "error", err)
	}
	matchRepo, err := repositories.NewMatchRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize match repository", "error", err)
	}

	// Initialize services
	employeeService := services.NewEmployeeService(employeeRepo)
	searchService := services.NewSearchService(employeeRepo)
	skillService := services.NewSkillService(skillRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo)
	dashboardService := services.NewDashboardService(employeeRepo, skillRepo, aiAgentRepo, matchRepo)
	uploadedFileService := services.NewUploadedFileService(uploadedFileRepo)
	notificationService := services.NewNotificationService(aiAgentRepo)
	aiAgentService := services.NewAIAgentService(aiAgentRepo, employeeRepo, skillRepo, matchRepo, notificationService)
	nerService := services.NewNERService()
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	llamaAIService := services.NewLlamaAIService()

	// Initialize handlers
	h := handlers.NewHandlers(employeeService, searchService, skillService, aiAgentService, nerService, llamaAIService)
	authHandlers := handlers.NewAuthHandlers(userService, roleService)
	dashboardHandlers := handlers.NewDashboardHandlers(dashboardService)
	fileUploadHandlers := handlers.NewFileUploadHandlers(uploadedFileService)
	apiKeyHandlers := handlers.NewAPIKeyHandlers(apiKeyService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup routes using enhanced structure
	app := routes.SetupAllRoutes(h, authHandlers, dashboardHandlers, fileUploadHandlers, apiKeyHandlers)

	log.Info("Server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
