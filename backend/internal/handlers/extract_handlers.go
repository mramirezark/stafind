package handlers

import (
	"encoding/json"
	"fmt"
	"stafind-backend/internal/logger"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ExtractHandlers handles candidate and resume extraction endpoints
type ExtractHandlers struct {
	extractionService       *services.CandidateExtractService
	aiAgentService          services.AIAgentService
	candidateStorageService *services.CandidateStorageService
}

// NewExtractHandlers creates new extraction handlers
func NewExtractHandlers(extractionService *services.CandidateExtractService, aiAgentService services.AIAgentService, candidateStorageService *services.CandidateStorageService) *ExtractHandlers {
	return &ExtractHandlers{
		extractionService:       extractionService,
		aiAgentService:          aiAgentService,
		candidateStorageService: candidateStorageService,
	}
}

// ProcessExtractedText handles extracted text using NER
func (h *ExtractHandlers) ProcessExtractedText(c *fiber.Ctx) error {
	var request models.ExtractAIRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if request.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "text field is required",
		})
	}

	if request.MessageText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "message_text field is required",
		})
	}

	// Process the request using pure NER
	result, err := h.extractionService.ProcessText(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process text with NER extraction",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// ExtractCandidateInfo processes resume text to extract structured candidate information
func (h *ExtractHandlers) ExtractCandidateInfo(c *fiber.Ctx) error {
	var request models.ExtractAIRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Set processing type for candidate extraction
	request.ProcessingType = "candidate_extraction"
	request.Metadata = map[string]interface{}{
		"language": "auto-detect",
		"format":   "resume",
	}

	result, err := h.extractionService.ProcessText(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract candidate information",
			"details": err.Error(),
		})
	}

	// Parse the structured result
	var candidateInfo map[string]interface{}
	if err := json.Unmarshal([]byte(result.ProcessedContent), &candidateInfo); err != nil {
		// If JSON parsing fails, return the raw result with a warning
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":           true,
			"candidate_info":    result.ProcessedContent,
			"warning":           "Could not parse structured data, returning raw result",
			"processing_time":   result.ProcessingTime,
			"model_used":        result.ModelUsed,
			"extraction_method": "Pure NER (Prose)",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":           true,
		"candidate_info":    candidateInfo,
		"processing_time":   result.ProcessingTime,
		"model_used":        result.ModelUsed,
		"extraction_method": "Pure NER (Prose)",
	})
}

// AnalyzeSearchRequest processes search requests to extract structured criteria using NER
func (h *ExtractHandlers) AnalyzeSearchRequest(c *fiber.Ctx) error {
	var request models.ExtractAIRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Set processing type for search request analysis
	request.ProcessingType = "search_analysis"
	request.Metadata = map[string]interface{}{
		"language": "auto-detect",
		"format":   "search_request",
	}

	result, err := h.extractionService.ProcessText(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to analyze search request",
			"details": err.Error(),
		})
	}

	// Parse the structured result
	var searchCriteria map[string]interface{}
	if err := json.Unmarshal([]byte(result.ProcessedContent), &searchCriteria); err != nil {
		// If JSON parsing fails, return the raw result with a warning
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":           true,
			"search_criteria":   result.ProcessedContent,
			"warning":           "Could not parse structured data, returning raw result",
			"processing_time":   result.ProcessingTime,
			"model_used":        result.ModelUsed,
			"extraction_method": "Pure NER (Prose)",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":           true,
		"search_criteria":   searchCriteria,
		"processing_time":   result.ProcessingTime,
		"model_used":        result.ModelUsed,
		"extraction_method": "Pure NER (Prose)",
	})
}

// MatchCandidateWithRequest performs candidate matching using NER analysis
func (h *ExtractHandlers) MatchCandidateWithRequest(c *fiber.Ctx) error {
	var request models.CandidateMatchingRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if request.CandidateInfo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "candidate_info field is required",
		})
	}

	if request.SearchCriteria == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search_criteria field is required",
		})
	}

	// Create extraction request for matching
	extractionRequest := models.ExtractAIRequest{
		Text: fmt.Sprintf("Candidate: %s\n\nSearch Criteria: %s",
			request.CandidateInfo, request.SearchCriteria),
		ProcessingType: "candidate_matching",
		Metadata: map[string]interface{}{
			"language":       "auto-detect",
			"format":         "matching",
			"candidate_name": request.CandidateName,
		},
	}

	result, err := h.extractionService.ProcessText(&extractionRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to match candidate with request",
			"details": err.Error(),
		})
	}

	// Parse the matching result
	var matchResult map[string]interface{}
	if err := json.Unmarshal([]byte(result.ProcessedContent), &matchResult); err != nil {
		// If JSON parsing fails, return the raw result with a warning
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":           true,
			"match_result":      result.ProcessedContent,
			"warning":           "Could not parse structured data, returning raw result",
			"processing_time":   result.ProcessingTime,
			"model_used":        result.ModelUsed,
			"extraction_method": "Pure NER (Prose)",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":           true,
		"match_result":      matchResult,
		"processing_time":   result.ProcessingTime,
		"model_used":        result.ModelUsed,
		"extraction_method": "Pure NER (Prose)",
	})
}

