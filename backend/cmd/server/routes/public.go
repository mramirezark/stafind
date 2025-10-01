package routes

import (
	"stafind-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupPublicRoutes configures public routes (no authentication required)
func SetupPublicRoutes(app *fiber.App, h *handlers.Handlers, apiKeyHandlers *handlers.APIKeyHandlers) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Test routes for debugging
	app.Get("/public-test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Public routes are working"})
	})

	app.Post("/public-test-post", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Public POST routes are working"})
	})

	// Public API key validation endpoints
	app.Post("/api-keys/validate", apiKeyHandlers.ValidateAPIKey)
	app.Get("/api-keys/test", apiKeyHandlers.TestAPIKey)

	// Public AI agent endpoint for testing
	app.Post("/ai-agent/process", h.AIAgentHandlers.ProcessAIAgentRequest)
}
