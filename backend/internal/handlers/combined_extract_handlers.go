package handlers

import (
	"encoding/json"
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CombinedExtractHandlers handles combined NER and Hugging Face extraction endpoints
type CombinedExtractHandlers struct {
	extractionService       *services.CandidateExtractService
	aiAgentService          services.AIAgentService
	candidateStorageService *services.CandidateStorageService
	cvExtractService        services.CVExtractService
	huggingFaceService      services.HuggingFaceSkillExtractionService
}

// NewCombinedExtractHandlers creates new combined extraction handlers
func NewCombinedExtractHandlers(
	extractionService *services.CandidateExtractService,
	aiAgentService services.AIAgentService,
	candidateStorageService *services.CandidateStorageService,
	cvExtractService services.CVExtractService,
	huggingFaceService services.HuggingFaceSkillExtractionService,
) *CombinedExtractHandlers {
	return &CombinedExtractHandlers{
		extractionService:       extractionService,
		aiAgentService:          aiAgentService,
		candidateStorageService: candidateStorageService,
		cvExtractService:        cvExtractService,
		huggingFaceService:      huggingFaceService,
	}
}

// ExtractProcessCombined handles both NER and Hugging Face extraction in one endpoint
func (h *CombinedExtractHandlers) ExtractProcessCombined(c *fiber.Ctx) error {
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

	// Start timing
	startTime := time.Now()

	// CV Extract Tracking: Create or update extract record
	var cvExtract *models.CVExtract
	if request.ExtractRequestId != "" {
		// Create metadata for CV extract
		metadata := map[string]interface{}{
			"extraction_source":  request.ExtractionSource,
			"processing_type":    request.ProcessingType,
			"resume_url":         request.ResumeURL,
			"file_number":        request.FileNumber,
			"total_files":        request.TotalFiles,
			"extraction_methods": []string{"ner", "huggingface"},
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
			fmt.Printf("Warning: Failed to create/update CV extract record: %v\n", err)
		} else {
			cvExtract = extract
		}
	}

	// Run both extractions in parallel
	nerResultChan := make(chan *models.ExtractAIResponse, 1)
	huggingFaceResultChan := make(chan *models.HuggingFaceSkillExtractionResponse, 1)
	errorChan := make(chan error, 2)

	// NER Extraction
	go func() {
		nerResult, err := h.extractionService.ProcessText(&request)
		if err != nil {
			errorChan <- fmt.Errorf("NER extraction failed: %w", err)
			return
		}
		nerResultChan <- nerResult
	}()

	// Hugging Face Extraction
	go func() {
		huggingFaceResult, err := h.huggingFaceService.ExtractSkillsFromText(request.Text)
		if err != nil {
			errorChan <- fmt.Errorf("Hugging Face extraction failed: %w", err)
			return
		}
		huggingFaceResultChan <- huggingFaceResult
	}()

	// Wait for both extractions to complete
	var nerResult *models.ExtractAIResponse
	var huggingFaceResult *models.HuggingFaceSkillExtractionResponse
	var extractionErrors []string

	// Collect results with timeout
	timeout := time.After(30 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case result := <-nerResultChan:
			nerResult = result
		case result := <-huggingFaceResultChan:
			huggingFaceResult = result
		case err := <-errorChan:
			extractionErrors = append(extractionErrors, err.Error())
		case <-timeout:
			extractionErrors = append(extractionErrors, "extraction timeout")
		}
	}

	// Check if we have at least one successful extraction
	if nerResult == nil && huggingFaceResult == nil {
		// CV Extract Tracking: Mark as failed if both extractions fail
		if cvExtract != nil && request.ExtractRequestId != "" {
			_, markErr := h.cvExtractService.MarkExtractFailed(
				request.ExtractRequestId,
				"Both NER and Hugging Face extractions failed: "+fmt.Sprintf("%v", extractionErrors),
			)
			if markErr != nil {
				fmt.Printf("Warning: Failed to mark CV extract as failed: %v\n", markErr)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Both extraction methods failed",
			"details": extractionErrors,
		})
	}

	// Parse NER result if available
	var nerExtractedData map[string]interface{}
	if nerResult != nil {
		if err := json.Unmarshal([]byte(nerResult.ProcessedContent), &nerExtractedData); err != nil {
			extractionErrors = append(extractionErrors, "Failed to parse NER result: "+err.Error())
		}
	}

	// Combine results
	combinedResult := h.combineExtractionResults(nerExtractedData, huggingFaceResult, extractionErrors)

	// Process candidate extraction and store in employees table
	extractionSource := request.ExtractionSource
	if extractionSource == "" {
		extractionSource = "resume"
	}

	candidateResult, err := h.candidateStorageService.ProcessCandidateExtraction(
		request.Text,
		combinedResult,
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

	// CV Extract Tracking: Update completion status based on results
	if cvExtract != nil && request.ExtractRequestId != "" {
		// Calculate total processing time
		totalTimeMs := int64(time.Since(startTime).Milliseconds())

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
			fmt.Printf("Warning: Failed to mark CV extract as successful: %v\n", err)
		}
	}

	// Return combined processing response
	response := fiber.Map{
		"success":     true,
		"request_id":  request.ExtractRequestId,
		"file_number": request.FileNumber,
		"total_files": request.TotalFiles,
		"extraction_methods": fiber.Map{
			"ner": fiber.Map{
				"success": nerResult != nil,
				"error": func() string {
					if nerResult == nil {
						for _, err := range extractionErrors {
							if fmt.Sprintf("%v", err) == "NER extraction failed" {
								return err
							}
						}
					}
					return ""
				}(),
			},
			"huggingface": fiber.Map{
				"success": huggingFaceResult != nil,
				"error": func() string {
					if huggingFaceResult == nil {
						for _, err := range extractionErrors {
							if fmt.Sprintf("%v", err) == "Hugging Face extraction failed" {
								return err
							}
						}
					}
					return ""
				}(),
			},
		},
		"extraction_result": fiber.Map{
			"combined_data":    combinedResult,
			"ner_data":         nerExtractedData,
			"huggingface_data": huggingFaceResult,
			"processing_time":  time.Since(startTime).Milliseconds(),
			"processing_type":  request.ProcessingType,
			"metadata":         request.Metadata,
			"timestamp":        time.Now(),
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
		"message": "Combined extraction processing completed successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// combineExtractionResults combines NER and Hugging Face extraction results
func (h *CombinedExtractHandlers) combineExtractionResults(
	nerData map[string]interface{},
	huggingFaceResult *models.HuggingFaceSkillExtractionResponse,
	errors []string,
) map[string]interface{} {
	combined := make(map[string]interface{})

	// Start with NER data as base
	for key, value := range nerData {
		combined[key] = value
	}

	// Add Hugging Face skills if available
	if huggingFaceResult != nil && huggingFaceResult.Success {
		// Convert Hugging Face skills to the expected format
		huggingFaceSkills := make(map[string][]string)
		for _, skill := range huggingFaceResult.Skills {
			category := skill.Category
			if category == "" {
				category = "Other"
			}
			if huggingFaceSkills[category] == nil {
				huggingFaceSkills[category] = []string{}
			}
			huggingFaceSkills[category] = append(huggingFaceSkills[category], skill.Name)
		}

		// Add Hugging Face skills to combined result
		combined["huggingface_skills"] = huggingFaceSkills
		combined["huggingface_confidence"] = huggingFaceResult.ConfidenceScore
		combined["huggingface_total_skills"] = huggingFaceResult.TotalSkills
		combined["huggingface_model_used"] = huggingFaceResult.ModelUsed

		// Merge skills with existing NER skills
		if existingSkills, ok := combined["skills"].(map[string]interface{}); ok {
			// Merge Hugging Face skills with existing skills
			for category, skills := range huggingFaceSkills {
				if existingSkills[category] == nil {
					existingSkills[category] = []string{}
				}
				// Add new skills that don't already exist
				existingSkillList := existingSkills[category].([]string)
				for _, newSkill := range skills {
					found := false
					for _, existingSkill := range existingSkillList {
						if existingSkill == newSkill {
							found = true
							break
						}
					}
					if !found {
						existingSkillList = append(existingSkillList, newSkill)
					}
				}
				existingSkills[category] = existingSkillList
			}
		} else {
			// No existing skills, use Hugging Face skills as base
			combined["skills"] = huggingFaceSkills
		}
	}

	// Add extraction metadata
	combined["extraction_errors"] = errors
	combined["extraction_timestamp"] = time.Now()
	combined["extraction_methods_used"] = []string{"ner", "huggingface"}

	return combined
}

// CompareExtractionMethods compares NER vs Hugging Face extraction results
func (h *CombinedExtractHandlers) CompareExtractionMethods(c *fiber.Ctx) error {
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

	startTime := time.Now()

	// Run both extractions
	nerResult, nerErr := h.extractionService.ProcessText(&request)
	huggingFaceResult, hfErr := h.huggingFaceService.ExtractSkillsFromText(request.Text)

	// Parse NER result
	var nerData map[string]interface{}
	if nerResult != nil {
		json.Unmarshal([]byte(nerResult.ProcessedContent), &nerData)
	}

	// Create comparison response
	comparison := fiber.Map{
		"success":         true,
		"text":            request.Text,
		"processing_time": time.Since(startTime).Milliseconds(),
		"methods": fiber.Map{
			"ner": fiber.Map{
				"success": nerErr == nil,
				"error": func() string {
					if nerErr != nil {
						return nerErr.Error()
					}
					return ""
				}(),
				"data": nerData,
				"processing_time": func() int64 {
					if nerResult != nil {
						return int64(nerResult.ProcessingTime)
					}
					return 0
				}(),
			},
			"huggingface": fiber.Map{
				"success": hfErr == nil,
				"error": func() string {
					if hfErr != nil {
						return hfErr.Error()
					}
					return ""
				}(),
				"data": huggingFaceResult,
				"processing_time": func() int64 {
					if huggingFaceResult != nil {
						return huggingFaceResult.ProcessingTime.Milliseconds()
					}
					return 0
				}(),
			},
		},
		"comparison": fiber.Map{
			"ner_skills_count": func() int {
				if skills, ok := nerData["skills"].(map[string]interface{}); ok {
					total := 0
					for _, skillList := range skills {
						if list, ok := skillList.([]string); ok {
							total += len(list)
						}
					}
					return total
				}
				return 0
			}(),
			"huggingface_skills_count": func() int {
				if huggingFaceResult != nil {
					return huggingFaceResult.TotalSkills
				}
				return 0
			}(),
			"ner_confidence": func() float64 {
				if confidence, ok := nerData["confidence_score"].(float64); ok {
					return confidence
				}
				return 0.0
			}(),
			"huggingface_confidence": func() float64 {
				if huggingFaceResult != nil {
					return huggingFaceResult.ConfidenceScore
				}
				return 0.0
			}(),
		},
	}

	return c.Status(fiber.StatusOK).JSON(comparison)
}
