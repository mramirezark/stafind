package services

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"strconv"
	"strings"
	"time"
)

// CandidateExtractService handles intelligent extraction from resumes and job descriptions
type CandidateExtractService struct {
	nerService *NERService
}

// NewCandidateExtractionService creates a new extraction service using pure NER
func NewCandidateExtractionService(skillRepo repositories.SkillRepository, categoryRepo repositories.CategoryRepository) *CandidateExtractService {
	log.Printf("Initializing Candidate Extraction Service with pure NER (Prose library)")

	return &CandidateExtractService{
		nerService: NewNERService(skillRepo, categoryRepo),
	}
}

// ProcessText processes text and returns structured results using pure NER
func (s *CandidateExtractService) ProcessText(request *models.ExtractProcessRequest) (*models.ExtractAIResponse, error) {
	startTime := time.Now()

	// Process with pure NER
	processedContent, err := s.processWithNER(request.ProcessingType, request.Text, request.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to process text with NER: %w", err)
	}

	processingTime := time.Since(startTime)

	response := &models.ExtractAIResponse{
		ProcessedContent: processedContent,
		ProcessingTime:   processingTime,
		ModelUsed:        "NER (Prose)",
		TokensProcessed:  s.estimateTokens(request.Text + processedContent),
		ProcessingType:   request.ProcessingType,
		Metadata:         request.Metadata,
		Timestamp:        time.Now(),
	}

	return response, nil
}

// processWithNER processes text using pure NER extraction
func (s *CandidateExtractService) processWithNER(processingType, text string, metadata map[string]interface{}) (string, error) {
	log.Printf("Processing with pure NER (Prose library)")
	log.Printf("Processing type: %s, Text length: %d characters", processingType, len(text))

	// Use NER service for all extraction
	switch processingType {
	case "candidate_extraction":
		log.Printf("Extracting candidate information using NER...")
		return s.extractCandidateWithNER(text)

	case "search_analysis":
		log.Printf("Analyzing search request using NER...")
		return s.analyzeSearchWithNER(text)

	case "candidate_matching":
		log.Printf("Matching candidate with NER analysis...")
		return s.matchCandidateWithNER(text)

	default:
		log.Printf("Processing generic text with NER...")
		return s.processGenericWithNER(text)
	}
}

// extractCandidateWithNER uses pure NER to extract candidate information
func (s *CandidateExtractService) extractCandidateWithNER(text string) (string, error) {
	// Use NER service to extract all information
	nerResult, err := s.nerService.ExtractSkillsFromText(text)
	if err != nil {
		return "", fmt.Errorf("NER extraction failed: %w", err)
	}

	// Extract contact information using NER entities
	contactInfo := s.extractContactInfoWithNER(text)

	// Build complete candidate info using pure NER
	candidateInfo := map[string]interface{}{
		"candidate_name":   s.extractNameWithNER(text),
		"contact_info":     contactInfo,
		"skills":           nerResult.Skills.Categories,
		"years_experience": nerResult.Skills.YearsOfExperience,
		"seniority_level":  s.extractAndNormalizeSeniorityLevelNER(text, nerResult.Skills.YearsOfExperience),
		"current_position": s.extractCurrentPosition(text),
		"education": map[string]interface{}{
			"level":        nerResult.Skills.EducationLevel,
			"institutions": s.extractInstitutionsWithNER(text),
		},
		"languages":         nerResult.Skills.LanguagesDetected,
		"summary":           nerResult.Skills.Summary,
		"total_skills":      nerResult.TotalSkillsFound,
		"confidence":        nerResult.Skills.ConfidenceScore,
		"extraction_method": "Pure NER (Prose)",
		"processing_time":   nerResult.ProcessingTime,
	}

	jsonData, err := json.MarshalIndent(candidateInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal candidate info: %w", err)
	}

	log.Printf("Successfully extracted candidate info using pure NER (confidence: %.2f, skills: %d)",
		nerResult.Skills.ConfidenceScore, nerResult.TotalSkillsFound)
	return string(jsonData), nil
}

