package handlers

import (
	"strconv"
	"strings"

	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// NERHandlers handles NER-related API endpoints
type NERHandlers struct {
	nerService    *services.NERService
	searchService services.SearchService
}

// NewNERHandlers creates a new NER handlers instance
func NewNERHandlers(nerService *services.NERService, searchService services.SearchService) *NERHandlers {
	return &NERHandlers{
		nerService:    nerService,
		searchService: searchService,
	}
}

// ExtractSkillsRequest represents the request payload for skill extraction
type ExtractSkillsRequest struct {
	Text string `json:"text" validate:"required"`
}

// ExtractSkills extracts skills from text using Go NER library
// @Summary Extract skills from text using Go NER library
// @Description Uses Prose NER library to extract programming languages, frameworks, databases, and other skills from job descriptions
// @Tags NER
// @Accept json
// @Produce json
// @Param request body ExtractSkillsRequest true "Text to analyze"
// @Success 200 {object} services.SkillExtractionResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ner/extract-skills [post]
func (h *NERHandlers) ExtractSkills(c *fiber.Ctx) error {
	var req ExtractSkillsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Text == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error": "Text field is required",
		})
	}

	// Extract skills using Go NER library
	result, err := h.nerService.ExtractSkillsFromText(req.Text)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract skills",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Skills extracted successfully using Go NER library",
		"data":    result,
	})
}

// ExtractSkillsFromJobDescription extracts skills from a job description with additional context
// @Summary Extract skills from job description
// @Description Extracts skills from job descriptions with enhanced context awareness
// @Tags NER
// @Accept json
// @Produce json
// @Param request body JobDescriptionRequest true "Job description to analyze"
// @Success 200 {object} services.SkillExtractionResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ner/extract-job-skills [post]
func (h *NERHandlers) ExtractSkillsFromJobDescription(c *fiber.Ctx) error {
	var req JobDescriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Description == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error": "Description field is required",
		})
	}

	// Combine title and description for better context
	fullText := req.Title
	if req.Title != "" && req.Description != "" {
		fullText += "\n\n" + req.Description
	} else if req.Description != "" {
		fullText = req.Description
	}

	// Extract skills using Go NER library
	result, err := h.nerService.ExtractSkillsFromText(fullText)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract skills from job description",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Skills extracted from job description using Go NER library",
		"data":    result,
		"job_info": fiber.Map{
			"title":        req.Title,
			"company":      req.Company,
			"location":     req.Location,
			"salary_range": req.SalaryRange,
		},
	})
}

// JobDescriptionRequest represents a job description analysis request
type JobDescriptionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description" validate:"required"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	SalaryRange string `json:"salary_range"`
}

// ExtractSkillsAndSearch extracts skills and finds matching employees in one call
// @Summary Extract skills and search for matching employees
// @Description Extracts skills from job description and immediately searches for matching employees
// @Tags NER
// @Accept json
// @Produce json
// @Param request body ExtractAndSearchRequest true "Job description to analyze and search"
// @Success 200 {object} ExtractAndSearchResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ner/extract-and-search [post]
func (h *NERHandlers) ExtractSkillsAndSearch(c *fiber.Ctx) error {
	var req ExtractAndSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Text == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error": "Text field is required",
		})
	}

	// Extract skills using Go NER library
	nerResult, err := h.nerService.ExtractSkillsFromText(req.Text)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract skills",
			"details": err.Error(),
		})
	}

	// Convert extracted skills to search request format
	var allSkills []string
	for _, skillList := range nerResult.Skills.Categories {
		allSkills = append(allSkills, skillList...)
	}

	searchRequest := models.SearchRequest{
		RequiredSkills:  allSkills,
		PreferredSkills: []string{},
		Department:      "",
		ExperienceLevel: func() string {
			if len(nerResult.Skills.YearsOfExperience) > 0 {
				// Get the maximum years from the array
				maxYears := 0
				for _, yearStr := range nerResult.Skills.YearsOfExperience {
					// Extract numeric value from string
					cleanStr := ""
					for _, char := range yearStr {
						if char >= '0' && char <= '9' {
							cleanStr += string(char)
						}
					}
					if cleanStr != "" {
						if years, err := strconv.Atoi(cleanStr); err == nil && years > maxYears {
							maxYears = years
						}
					}
				}

				if maxYears >= 5 {
					return "senior"
				} else if maxYears >= 3 {
					return "mid"
				} else {
					return "junior"
				}
			}
			return ""
		}(),
		Location:      req.Location,
		MinMatchScore: 0.3,
	}

	// Search for matching employees
	matches, err := h.searchService.SearchEmployees(&searchRequest)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to search for matching employees",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Skills extracted and search completed",
		"data": ExtractAndSearchResult{
			ExtractedSkills:   *nerResult,
			MatchingEmployees: matches,
			SearchCriteria:    searchRequest,
			TotalMatches:      len(matches),
			ProcessingTime:    nerResult.ProcessingTime,
		},
	})
}