// HealthCheck provides a health check endpoint for the extraction service
func (h *ExtractHandlers) HealthCheck(c *fiber.Ctx) error {
	status, err := h.extractionService.GetHealthStatus()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "unhealthy",
			"error":     err.Error(),
			"timestamp": time.Now(),
			"service":   "Candidate Extraction Service (NER)",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    status,
		"timestamp": time.Now(),
		"service":   "Candidate Extraction Service (NER)",
		"method":    "Pure NER (Prose)",
	})
}

// ConsolidatedExtract handles both AI agent creation/retrieval and text extraction in one endpoint
func (h *ExtractHandlers) ConsolidatedExtract(c *fiber.Ctx) error {
	var request models.ConsolidatedExtractRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if request.MessageText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "message_text field is required",
		})
	}

	if request.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "text field is required",
		})
	}

	// Step 1: Handle AI Agent Request (create if not exists, get if exists)
	var aiAgentRequest *models.AIAgentRequest
	var err error

	// Check if we have a Teams message ID to look up existing request
	if request.TeamsMessageID != "" {
		// Try to get existing AI agent request by Teams message ID
		aiAgentRequest, err = h.aiAgentService.GetAIAgentRequestByTeamsMessageID(request.TeamsMessageID)
		if err != nil {
			// If not found, create a new one
			createRequest := models.CreateAIAgentRequest{
				TeamsMessageID: request.TeamsMessageID,
				ChannelID:      request.ChannelID,
				UserID:         request.UserID,
				UserName:       request.UserName,
				MessageText:    request.MessageText,
				AttachmentURL:  request.AttachmentURL,
			}

			// Generate defaults for public endpoint if not provided
			if createRequest.ChannelID == "" {
				createRequest.ChannelID = "consolidated-channel"
			}
			if createRequest.UserID == "" {
				createRequest.UserID = "consolidated-user"
			}
			if createRequest.UserName == "" {
				createRequest.UserName = "Consolidated User"
			}

			aiAgentRequest, err = h.aiAgentService.CreateAIAgentRequest(&createRequest)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Failed to create AI agent request",
					"details": err.Error(),
				})
			}
		}
	} else {
		// Create new AI agent request with generated Teams message ID
		createRequest := models.CreateAIAgentRequest{
			TeamsMessageID: fmt.Sprintf("consolidated-%d", time.Now().UnixNano()),
			ChannelID:      request.ChannelID,
			UserID:         request.UserID,
			UserName:       request.UserName,
			MessageText:    request.MessageText,
			AttachmentURL:  request.AttachmentURL,
		}

		// Generate defaults for public endpoint if not provided
		if createRequest.ChannelID == "" {
			createRequest.ChannelID = "consolidated-channel"
		}
		if createRequest.UserID == "" {
			createRequest.UserID = "consolidated-user"
		}
		if createRequest.UserName == "" {
			createRequest.UserName = "Consolidated User"
		}

		aiAgentRequest, err = h.aiAgentService.CreateAIAgentRequest(&createRequest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to create AI agent request",
				"details": err.Error(),
			})
		}
	}

	// Step 2: Process text extraction using NER
	extractRequest := models.ExtractAIRequest{
		TeamsMessageID: aiAgentRequest.TeamsMessageID,
		ChannelID:      aiAgentRequest.ChannelID,
		UserID:         aiAgentRequest.UserID,
		UserName:       aiAgentRequest.UserName,
		MessageText:    request.MessageText,
		Text:           request.Text,
		FileName:       request.FileName,
		FileURL:        request.FileURL,
		ProcessingType: request.ProcessingType,
		Metadata:       request.Metadata,
	}

	// Set default processing type if not provided
	if extractRequest.ProcessingType == "" {
		extractRequest.ProcessingType = "candidate_extraction"
	}

	// Set default metadata if not provided
	if extractRequest.Metadata == nil {
		extractRequest.Metadata = map[string]interface{}{
			"language": "auto-detect",
			"format":   "consolidated",
		}
	}

	result, err := h.extractionService.ProcessText(&extractRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process text with NER extraction",
			"details": err.Error(),
		})
	}

	// Step 3: Parse extracted data and store/update employee
	var extractedData map[string]interface{}
	if err := json.Unmarshal([]byte(result.ProcessedContent), &extractedData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to parse extracted data",
			"details": err.Error(),
		})
	}

	// Process candidate extraction and store in employees table
	extractionSource := request.ExtractionSource
	if extractionSource == "" {
		extractionSource = "resume"
	}

	candidateResult, err := h.candidateStorageService.ProcessCandidateExtraction(
		request.Text,
		extractedData,
		extractionSource,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process candidate extraction",
			"details": err.Error(),
		})
	}

	// Step 4: Update AI agent request with extracted text
	aiAgentRequest.ExtractedText = &result.ProcessedContent
	if request.TotalFiles == request.FileNumber {
		aiAgentRequest.Status = "completed"
	} else {
		aiAgentRequest.Status = "processing"
	}
	now := time.Now()
	aiAgentRequest.ProcessedAt = &now

	// Step 5: If all files are processed, trigger matching process using AI agent service
	var aiMatches []models.AIAgentMatch
	var matchSummary string
	var savedMatches []models.Match
	//Temp fix for testing
	if request.TotalFiles > 0 {
		//if request.TotalFiles > 0 && request.TotalFiles == request.FileNumber {

		// All files processed, extract skills from the search request message text
		skillExtractResult, err := h.aiAgentService.ExtractSkillsFromText(request.MessageText)
		if err != nil {
			aiAgentRequest.Status = "failed"
			fmt.Printf("Warning: Failed to extract skills from search request: %v\n", err)
		} else {
			// Find matching employees using AI agent service
			matches, err := h.aiAgentService.FindMatchingEmployees(skillExtractResult.Skills)
			if err != nil {
				aiAgentRequest.Status = "failed"
				fmt.Printf("Warning: Failed to find matching employees: %v\n", err)
			} else {
				// Save matches to database
				for _, match := range matches {
					savedMatch, err := h.aiAgentService.SaveMatch(&match)
					if err != nil {
						fmt.Printf("Warning: Failed to save match for employee %d: %v\n", match.EmployeeID, err)
					} else {
						savedMatches = append(savedMatches, *savedMatch)
					}
				}

				// Generate match explanations using AI agent service
				aiMatches = h.aiAgentService.GenerateMatchExplanations(matches, skillExtractResult.Skills)
				matchSummary = h.aiAgentService.GenerateMatchSummary(aiMatches, skillExtractResult.Skills)

				// Save AI agent response to database
				response := &models.AIAgentResponse{
					RequestID:      aiAgentRequest.ID,
					Matches:        aiMatches,
					Summary:        matchSummary,
					ProcessingTime: time.Since(aiAgentRequest.CreatedAt).Milliseconds(),
					Status:         "completed",
				}

				if err := h.aiAgentService.SaveResponse(response); err != nil {
					fmt.Printf("Warning: Failed to save AI agent response: %v\n", err)
				}

				aiAgentRequest.Status = "completed"
			}
		}
	}
	h.aiAgentService.UpdateAIAgentStatus(aiAgentRequest.ID, aiAgentRequest.Status)

	// Return consolidated response
	response := fiber.Map{
		"success": true,
		"ai_agent_request": fiber.Map{
			"id":               aiAgentRequest.ID,
			"teams_message_id": aiAgentRequest.TeamsMessageID,
			"channel_id":       aiAgentRequest.ChannelID,
			"user_id":          aiAgentRequest.UserID,
			"user_name":        aiAgentRequest.UserName,
			"message_text":     aiAgentRequest.MessageText,
			"status":           aiAgentRequest.Status,
			"created_at":       aiAgentRequest.CreatedAt,
			"processed_at":     aiAgentRequest.ProcessedAt,
		},
		"extraction_result": fiber.Map{
			"processed_content": result.ProcessedContent,
			"processing_time":   result.ProcessingTime,
			"model_used":        result.ModelUsed,
			"processing_type":   result.ProcessingType,
			"metadata":          result.Metadata,
			"timestamp":         result.Timestamp,
		},
		"candidate_result": fiber.Map{
			"employee_id":      candidateResult.EmployeeID,
			"action":           candidateResult.Action,
			"changes_detected": candidateResult.ChangesDetected,
			"changes_summary":  candidateResult.ChangesSummary,
			"processing_time":  candidateResult.ProcessingTime,
			"status":           candidateResult.Status,
			"message":          candidateResult.Message,
		},
		"message": "Consolidated extraction completed successfully",
	}

	// Add AI agent matching result if available
	if len(aiMatches) > 0 {
		response["ai_matching_result"] = fiber.Map{
			"matches":        aiMatches,
			"summary":        matchSummary,
			"total_matches":  len(aiMatches),
			"saved_matches":  len(savedMatches),
			"response_saved": true,
		}
		response["message"] = "All files processed and AI matching completed successfully"
	}

	logger.Info("Status, ", aiAgentRequest.Status, "response", response)
	return c.Status(fiber.StatusOK).JSON(response)
}
