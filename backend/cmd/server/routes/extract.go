package routes

import (
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupExtractRoutes configures extraction routes using pure NER with API key authentication
func SetupExtractRoutes(app *fiber.App, h *handlers.ExtractHandlers) {
	app.Get("/api/v1/extract/health", h.HealthCheck)

	apiShort := app.Group("/api/v1/extract", middleware.APIKeyMiddleware())
	{
		apiShort.Post("/process", h.ConsolidatedExtract)
	}
}