// ExtractAndSearchRequest represents a request to extract skills and search for employees
type ExtractAndSearchRequest struct {
	Text        string `json:"text" validate:"required"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	SalaryRange string `json:"salary_range"`
	Limit       int    `json:"limit"`
}

// ExtractAndSearchResult represents the result of skill extraction and employee search
type ExtractAndSearchResult struct {
	ExtractedSkills   services.SkillExtractionResult `json:"extracted_skills"`
	MatchingEmployees []models.Match                 `json:"matching_employees"`
	SearchCriteria    models.SearchRequest           `json:"search_criteria"`
	TotalMatches      int                            `json:"total_matches"`
	ProcessingTime    string                         `json:"processing_time"`
}

// CompareSkills compares two sets of extracted skills
// @Summary Compare two sets of skills
// @Description Compares skills extracted from two different texts (e.g., job description vs candidate resume)
// @Tags NER
// @Accept json
// @Produce json
// @Param request body SkillComparisonRequest true "Skills to compare"
// @Success 200 {object} SkillComparisonResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ner/compare-skills [post]
func (h *NERHandlers) CompareSkills(c *fiber.Ctx) error {
	var req SkillComparisonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Text1 == "" || req.Text2 == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error": "Both text1 and text2 fields are required",
		})
	}

	// Extract skills from both texts
	result1, err := h.nerService.ExtractSkillsFromText(req.Text1)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract skills from first text",
			"details": err.Error(),
		})
	}

	result2, err := h.nerService.ExtractSkillsFromText(req.Text2)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract skills from second text",
			"details": err.Error(),
		})
	}

	// Compare skills
	comparison := h.compareSkillSets(result1.Skills, result2.Skills)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Skills compared successfully",
		"data": SkillComparisonResult{
			Text1Skills: result1.Skills,
			Text2Skills: result2.Skills,
			Comparison:  comparison,
			MatchScore:  comparison.OverallMatchScore,
		},
	})
}

// SkillComparisonRequest represents a request to compare two sets of skills
type SkillComparisonRequest struct {
	Text1      string `json:"text1" validate:"required"`
	Text2      string `json:"text2" validate:"required"`
	Text1Label string `json:"text1_label"`
	Text2Label string `json:"text2_label"`
}

// SkillComparisonResult represents the result of skill comparison
type SkillComparisonResult struct {
	Text1Skills services.ExtractedSkills `json:"text1_skills"`
	Text2Skills services.ExtractedSkills `json:"text2_skills"`
	Comparison  SkillComparison          `json:"comparison"`
	MatchScore  float64                  `json:"match_score"`
}

// SkillComparison represents the detailed comparison between two skill sets
type SkillComparison struct {
	CommonSkills      []string `json:"common_skills"`
	Text1OnlySkills   []string `json:"text1_only_skills"`
	Text2OnlySkills   []string `json:"text2_only_skills"`
	ProgrammingMatch  float64  `json:"programming_match"`
	WebTechMatch      float64  `json:"web_tech_match"`
	DatabaseMatch     float64  `json:"database_match"`
	CloudDevOpsMatch  float64  `json:"cloud_devops_match"`
	SoftSkillsMatch   float64  `json:"soft_skills_match"`
	OverallMatchScore float64  `json:"overall_match_score"`
}

// compareSkillSets compares two sets of extracted skills using dynamic categories
func (h *NERHandlers) compareSkillSets(skills1, skills2 services.ExtractedSkills) SkillComparison {
	comparison := SkillComparison{}

	// Get all unique categories from both skill sets
	allCategories := make(map[string]bool)
	for category := range skills1.Categories {
		allCategories[category] = true
	}
	for category := range skills2.Categories {
		allCategories[category] = true
	}

	// Compare skills in each category
	totalMatch := 0.0
	categoryCount := 0

	for category := range allCategories {
		skills1List := skills1.Categories[category]
		skills2List := skills2.Categories[category]

		if len(skills1List) == 0 && len(skills2List) == 0 {
			continue
		}

		common, _, _ := h.compareStringSlices(skills1List, skills2List)
		total := len(skills1List) + len(skills2List) - len(common)

		if total > 0 {
			match := float64(len(common)) / float64(total)
			totalMatch += match
			categoryCount++
		}
	}

	// Calculate overall match percentage
	if categoryCount > 0 {
		comparison.OverallMatchScore = totalMatch / float64(categoryCount)
	}

	// Set individual category matches to overall match for simplicity
	comparison.ProgrammingMatch = comparison.OverallMatchScore
	comparison.WebTechMatch = comparison.OverallMatchScore
	comparison.DatabaseMatch = comparison.OverallMatchScore
	comparison.CloudDevOpsMatch = comparison.OverallMatchScore
	comparison.SoftSkillsMatch = comparison.OverallMatchScore

	// Combine all common skills from all categories
	var allCommonSkills []string
	var allText1OnlySkills []string
	var allText2OnlySkills []string

	for category := range allCategories {
		skills1List := skills1.Categories[category]
		skills2List := skills2.Categories[category]

		common, only1, only2 := h.compareStringSlices(skills1List, skills2List)
		allCommonSkills = append(allCommonSkills, common...)
		allText1OnlySkills = append(allText1OnlySkills, only1...)
		allText2OnlySkills = append(allText2OnlySkills, only2...)
	}

	comparison.CommonSkills = allCommonSkills
	comparison.Text1OnlySkills = allText1OnlySkills
	comparison.Text2OnlySkills = allText2OnlySkills

	// Overall match score is already calculated above

	return comparison
}

// compareStringSlices compares two string slices and returns common, only1, and only2 elements
func (h *NERHandlers) compareStringSlices(slice1, slice2 []string) (common, only1, only2 []string) {
	// Create maps for efficient lookup
	map1 := make(map[string]bool)
	map2 := make(map[string]bool)

	for _, item := range slice1 {
		map1[strings.ToLower(item)] = true
	}
	for _, item := range slice2 {
		map2[strings.ToLower(item)] = true
	}

	// Find common elements
	for _, item := range slice1 {
		if map2[strings.ToLower(item)] {
			common = append(common, item)
		} else {
			only1 = append(only1, item)
		}
	}

	// Find elements only in slice2
	for _, item := range slice2 {
		if !map1[strings.ToLower(item)] {
			only2 = append(only2, item)
		}
	}

	return common, only1, only2
}
