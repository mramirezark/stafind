package handlers

import (
	"encoding/json"
	"fmt"
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
	cvExtractService        services.CVExtractService
}

// NewExtractHandlers creates new extraction handlers
func NewExtractHandlers(extractionService *services.CandidateExtractService, aiAgentService services.AIAgentService, candidateStorageService *services.CandidateStorageService, cvExtractService services.CVExtractService) *ExtractHandlers {
	return &ExtractHandlers{
		extractionService:       extractionService,
		aiAgentService:          aiAgentService,
		candidateStorageService: candidateStorageService,
		cvExtractService:        cvExtractService,
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

	// Convert ExtractAIRequest to ExtractProcessRequest
	processRequest := models.ExtractProcessRequest{
		Text:             request.Text,
		ResumeURL:        request.FileURL, // Use FileURL as ResumeURL
		Metadata:         request.Metadata,
		ExtractionSource: "resume", // Default extraction source
		ProcessingType:   request.ProcessingType,
	}

	// Process the request using pure NER
	result, err := h.extractionService.ProcessText(&processRequest)
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

	// Convert ExtractAIRequest to ExtractProcessRequest
	processRequest := models.ExtractProcessRequest{
		Text:             request.Text,
		ResumeURL:        request.FileURL, // Use FileURL as ResumeURL
		Metadata:         request.Metadata,
		ExtractionSource: "resume", // Default extraction source
		ProcessingType:   request.ProcessingType,
	}

	result, err := h.extractionService.ProcessText(&processRequest)
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

	// Convert ExtractAIRequest to ExtractProcessRequest
	processRequest := models.ExtractProcessRequest{
		Text:             request.Text,
		ResumeURL:        request.FileURL, // Use FileURL as ResumeURL
		Metadata:         request.Metadata,
		ExtractionSource: "resume", // Default extraction source
		ProcessingType:   request.ProcessingType,
	}

	result, err := h.extractionService.ProcessText(&processRequest)
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

	// Convert ExtractAIRequest to ExtractProcessRequest
	processRequest := models.ExtractProcessRequest{
		Text:             extractionRequest.Text,
		ResumeURL:        extractionRequest.FileURL, // Use FileURL as ResumeURL
		Metadata:         extractionRequest.Metadata,
		ExtractionSource: "resume", // Default extraction source
		ProcessingType:   extractionRequest.ProcessingType,
	}

	result, err := h.extractionService.ProcessText(&processRequest)
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

// ExtractProcess handles both AI agent creation/retrieval and text extraction in one endpoint
func (h *ExtractHandlers) ExtractProcess(c *fiber.Ctx) error {
	var request models.ExtractProcessRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if request.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "text field is required",
		})
	}

	// CV Extract Tracking: Create or update extract record
	var cvExtract *models.CVExtract
	if request.ExtractRequestId != "" {
		// Create metadata for CV extract
		metadata := map[string]interface{}{
			"extraction_source": request.ExtractionSource,
			"processing_type":   request.ProcessingType,
			"resume_url":        request.ResumeURL,
			"file_number":       request.FileNumber,
			"total_files":       request.TotalFiles,
		}
		metadataJSON, _ := json.Marshal(metadata)
		metadataStr := string(metadataJSON)

		// Create or update CV extract record
		extract, err := h.cvExtractService.CreateOrUpdateExtract(
			request.ExtractRequestId,
			models.CVExtractStatusProcessing,
			request.TotalFiles,
			request.FileNumber,
			&metadataStr,
		)
		if err != nil {
			// Log error but don't fail the extraction process
			fmt.Printf("Warning: Failed to create/update CV extract record: %v\n", err)
		} else {
			cvExtract = extract
		}
	}

	result, err := h.extractionService.ProcessText(&request)
	if err != nil {
		// CV Extract Tracking: Mark as failed if extraction fails
		if cvExtract != nil && request.ExtractRequestId != "" {
			_, markErr := h.cvExtractService.MarkExtractFailed(
				request.ExtractRequestId,
				"Failed to process text with NER extraction: "+err.Error(),
			)
			if markErr != nil {
				fmt.Printf("Warning: Failed to mark CV extract as failed: %v\n", markErr)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process text with NER extraction",
			"details": err.Error(),
		})
	}

	var extractedData map[string]interface{}
	if err := json.Unmarshal([]byte(result.ProcessedContent), &extractedData); err != nil {
		// CV Extract Tracking: Mark as failed if parsing fails
		if cvExtract != nil && request.ExtractRequestId != "" {
			_, markErr := h.cvExtractService.MarkExtractFailed(
				request.ExtractRequestId,
				"Failed to parse extracted data: "+err.Error(),
			)
			if markErr != nil {
				fmt.Printf("Warning: Failed to mark CV extract as failed: %v\n", markErr)
			}
		}

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
		request.ResumeURL,
	)
	if err != nil {
		// CV Extract Tracking: Mark as failed if candidate extraction fails
		if cvExtract != nil && request.ExtractRequestId != "" {
			_, markErr := h.cvExtractService.MarkExtractFailed(
				request.ExtractRequestId,
				"Failed to process candidate extraction: "+err.Error(),
			)
			if markErr != nil {
				fmt.Printf("Warning: Failed to mark CV extract as failed: %v\n", markErr)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process candidate extraction",
			"details": err.Error(),
		})
	}

	//h.aiAgentService.UpdateAIAgentStatus(aiAgentRequest.ID, aiAgentRequest.Status)

	// CV Extract Tracking: Update completion status based on results
	if cvExtract != nil && request.ExtractRequestId != "" {
		// Calculate total processing time
		totalTimeMs := int64(result.ProcessingTime + candidateResult.ProcessingTime)

		// Update file progress
		_, err := h.cvExtractService.UpdateFileProgress(
			request.ExtractRequestId,
			request.FileNumber,
			1, // files processed (current file)
			0, // files failed
		)
		if err != nil {
			fmt.Printf("Warning: Failed to update CV extract progress: %v\n", err)
		}

		// Mark as successful completion
		_, err = h.cvExtractService.MarkExtractSuccess(
			request.ExtractRequestId,
			totalTimeMs,
		)
		if err != nil {
			// Log error but don't fail the response
			fmt.Printf("Warning: Failed to mark CV extract as successful: %v\n", err)
		}
	}

	// Return simplified processing response
	response := fiber.Map{
		"success":     true,
		"request_id":  request.ExtractRequestId,
		"file_number": request.FileNumber,
		"total_files": request.TotalFiles,
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
		"message": "File processing completed successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
