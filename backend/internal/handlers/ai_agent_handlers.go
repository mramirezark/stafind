package handlers

import (
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AIAgentHandlers struct {
	aiAgentService services.AIAgentService
}

func NewAIAgentHandlers(aiAgentService services.AIAgentService) *AIAgentHandlers {
	return &AIAgentHandlers{
		aiAgentService: aiAgentService,
	}
}

// ProcessAIAgentRequest creates and processes AI agent requests synchronously
func (h *AIAgentHandlers) ProcessAIAgentRequest(c *fiber.Ctx) error {
	// Parse request
	var req models.CreateAIAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, err.Error())
	}

	// For public endpoint, generate unique Teams message ID if not provided
	if req.TeamsMessageID == "" {
		req.TeamsMessageID = fmt.Sprintf("public-test-%d", time.Now().UnixNano())
	}
	if req.ChannelID == "" {
		req.ChannelID = "public-test-channel"
	}
	if req.UserID == "" {
		req.UserID = "public-test-user"
	}
	if req.UserName == "" {
		req.UserName = "Public Test User"
	}

	// Use description and skills from public endpoint
	if req.Description != "" && req.MessageText == "" {
		req.MessageText = req.Description
	}

	// Create the AI agent request
	aiRequest, err := h.aiAgentService.CreateAIAgentRequest(&req)
	if err != nil {
		return InternalServerError(c, err.Error())
	}

	// Process synchronously and return results
	response, err := h.aiAgentService.ProcessAIAgentRequest(aiRequest.ID)
	if err != nil {
		return InternalServerErrorWithDetails(c, "Failed to process request", err.Error())
	}

	return Success(c, "Request processed successfully", response)
}

// GetAIAgentRequest returns a specific AI agent request by ID
func (h *AIAgentHandlers) GetAIAgentRequest(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid request ID")
	}

	aiRequest, err := h.aiAgentService.GetAIAgentRequest(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "AI agent request not found"})
	}

	return c.JSON(aiRequest)
}

// GetAIAgentResponse returns the AI agent response for a specific request ID
func (h *AIAgentHandlers) GetAIAgentResponse(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid request ID")
	}

	response, err := h.aiAgentService.GetAIAgentResponse(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "AI agent response not found"})
	}

	return c.JSON(response)
}

// ProcessAIAgentRequestByID processes an existing AI agent request by ID
func (h *AIAgentHandlers) ProcessAIAgentRequestByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "Invalid request ID")
	}

	response, err := h.aiAgentService.ProcessAIAgentRequest(id)
	if err != nil {
		return InternalServerError(c, err.Error())
	}

	return c.JSON(response)
}

// ExtractSkills extracts skills from text
func (h *AIAgentHandlers) ExtractSkills(c *fiber.Ctx) error {
	var req models.SkillExtractionRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, err.Error())
	}

	response, err := h.aiAgentService.ExtractSkillsFromText(req.Text)
	if err != nil {
		return InternalServerError(c, err.Error())
	}

	return c.JSON(response)
}

// GetAIAgentRequests returns all AI agent requests with pagination
func (h *AIAgentHandlers) GetAIAgentRequests(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	requests, err := h.aiAgentService.GetAIAgentRequests(limit, offset)
	if err != nil {
		return InternalServerError(c, err.Error())
	}

	return c.JSON(requests)
}