// analyzeSearchWithNER uses pure NER to analyze search requests
func (s *CandidateExtractService) analyzeSearchWithNER(text string) (string, error) {
	// Use NER service to extract required skills
	nerResult, err := s.nerService.ExtractSkillsFromText(text)
	if err != nil {
		return "", fmt.Errorf("NER extraction failed: %w", err)
	}

	// Determine experience level from years
	experienceLevel := "Any"
	if len(nerResult.Skills.YearsOfExperience) > 0 {
		// Get the highest years of experience found
		maxYears := s.getMaxYearsFromStrings(nerResult.Skills.YearsOfExperience)
		if maxYears >= 7 {
			experienceLevel = "Senior"
		} else if maxYears >= 3 {
			experienceLevel = "Mid-Level"
		} else if maxYears > 0 {
			experienceLevel = "Junior"
		}
	}

	// Build comprehensive search criteria using NER results
	searchCriteria := map[string]interface{}{
		"original_request":  text,
		"detected_language": strings.Join(nerResult.Skills.LanguagesDetected, ", "),
		"search_criteria": map[string]interface{}{
			"skills":               nerResult.Skills.Categories,
			"experience_level":     experienceLevel,
			"years_experience_min": nerResult.Skills.YearsOfExperience,
		},
		"total_skills_found":  nerResult.TotalSkillsFound,
		"confidence":          nerResult.Skills.ConfidenceScore,
		"response_suggestion": fmt.Sprintf("Found %d technical skills. Ready to search for matching candidates.", nerResult.TotalSkillsFound),
		"extraction_method":   "Pure NER (Prose)",
		"processing_time":     nerResult.ProcessingTime,
	}

	jsonData, err := json.MarshalIndent(searchCriteria, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal search criteria: %w", err)
	}

	log.Printf("Successfully analyzed search with pure NER (%d skills found)", nerResult.TotalSkillsFound)
	return string(jsonData), nil
}

// matchCandidateWithNER uses NER for candidate matching analysis
func (s *CandidateExtractService) matchCandidateWithNER(text string) (string, error) {
	// Extract skills from the combined text
	nerResult, err := s.nerService.ExtractSkillsFromText(text)
	if err != nil {
		return "", fmt.Errorf("NER matching failed: %w", err)
	}

	// Calculate match score based on skill overlap
	matchScore := s.calculateMatchScore(nerResult)

	matchResult := map[string]interface{}{
		"match_score":       matchScore,
		"match_percentage":  fmt.Sprintf("%d%%", matchScore),
		"extracted_skills":  nerResult.Skills.Categories,
		"total_skills":      nerResult.TotalSkillsFound,
		"years_experience":  nerResult.Skills.YearsOfExperience,
		"recommendation":    s.generateRecommendation(matchScore),
		"confidence":        nerResult.Skills.ConfidenceScore,
		"extraction_method": "Pure NER (Prose)",
		"processing_time":   nerResult.ProcessingTime,
	}

	jsonData, err := json.MarshalIndent(matchResult, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal match result: %w", err)
	}

	return string(jsonData), nil
}

// processGenericWithNER uses pure NER for generic text processing
func (s *CandidateExtractService) processGenericWithNER(text string) (string, error) {
	nerResult, err := s.nerService.ExtractSkillsFromText(text)
	if err != nil {
		return "", fmt.Errorf("NER processing failed: %w", err)
	}

	analysis := map[string]interface{}{
		"text_analysis": map[string]interface{}{
			"word_count":      len(strings.Fields(text)),
			"character_count": len(text),
			"languages":       nerResult.Skills.LanguagesDetected,
		},
		"extracted_skills":  nerResult.Skills.Categories,
		"years_experience":  nerResult.Skills.YearsOfExperience,
		"education_level":   nerResult.Skills.EducationLevel,
		"processing_method": "Pure NER (Prose)",
		"confidence_score":  nerResult.Skills.ConfidenceScore,
		"processing_time":   nerResult.ProcessingTime,
	}

	jsonData, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal analysis: %w", err)
	}

	return string(jsonData), nil
}

// Helper methods for NER-based extraction

