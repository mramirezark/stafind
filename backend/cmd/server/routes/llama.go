package routes

import (
	"stafind-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupLlamaAIRoutes configures Llama AI routes with real data processing and API key authentication
func SetupLlamaAIRoutes(app *fiber.App) {
	// DEBUG: Add a simple test route first
	app.Get("/llama/debug", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Llama routes are working"})
	})

	// Initialize the real data processor with your API key
	processor := services.NewRealLlamaProcessor("8264dc9f07e749d9c2ffead0b25de8cb22bed7af774e189ef224ae015908776b")

	// DEBUG: Check if processor is nil
	if processor == nil {
		panic("Processor is nil!")
	}

	// Llama AI routes group with API key authentication
	llama := app.Group("/llama")
	{
		// Health check endpoint
		llama.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status":    "healthy",
				"service":   "llama-ai",
				"timestamp": "2025-09-30T15:44:13Z",
			})
		})

		// Generic text processing endpoint
		llama.Post("/process", func(c *fiber.Ctx) error {
			// Validate API key
			apiKey := c.Get("X-API-Key")
			if apiKey != "8264dc9f07e749d9c2ffead0b25de8cb22bed7af774e189ef224ae015908776b" {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid API key"})
			}

			var request struct {
				Text           string                 `json:"text" binding:"required"`
				ProcessingType string                 `json:"processing_type" binding:"required"`
				Metadata       map[string]interface{} `json:"metadata,omitempty"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
			}

			// Process with real data
			response := processor.ProcessWithRealData(request.Text, request.ProcessingType)
			return c.JSON(response)
		})

		// Candidate extraction endpoint
		llama.Post("/extract-candidate", func(c *fiber.Ctx) error {
			// Validate API key
			apiKey := c.Get("X-API-Key")
			if apiKey != "8264dc9f07e749d9c2ffead0b25de8cb22bed7af774e189ef224ae015908776b" {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid API key"})
			}

			var request struct {
				Text     string `json:"text" binding:"required"`
				Language string `json:"language,omitempty"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
			}

			// Process with real data
			response := processor.ProcessWithRealData(request.Text, "candidate_extraction")
			return c.JSON(response)
		})

		// Search analysis endpoint
		llama.Post("/analyze-search", func(c *fiber.Ctx) error {
			// Validate API key
			apiKey := c.Get("X-API-Key")
			if apiKey != "8264dc9f07e749d9c2ffead0b25de8cb22bed7af774e189ef224ae015908776b" {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid API key"})
			}

			var request struct {
				Text     string `json:"text" binding:"required"`
				Language string `json:"language,omitempty"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
			}

			// Process with real data
			response := processor.ProcessWithRealData(request.Text, "search_analysis")
			return c.JSON(response)
		})

		// Candidate matching endpoint
		llama.Post("/match-candidate", func(c *fiber.Ctx) error {
			// Validate API key
			apiKey := c.Get("X-API-Key")
			if apiKey != "8264dc9f07e749d9c2ffead0b25de8cb22bed7af774e189ef224ae015908776b" {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid API key"})
			}

			var request struct {
				CandidateName  string `json:"candidate_name" binding:"required"`
				CandidateInfo  string `json:"candidate_info" binding:"required"`
				SearchCriteria string `json:"search_criteria" binding:"required"`
				Language       string `json:"language,omitempty"`
			}

			if err := c.BodyParser(&request); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
			}

			// Process with real data
			response := processor.ProcessWithRealData(request.CandidateInfo, "candidate_matching")
			return c.JSON(response)
		})
	}
}
