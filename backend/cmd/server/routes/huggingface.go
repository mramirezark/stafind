package routes

import (
	"stafind-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupHuggingFaceRoutes sets up routes for Hugging Face skill extraction
func SetupHuggingFaceRoutes(app *fiber.App, huggingFaceHandlers *handlers.HuggingFaceHandlers) {
	// Create a group for Hugging Face routes
	huggingface := app.Group("/api/huggingface")

	// Skill extraction endpoints
	huggingface.Post("/extract-skills", huggingFaceHandlers.ExtractSkills)
	huggingface.Post("/extract-skills-simple", huggingFaceHandlers.ExtractSkillsFromText)
	huggingface.Post("/batch-extract-skills", huggingFaceHandlers.BatchExtractSkills)
	huggingface.Post("/compare-models", huggingFaceHandlers.CompareModels)

	// Model management endpoints
	huggingface.Get("/models", huggingFaceHandlers.GetAvailableModels)
	huggingface.Get("/models/:model/config", huggingFaceHandlers.GetModelConfig)

	// Statistics and monitoring endpoints
	huggingface.Get("/stats", huggingFaceHandlers.GetStats)
	huggingface.Get("/health", huggingFaceHandlers.HealthCheck)
}