func (s *CandidateExtractService) extractNameWithNER(text string) string {
	// Extract name from first line or using simple heuristics
	lines := strings.Split(text, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Check if first line looks like a name (2-4 words, proper case)
		words := strings.Fields(firstLine)
		if len(words) >= 2 && len(words) <= 4 {
			allProperCase := true
			for _, word := range words {
				if len(word) > 0 && !strings.HasPrefix(word, strings.ToUpper(string(word[0]))) {
					allProperCase = false
					break
				}
			}
			if allProperCase {
				return firstLine
			}
		}
	}

	// Try to extract name from the beginning of the text
	words := strings.Fields(text)
	if len(words) >= 2 {
		// Check if first two words look like a name
		firstName := words[0]
		lastName := words[1]
		if len(firstName) > 1 && len(lastName) > 1 &&
			strings.HasPrefix(firstName, strings.ToUpper(string(firstName[0]))) &&
			strings.HasPrefix(lastName, strings.ToUpper(string(lastName[0]))) {
			return firstName + " " + lastName
		}
	}

	return "Name not found"
}

func (s *CandidateExtractService) extractContactInfoWithNER(text string) map[string]string {
	contactInfo := map[string]string{
		"email":    s.findEmail(text),
		"phone":    s.findPhone(text),
		"location": s.findLocation(text),
	}
	return contactInfo
}

func (s *CandidateExtractService) findEmail(text string) string {
	// Preprocess text to handle split emails (common in PDFs and formatted text)
	// Replace line breaks that might split email addresses
	processedText := strings.ReplaceAll(text, "\n", " ")
	processedText = strings.ReplaceAll(processedText, "\r", " ")
	processedText = strings.ReplaceAll(processedText, "\t", " ")
	// Remove multiple spaces
	processedText = regexp.MustCompile(`\s+`).ReplaceAllString(processedText, " ")

	// Handle split email domains (e.g., "gmail.c om" -> "gmail.com")
	// Look for patterns like "domain.c om" and fix them
	domainFixRegex := regexp.MustCompile(`([a-zA-Z0-9.-]+)\.([a-zA-Z]{1,3})\s+([a-zA-Z]{2,})`)
	processedText = domainFixRegex.ReplaceAllString(processedText, "$1.$2$3")

	fmt.Printf("DEBUG: Processed text for email detection: %s\n", processedText)

	// Simple email detection using regex for better accuracy
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(processedText, -1)
	fmt.Printf("DEBUG: Email regex matches: %v\n", matches)
	if len(matches) > 0 {
		result := strings.Trim(matches[0], ",.;:()[]")
		fmt.Printf("DEBUG: Found email: %s\n", result)
		return result
	}

	// If no email found, generate a placeholder based on name
	lines := strings.Split(text, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		words := strings.Fields(firstLine)
		if len(words) >= 2 {
			// Generate email from name: firstname.lastname@example.com
			firstName := strings.ToLower(words[0])
			lastName := strings.ToLower(words[1])
			return fmt.Sprintf("%s.%s@example.com", firstName, lastName)
		}
	}

	// Fallback to a generic email
	return "candidate@example.com"
}

func (s *CandidateExtractService) findPhone(text string) string {
	// Look for phone patterns
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "phone") || strings.Contains(lower, "tel") || strings.Contains(lower, "móvil") {
			// Extract the number part
			words := strings.Fields(line)
			for _, word := range words {
				if strings.ContainsAny(word, "0123456789") && len(word) >= 9 {
					return word
				}
			}
		}
	}
	// Look for standalone numbers
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "+") && len(word) > 10 {
			return word
		}
	}
	return "Not provided"
}

func (s *CandidateExtractService) findLocation(text string) string {
	// Look for location indicators
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "location") || strings.Contains(lower, "ubicación") ||
			strings.Contains(lower, "ciudad") || strings.Contains(lower, "city") {
			return strings.TrimSpace(strings.Split(line, ":")[len(strings.Split(line, ":"))-1])
		}
	}
	return "Not specified"
}

func (s *CandidateExtractService) extractInstitutionsWithNER(text string) []string {
	institutions := []string{}
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "university") || strings.Contains(lower, "universidad") ||
			strings.Contains(lower, "college") || strings.Contains(lower, "instituto") {
			institutions = append(institutions, strings.TrimSpace(line))
		}
	}

	if len(institutions) == 0 {
		institutions = append(institutions, "Not specified")
	}
	return institutions
}

