package main

import (
	"os"
	"stafind-backend/cmd/server/routes"
	"stafind-backend/internal/constants"
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
	// Try .env first (standard), then fall back to config.env (legacy)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("config.env"); err != nil {
			log.Info("No .env or config.env file found, using environment variables")
		} else {
			log.Info("Loaded environment from config.env (consider renaming to .env)")
		}
	} else {
		log.Info("Loaded environment from .env")
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
	categoryRepo, err := repositories.NewCategoryRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize category repository", "error", err)
	}

	userRepo, err := repositories.NewUserRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize user repository", "error", err)
	}

	roleRepo, err := repositories.NewRoleRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize role repository", "error", err)
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

	cvExtractRepo, err := repositories.NewCVExtractRepository(db.DB)
	if err != nil {
		log.Fatal("Failed to initialize CV extract repository", "error", err)
	}

	// Initialize services
	employeeService := services.NewEmployeeService(employeeRepo)
	searchService := services.NewSearchService(employeeRepo)
	skillService := services.NewSkillService(skillRepo, employeeRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo)
	dashboardService := services.NewDashboardService(employeeRepo, skillRepo, aiAgentRepo, matchRepo)
	notificationService := services.NewNotificationService(aiAgentRepo)
	aiAgentService := services.NewAIAgentService(aiAgentRepo, employeeRepo, skillRepo, categoryRepo, matchRepo, notificationService)
	nerService := services.NewNERService(skillRepo, categoryRepo)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	extractionService := services.NewCandidateExtractionService(skillRepo, categoryRepo)
	candidateStorageService := services.NewCandidateStorageService(employeeRepo, skillRepo)
	cvExtractService := services.NewCVExtractService(cvExtractRepo)

	// Initialize Hugging Face service
	huggingFaceAPIKey := os.Getenv(constants.EnvHuggingFaceAPIKey)
	if huggingFaceAPIKey == "" {
		log.Warn("HUGGINGFACE_API_KEY not set, Hugging Face service will not be available")
	}
	huggingFaceService := services.NewHuggingFaceSkillService(huggingFaceAPIKey)

	// Initialize handlers
	h := handlers.NewHandlers(employeeService, searchService, skillService, categoryService, aiAgentService, nerService)
	authHandlers := handlers.NewAuthHandlers(userService, roleService)
	dashboardHandlers := handlers.NewDashboardHandlers(dashboardService)
	apiKeyHandlers := handlers.NewAPIKeyHandlers(apiKeyService)
	extractionHandlers := handlers.NewExtractHandlers(extractionService, aiAgentService, candidateStorageService, cvExtractService)
	matchingHandlers := handlers.NewMatchingHandler(aiAgentService)
	cvExtractHandlers := handlers.NewCVExtractHandlers(cvExtractService)
	huggingFaceHandlers := handlers.NewHuggingFaceHandlers(huggingFaceService)
	combinedExtractHandlers := handlers.NewCombinedExtractHandlers(extractionService, aiAgentService, candidateStorageService, cvExtractService, huggingFaceService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup routes using enhanced structure
	app := routes.SetupAllRoutes(h, authHandlers, dashboardHandlers, apiKeyHandlers, extractionHandlers, matchingHandlers, cvExtractHandlers, huggingFaceHandlers, combinedExtractHandlers)

	log.Info("Server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}
