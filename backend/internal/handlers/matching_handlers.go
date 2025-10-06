package handlers

import (
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MatchingHandler struct {
	aiAgentService services.AIAgentService
}

func NewMatchingHandler(aiAgentService services.AIAgentService) *MatchingHandler {
	return &MatchingHandler{
		aiAgentService: aiAgentService,
	}
}

// MatchEmployeesRequest represents the request to find matching employees
type MatchEmployeesRequest struct {
	Skills []string `json:"skills"`
	Text   string   `json:"text,omitempty"` // Optional text to extract skills from
}

// MatchEmployeesResponse represents the response with matching employees
type MatchEmployeesResponse struct {
	Matches        []models.AIAgentMatch `json:"matches"`
	Summary        string                `json:"summary"`
	ProcessingTime int64                 `json:"processing_time_ms"`
	Status         string                `json:"status"`
	Message        string                `json:"message"`
}

// FindMatchingEmployees handles finding employees that match given skills
func (h *MatchingHandler) FindMatchingEmployees(c *fiber.Ctx) error {
	var request MatchEmployeesRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	startTime := time.Now()
	var skills []string
	var err error

	// Extract skills from text if provided, otherwise use provided skills
	if request.Text != "" {
		fmt.Printf("DEBUG: Extracting skills from text: %s\n", request.Text)
		skillExtractResult, extractErr := h.aiAgentService.ExtractSkillsFromText(request.Text)
		if extractErr != nil {
			fmt.Printf("Warning: Failed to extract skills from text: %v\n", extractErr)
			// Fall back to provided skills if extraction fails
			skills = request.Skills
		} else {
			skills = skillExtractResult.Skills
			fmt.Printf("DEBUG: Extracted skills: %v\n", skills)
		}
	} else {
		skills = request.Skills
		fmt.Printf("DEBUG: Using provided skills: %v\n", skills)
	}

	if len(skills) == 0 {
		response := MatchEmployeesResponse{
			Matches:        []models.AIAgentMatch{},
			Summary:        "No skills provided for matching",
			ProcessingTime: time.Since(startTime).Milliseconds(),
			Status:         "completed",
			Message:        "No skills to match against",
		}
		return c.JSON(response)
	}

	// Find matching employees using AI agent service
	fmt.Printf("DEBUG: Finding employees with skills: %v\n", skills)
	matches, err := h.aiAgentService.FindMatchingEmployees(skills)
	if err != nil {
		fmt.Printf("ERROR: Failed to find matching employees: %v\n", err)
		response := MatchEmployeesResponse{
			Matches:        []models.AIAgentMatch{},
			Summary:        "Error finding matching employees",
			ProcessingTime: time.Since(startTime).Milliseconds(),
			Status:         "failed",
			Message:        fmt.Sprintf("Failed to find matching employees: %v", err),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	fmt.Printf("DEBUG: Found %d matching employees\n", len(matches))

	// Save matches to database
	var savedMatches []models.Match
	for _, match := range matches {
		savedMatch, saveErr := h.aiAgentService.SaveMatch(&match)
		if saveErr != nil {
			fmt.Printf("Warning: Failed to save match for employee %d: %v\n", match.EmployeeID, saveErr)
		} else {
			savedMatches = append(savedMatches, *savedMatch)
		}
	}

	// Generate match explanations using AI agent service
	aiMatches := h.aiAgentService.GenerateMatchExplanations(matches, skills)
	matchSummary := h.aiAgentService.GenerateMatchSummary(aiMatches, skills)

	// Create AI agent request for tracking
	aiAgentRequest := &models.AIAgentRequest{
		Status:    "completed",
		CreatedAt: startTime,
	}

	// Save AI agent response to database
	response := &models.AIAgentResponse{
		RequestID:      aiAgentRequest.ID,
		Matches:        aiMatches,
		Summary:        matchSummary,
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Status:         "completed",
	}

	if err := h.aiAgentService.SaveResponse(response); err != nil {
		fmt.Printf("Warning: Failed to save AI agent response: %v\n", err)
	}

	fmt.Printf("Successfully completed employee matching with %d matches\n", len(matches))

	// Return response
	matchResponse := MatchEmployeesResponse{
		Matches:        aiMatches,
		Summary:        matchSummary,
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Status:         "completed",
		Message:        fmt.Sprintf("Found %d matching employees", len(matches)),
	}

	return c.JSON(matchResponse)
}

// GetMatchingHistory handles retrieving matching history
func (h *MatchingHandler) GetMatchingHistory(c *fiber.Ctx) error {
	// TODO: Implement history retrieval
	response := map[string]interface{}{
		"message": "Matching history endpoint - to be implemented",
		"status":  "pending",
	}

	return c.JSON(response)
}