// getMaxYearsFromStrings extracts the maximum years from a string array
func (s *CandidateExtractService) getMaxYearsFromStrings(yearsStrings []string) int {
	maxYears := 0
	for _, yearStr := range yearsStrings {
		// Try to extract numeric value from string
		if years := s.extractNumericYears(yearStr); years > maxYears {
			maxYears = years
		}
	}
	return maxYears
}

// extractNumericYears extracts numeric years from a string
func (s *CandidateExtractService) extractNumericYears(yearStr string) int {
	// Remove non-numeric characters and try to parse
	cleanStr := ""
	for _, char := range yearStr {
		if char >= '0' && char <= '9' {
			cleanStr += string(char)
		}
	}

	if cleanStr != "" {
		// Try to parse as integer
		if years, err := strconv.Atoi(cleanStr); err == nil {
			return years
		}
	}

	return 0
}

// extractAndNormalizeSeniorityLevelNER extracts and normalizes seniority level using NER data
func (s *CandidateExtractService) extractAndNormalizeSeniorityLevelNER(text string, yearsExperience []string) string {
	// First try to extract from text directly using improved logic
	seniorityFromText := s.extractSeniorityFromTextNER(text)
	if seniorityFromText != "" {
		return seniorityFromText
	}

	// If not found in text, determine based on years of experience
	if len(yearsExperience) > 0 {
		maxYears := s.getMaxYearsFromStrings(yearsExperience)
		return s.normalizeSeniorityFromYears(maxYears)
	}

	// Default fallback
	return "Mid"
}

// extractSeniorityFromTextNER extracts seniority level from text using improved pattern matching
func (s *CandidateExtractService) extractSeniorityFromTextNER(text string) string {
	lines := strings.Split(text, "\n")

	// Look for common patterns in resume text
	patterns := []string{
		"senior", "sr", "lead", "principal", "staff",
		"mid", "middle", "intermediate",
		"junior", "jr", "entry", "associate",
		"intern", "trainee", "apprentice",
		"manager", "director", "vp", "executive",
		"architect", "expert",
	}

	for _, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			continue
		}

		for _, pattern := range patterns {
			if strings.Contains(line, pattern) {
				// Found a potential seniority indicator
				return s.normalizeSeniorityLevelNER(line)
			}
		}
	}

	return ""
}

// normalizeSeniorityLevelNER normalizes various seniority level formats to standard values
func (s *CandidateExtractService) normalizeSeniorityLevelNER(level string) string {
	level = strings.ToLower(strings.TrimSpace(level))

	// Handle various seniority level formats
	switch {
	case strings.Contains(level, "senior") || strings.Contains(level, "sr"):
		return "Senior"
	case strings.Contains(level, "lead") || strings.Contains(level, "principal") || strings.Contains(level, "staff"):
		return "Senior"
	case strings.Contains(level, "mid") || strings.Contains(level, "middle") || strings.Contains(level, "intermediate"):
		return "Mid"
	case strings.Contains(level, "junior") || strings.Contains(level, "jr") || strings.Contains(level, "entry") || strings.Contains(level, "associate"):
		return "Junior"
	case strings.Contains(level, "intern") || strings.Contains(level, "trainee") || strings.Contains(level, "apprentice"):
		return "Junior"
	case strings.Contains(level, "manager") || strings.Contains(level, "director") || strings.Contains(level, "vp") || strings.Contains(level, "executive"):
		return "Senior"
	case strings.Contains(level, "architect") || strings.Contains(level, "expert"):
		return "Senior"
	case strings.Contains(level, "level") && (strings.Contains(level, "3") || strings.Contains(level, "iii")):
		return "Senior"
	case strings.Contains(level, "level") && (strings.Contains(level, "2") || strings.Contains(level, "ii")):
		return "Mid"
	case strings.Contains(level, "level") && (strings.Contains(level, "1") || strings.Contains(level, "i")):
		return "Junior"
	default:
		return ""
	}
}

