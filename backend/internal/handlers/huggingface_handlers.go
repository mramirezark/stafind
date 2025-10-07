package handlers

import (
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HuggingFaceHandlers handles Hugging Face skill extraction endpoints
type HuggingFaceHandlers struct {
	huggingFaceService services.HuggingFaceSkillExtractionService
}

// NewHuggingFaceHandlers creates new Hugging Face handlers
func NewHuggingFaceHandlers(huggingFaceService services.HuggingFaceSkillExtractionService) *HuggingFaceHandlers {
	return &HuggingFaceHandlers{
		huggingFaceService: huggingFaceService,
	}
}

// ExtractSkills extracts skills from text using Hugging Face models
func (h *HuggingFaceHandlers) ExtractSkills(c *fiber.Ctx) error {
	var request models.HuggingFaceSkillExtractionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if request.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "text field is required",
		})
	}

	// Set default values if not provided
	if request.ConfidenceThreshold == 0 {
		request.ConfidenceThreshold = 0.5
	}
	if request.MaxSkills == 0 {
		request.MaxSkills = 50
	}

	// Extract skills
	response, err := h.huggingFaceService.ExtractSkills(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to extract skills",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ExtractSkillsFromText is a simplified endpoint for basic skill extraction
func (h *HuggingFaceHandlers) ExtractSkillsFromText(c *fiber.Ctx) error {
	var request struct {
		Text string `json:"text" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if request.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "text field is required",
		})
	}

	// Extract skills using default settings
	response, err := h.huggingFaceService.ExtractSkillsFromText(request.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to extract skills",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetAvailableModels returns the list of available Hugging Face models
func (h *HuggingFaceHandlers) GetAvailableModels(c *fiber.Ctx) error {
	models, err := h.huggingFaceService.GetAvailableModels()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to get available models",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"models":  models,
	})
}

// GetModelConfig returns configuration for a specific model
func (h *HuggingFaceHandlers) GetModelConfig(c *fiber.Ctx) error {
	modelName := c.Params("model")
	if modelName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "model parameter is required",
		})
	}

	config, err := h.huggingFaceService.GetModelConfig(modelName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Model not found",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"config":  config,
	})
}

// GetStats returns statistics for the skill extraction service
func (h *HuggingFaceHandlers) GetStats(c *fiber.Ctx) error {
	stats, err := h.huggingFaceService.GetStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to get statistics",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

// HealthCheck provides a health check endpoint for the Hugging Face service
func (h *HuggingFaceHandlers) HealthCheck(c *fiber.Ctx) error {
	err := h.huggingFaceService.HealthCheck()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "unhealthy",
			"error":     err.Error(),
			"timestamp": time.Now(),
			"service":   "Hugging Face Skill Extraction Service",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "Hugging Face Skill Extraction Service",
	})
}

// BatchExtractSkills extracts skills from multiple texts in batch
func (h *HuggingFaceHandlers) BatchExtractSkills(c *fiber.Ctx) error {
	var request struct {
		Texts               []string `json:"texts" validate:"required"`
		ModelName           string   `json:"model_name,omitempty"`
		ConfidenceThreshold float64  `json:"confidence_threshold,omitempty"`
		MaxSkills           int      `json:"max_skills,omitempty"`
		Categories          []string `json:"categories,omitempty"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if len(request.Texts) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "texts array is required and cannot be empty",
		})
	}

	if len(request.Texts) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Maximum 10 texts allowed per batch request",
		})
	}

	// Process each text
	var results []models.HuggingFaceSkillExtractionResponse
	var errors []string

	for i, text := range request.Texts {
		extractRequest := models.HuggingFaceSkillExtractionRequest{
			Text:                text,
			ModelName:           request.ModelName,
			ConfidenceThreshold: request.ConfidenceThreshold,
			MaxSkills:           request.MaxSkills,
			Categories:          request.Categories,
		}

		response, err := h.huggingFaceService.ExtractSkills(&extractRequest)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Text %d: %s", i+1, err.Error()))
			results = append(results, models.HuggingFaceSkillExtractionResponse{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			results = append(results, *response)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"results":    results,
		"errors":     errors,
		"total":      len(request.Texts),
		"successful": len(results) - len(errors),
		"failed":     len(errors),
	})
}

// CompareModels compares skill extraction results from different models
func (h *HuggingFaceHandlers) CompareModels(c *fiber.Ctx) error {
	var request struct {
		Text   string   `json:"text" validate:"required"`
		Models []string `json:"models" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if request.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "text field is required",
		})
	}

	if len(request.Models) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "models array is required",
		})
	}

	if len(request.Models) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Maximum 5 models allowed for comparison",
		})
	}

	// Extract skills using each model
	var comparisons []fiber.Map
	var errors []string

	for _, modelName := range request.Models {
		extractRequest := models.HuggingFaceSkillExtractionRequest{
			Text:      request.Text,
			ModelName: modelName,
		}

		response, err := h.huggingFaceService.ExtractSkills(&extractRequest)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Model %s: %s", modelName, err.Error()))
			comparisons = append(comparisons, fiber.Map{
				"model":   modelName,
				"success": false,
				"error":   err.Error(),
			})
		} else {
			comparisons = append(comparisons, fiber.Map{
				"model":           modelName,
				"success":         true,
				"total_skills":    response.TotalSkills,
				"confidence":      response.ConfidenceScore,
				"processing_time": response.ProcessingTime,
				"skills":          response.Skills,
				"categories":      response.Categories,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":      true,
		"text":         request.Text,
		"comparisons":  comparisons,
		"errors":       errors,
		"total_models": len(request.Models),
		"successful":   len(comparisons) - len(errors),
		"failed":       len(errors),
	})
}
