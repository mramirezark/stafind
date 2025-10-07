package routes

import (
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupExtractRoutes configures extraction routes using pure NER with API key authentication
func SetupExtractRoutes(app *fiber.App, h *handlers.ExtractHandlers) {
	apiShort := app.Group("/api/v1/extract", middleware.APIKeyMiddleware())
	{
		apiShort.Post("/process", h.ExtractProcess)
	}
}

// SetupCombinedExtractRoutes configures combined NER and Hugging Face extraction routes
func SetupCombinedExtractRoutes(app *fiber.App, h *handlers.CombinedExtractHandlers) {
	apiShort := app.Group("/api/v1/extract", middleware.APIKeyMiddleware())
	{
		apiShort.Post("/process-combined", h.ExtractProcessCombined)
		apiShort.Post("/compare-methods", h.CompareExtractionMethods)
	}
}

// SetupMatchingRoutes configures employee matching routes
func SetupMatchingRoutes(app *fiber.App, h *handlers.MatchingHandler) {
	apiShort := app.Group("/api/v1/matching", middleware.APIKeyMiddleware())
	{
		apiShort.Post("/employees", h.FindMatchingEmployees)
		apiShort.Get("/history", h.GetMatchingHistory)
	}
}