// normalizeSeniorityFromYears normalizes seniority based on years of experience
func (s *CandidateExtractService) normalizeSeniorityFromYears(years int) string {
	if years >= 7 {
		return "Senior"
	} else if years >= 3 {
		return "Mid"
	} else if years > 0 {
		return "Junior"
	}
	return "Mid" // Default fallback
}

// extractSeniorityLevelFromStrings determines seniority level from years strings
func (s *CandidateExtractService) extractSeniorityLevelFromStrings(text string, yearsExperience []string) string {
	if len(yearsExperience) == 0 {
		return "Any"
	}

	maxYears := s.getMaxYearsFromStrings(yearsExperience)
	if maxYears >= 7 {
		return "Senior"
	} else if maxYears >= 3 {
		return "Mid-Level"
	} else if maxYears > 0 {
		return "Junior"
	}

	return "Any"
}

func (s *CandidateExtractService) extractSeniorityLevel(text string, yearsExperience *int) string {
	// First try to extract from text directly using improved logic
	seniorityFromText := s.extractSeniorityFromTextNER(text)
	if seniorityFromText != "" {
		return seniorityFromText
	}

	// If not found in text, determine based on years of experience
	if yearsExperience != nil {
		years := *yearsExperience
		return s.normalizeSeniorityFromYears(years)
	}

	// Default fallback
	return "Mid"
}

// extractCurrentPosition extracts the job title/role from text (e.g., "Senior Software Engineer")
// This is different from current_project which would be the specific project name
func (s *CandidateExtractService) extractCurrentPosition(text string) string {
	// Look for current position indicators
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "currently") || strings.Contains(lower, "working as") ||
			strings.Contains(lower, "position") || strings.Contains(lower, "role") {
			// Extract the position from the line
			words := strings.Fields(line)
			for i, word := range words {
				if strings.Contains(strings.ToLower(word), "developer") ||
					strings.Contains(strings.ToLower(word), "engineer") ||
					strings.Contains(strings.ToLower(word), "manager") ||
					strings.Contains(strings.ToLower(word), "analyst") {
					// Try to get the full position title
					if i > 0 && i < len(words)-1 {
						return strings.Join(words[i-1:i+2], " ")
					}
					return word
				}
			}
		}
	}

	// Look for common position patterns
	words := strings.Fields(text)
	for i, word := range words {
		lower := strings.ToLower(word)
		if strings.Contains(lower, "developer") || strings.Contains(lower, "engineer") {
			// Try to get the full position title
			if i > 0 {
				return strings.Join(words[i-1:i+1], " ")
			}
			return word
		}
	}

	// Default fallback
	return "Software Developer"
}

func (s *CandidateExtractService) calculateMatchScore(nerResult *SkillExtractionResult) int {
	// Base score
	score := 50

	// Add points for skills based on dynamic categories
	totalSkills := 0
	for _, skillList := range nerResult.Skills.Categories {
		totalSkills += len(skillList)
	}

	// Add points based on total skills found (simplified scoring)
	score += totalSkills * 2

	// Add points for experience
	if len(nerResult.Skills.YearsOfExperience) > 0 {
		maxYears := s.getMaxYearsFromStrings(nerResult.Skills.YearsOfExperience)
		if maxYears >= 5 {
			score += 15
		} else if maxYears >= 3 {
			score += 10
		} else if maxYears > 0 {
			score += 5
		}
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

func (s *CandidateExtractService) generateRecommendation(matchScore int) string {
	if matchScore >= 80 {
		return "Highly Recommended"
	} else if matchScore >= 60 {
		return "Recommended"
	} else if matchScore >= 40 {
		return "Consider with reservations"
	}
	return "Not recommended"
}

// GetHealthStatus returns the health status of the extraction service
func (s *CandidateExtractService) GetHealthStatus() (string, error) {
	// Check if NER service is initialized
	if s.nerService == nil {
		return "unhealthy", fmt.Errorf("NER service not initialized")
	}

	// Service is healthy and using pure NER
	return "healthy (Pure NER)", nil
}

// estimateTokens provides a rough estimate of token count
func (s *CandidateExtractService) estimateTokens(text string) int {
	// Rough estimation: 1 token ≈ 3-4 characters
	return len(text) / 3
}
