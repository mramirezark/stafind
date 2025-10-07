package routes

import (
	"stafind-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupAllRoutes configures all application routes with enhanced structure
func SetupAllRoutes(
	h *handlers.Handlers,
	authHandlers *handlers.AuthHandlers,
	dashboardHandlers *handlers.DashboardHandlers,
	apiKeyHandlers *handlers.APIKeyHandlers,
	extractionHandlers *handlers.ExtractHandlers,
	matchingHandlers *handlers.MatchingHandler,
	cvExtractHandlers *handlers.CVExtractHandlers,
	huggingFaceHandlers *handlers.HuggingFaceHandlers,
	combinedExtractHandlers *handlers.CombinedExtractHandlers,
) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Add global middleware
	app.Use(fiberLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-API-Key",
	}))

	// Setup route groups in order of priority
	SetupPublicRoutes(app, h, apiKeyHandlers)
	SetupAuthRoutes(app, authHandlers)
	SetupExtractRoutes(app, extractionHandlers)
	SetupCombinedExtractRoutes(app, combinedExtractHandlers)
	SetupMatchingRoutes(app, matchingHandlers)
	SetupCVExtractRoutes(app, cvExtractHandlers)
	SetupHuggingFaceRoutes(app, huggingFaceHandlers)
	SetupAPIRoutes(app, h, authHandlers, dashboardHandlers, apiKeyHandlers)
	SetupAdminRoutes(app, authHandlers, apiKeyHandlers)

	return app
}

// SetupCVExtractRoutes configures CV extract tracking routes
func SetupCVExtractRoutes(app *fiber.App, cvExtractHandlers *handlers.CVExtractHandlers) {
	// Register CV extract routes
	cvExtractHandlers.RegisterCVExtractRoutes(app)
}
