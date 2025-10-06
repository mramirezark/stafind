package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"strings"
	"time"
)

// CandidateStorageService handles candidate extraction and storage
type CandidateStorageService struct {
	employeeRepo repositories.EmployeeRepository
	skillRepo    repositories.SkillRepository
}

// NewCandidateStorageService creates a new candidate storage service
func NewCandidateStorageService(employeeRepo repositories.EmployeeRepository, skillRepo repositories.SkillRepository) *CandidateStorageService {
	return &CandidateStorageService{
		employeeRepo: employeeRepo,
		skillRepo:    skillRepo,
	}
}

// ProcessCandidateExtraction processes candidate extraction and stores/updates employee data
func (s *CandidateStorageService) ProcessCandidateExtraction(
	originalText string,
	extractedData map[string]interface{},
	extractionSource string,
	resumeURL string,
) (*models.CandidateExtractionResult, error) {
	startTime := time.Now()

	// Extract candidate name and email from extracted data
	candidateName, _ := extractedData["candidate_name"].(string)
	candidateEmail, _ := extractedData["contact_info"].(map[string]interface{})["email"].(string)

	if candidateName == "" || candidateEmail == "" {
		return &models.CandidateExtractionResult{
			Status:         "failed",
			Message:        "Candidate name and email are required",
			ProcessingTime: time.Since(startTime),
		}, fmt.Errorf("candidate name and email are required")
	}

	// Check if employee already exists by email
	fmt.Printf("DEBUG: Checking for existing employee with email: %s\n", candidateEmail)
	existingEmployee, err := s.employeeRepo.GetByEmail(candidateEmail)
	if err != nil {
		fmt.Printf("DEBUG: GetByEmail error: %v\n", err)
		// Check if it's a "no rows" error (employee doesn't exist)
		if err == sql.ErrNoRows {
			fmt.Printf("DEBUG: Employee not found, creating new one\n")
			// Employee doesn't exist, create new one
			return s.createNewEmployee(originalText, extractedData, extractionSource, startTime, resumeURL)
		}
		// Some other database error occurred
		fmt.Printf("DEBUG: Database error: %v\n", err)
		return &models.CandidateExtractionResult{
			Status:         "failed",
			Message:        fmt.Sprintf("Database error while checking for existing employee: %v", err),
			ProcessingTime: time.Since(startTime),
		}, err
	}
	fmt.Printf("DEBUG: Found existing employee: %s\n", existingEmployee.Name)

	// Employee exists, check for changes
	return s.updateExistingEmployee(existingEmployee, originalText, extractedData, extractionSource, startTime, resumeURL)
}

// createNewEmployee creates a new employee from extracted data
func (s *CandidateStorageService) createNewEmployee(
	originalText string,
	extractedData map[string]interface{},
	extractionSource string,
	startTime time.Time,
	resumeURL string,
) (*models.CandidateExtractionResult, error) {
	// Extract basic information
	candidateName, _ := extractedData["candidate_name"].(string)
	contactInfo, _ := extractedData["contact_info"].(map[string]interface{})
	email, _ := contactInfo["email"].(string)
	location, _ := contactInfo["location"].(string)

	summary, _ := extractedData["summary"].(string)

	// Extract and normalize seniority level using improved extraction
	seniorityLevel := s.extractAndNormalizeSeniorityLevel(extractedData, originalText)

	// Extract last project from work experience
	lastProject := s.extractLastProject(originalText)

	// Create employee request
	createReq := &models.CreateEmployeeRequest{
		Name:           candidateName,
		Email:          email,
		Department:     s.mapSeniorityToDepartment(seniorityLevel),
		Level:          seniorityLevel, // Seniority level (e.g., "Senior", "Mid", "Junior")
		Location:       location,
		Bio:            summary,
		CurrentProject: lastProject, // Last project they worked on
		ResumeUrl:      resumeURL,   // Resume URL from the extraction request
		Skills:         s.extractSkillsFromData(extractedData),
	}

	// Create employee with extraction data
	fmt.Printf("DEBUG: Creating new employee with extraction data, email: %s\n", email)
	employee, err := s.employeeRepo.CreateWithExtraction(
		createReq,
		originalText,
		extractedData,
		extractionSource,
		"completed",
		resumeURL,
	)
	if err != nil {
		fmt.Printf("DEBUG: Failed to create employee with extraction data: %v\n", err)
		return &models.CandidateExtractionResult{
			Status:         "failed",
			Message:        fmt.Sprintf("Failed to create employee with extraction data: %v", err),
			ProcessingTime: time.Since(startTime),
		}, err
	}
	fmt.Printf("DEBUG: Successfully created employee with extraction data, ID: %d\n", employee.ID)

	return &models.CandidateExtractionResult{
		EmployeeID:      employee.ID,
		Action:          "created",
		Employee:        employee,
		ExtractedData:   extractedData,
		ChangesDetected: true,
		ChangesSummary:  []string{"New employee created"},
		ProcessingTime:  time.Since(startTime),
		Status:          "completed",
		Message:         "Employee created successfully",
	}, nil
}

// updateExistingEmployee updates an existing employee and checks for changes
func (s *CandidateStorageService) updateExistingEmployee(
	existingEmployee *models.Employee,
	originalText string,
	extractedData map[string]interface{},
	extractionSource string,
	startTime time.Time,
	resumeURL string,
) (*models.CandidateExtractionResult, error) {
	changesDetected := false
	changesSummary := []string{}

	// Check if original text has changed
	if existingEmployee.OriginalText == nil || *existingEmployee.OriginalText != originalText {
		changesDetected = true
		changesSummary = append(changesSummary, "Original text updated")
	}

	// Check for changes in extracted data
	if existingEmployee.ExtractedData != nil {
		existingDataJSON, _ := json.Marshal(existingEmployee.ExtractedData)
		newDataJSON, _ := json.Marshal(extractedData)
		if string(existingDataJSON) != string(newDataJSON) {
			changesDetected = true
			changesSummary = append(changesSummary, "Extracted data updated")
		}
	} else {
		changesDetected = true
		changesSummary = append(changesSummary, "Extracted data added")
	}

	// If no changes detected, return without updating
	if !changesDetected {
		return &models.CandidateExtractionResult{
			EmployeeID:      existingEmployee.ID,
			Action:          "no_changes",
			Employee:        existingEmployee,
			ExtractedData:   extractedData,
			ChangesDetected: false,
			ChangesSummary:  []string{"No changes detected"},
			ProcessingTime:  time.Since(startTime),
			Status:          "completed",
			Message:         "No changes detected, employee data is up to date",
		}, nil
	}

	// Update employee with new data
	now := time.Now()
	existingEmployee.OriginalText = &originalText
	existingEmployee.ExtractedData = extractedData
	existingEmployee.ExtractionTimestamp = &now
	existingEmployee.ExtractionSource = &extractionSource
	status := "completed"
	existingEmployee.ExtractionStatus = &status

	// Extract updated information
	contactInfo, _ := extractedData["contact_info"].(map[string]interface{})
	location, _ := contactInfo["location"].(string)
	summary, _ := extractedData["summary"].(string)

	// Extract and normalize seniority level using improved extraction
	seniorityLevel := s.extractAndNormalizeSeniorityLevel(extractedData, originalText)

	// Extract last project from work experience
	lastProject := s.extractLastProject(originalText)

	// Update basic information with extraction data
	updateReq := &models.CreateEmployeeRequest{
		Name:       existingEmployee.Name,
		Email:      existingEmployee.Email,
		Department: s.mapSeniorityToDepartment(seniorityLevel),
		Level:      seniorityLevel, // Seniority level (e.g., "Senior", "Mid", "Junior")
		Location:   location,
		Bio:        summary,
		CurrentProject: func() string {
			// Use extracted project if available, otherwise keep existing
			if lastProject != "" {
				return lastProject
			}
			if existingEmployee.CurrentProject != nil {
				return *existingEmployee.CurrentProject
			}
			return ""
		}(), // Last project they worked on or existing project
		ResumeUrl: resumeURL, // Resume URL from the extraction request
		Skills:    s.extractSkillsFromData(extractedData),
	}

	updatedEmployee, err := s.employeeRepo.UpdateWithExtraction(
		existingEmployee.ID,
		updateReq,
		originalText,
		extractedData,
		extractionSource,
		"completed",
		resumeURL,
	)
	if err != nil {
		return &models.CandidateExtractionResult{
			EmployeeID:      existingEmployee.ID,
			Action:          "failed",
			Employee:        existingEmployee,
			ExtractedData:   extractedData,
			ChangesDetected: changesDetected,
			ChangesSummary:  changesSummary,
			ProcessingTime:  time.Since(startTime),
			Status:          "failed",
			Message:         fmt.Sprintf("Failed to update employee: %v", err),
		}, err
	}

	return &models.CandidateExtractionResult{
		EmployeeID:      updatedEmployee.ID,
		Action:          "updated",
		Employee:        updatedEmployee,
		ExtractedData:   extractedData,
		ChangesDetected: changesDetected,
		ChangesSummary:  changesSummary,
		ProcessingTime:  time.Since(startTime),
		Status:          "completed",
		Message:         "Employee updated successfully",
	}, nil
}

// extractAndNormalizeSeniorityLevel extracts and normalizes seniority level from various formats
func (s *CandidateStorageService) extractAndNormalizeSeniorityLevel(extractedData map[string]interface{}, originalText string) string {
	// First try to get from extracted data
	if seniorityLevel, ok := extractedData["seniority_level"].(string); ok && seniorityLevel != "" {
		return s.normalizeSeniorityLevel(seniorityLevel)
	}

	// If not found in extracted data, try to extract from original text
	return s.extractSeniorityFromText(originalText)
}

// normalizeSeniorityLevel normalizes various seniority level formats to standard values
func (s *CandidateStorageService) normalizeSeniorityLevel(level string) string {
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
	case strings.Contains(level, "years") && strings.Contains(level, "5"):
		return "Senior"
	case strings.Contains(level, "years") && (strings.Contains(level, "3") || strings.Contains(level, "4")):
		return "Mid"
	case strings.Contains(level, "years") && (strings.Contains(level, "1") || strings.Contains(level, "2")):
		return "Junior"
	default:
		// If we can't determine, try to infer from years of experience
		return s.inferSeniorityFromExperience(level)
	}
}

// extractSeniorityFromText extracts seniority level from original resume text
func (s *CandidateStorageService) extractSeniorityFromText(text string) string {
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
				return s.normalizeSeniorityLevel(line)
			}
		}
	}

	// Default fallback
	return "Mid"
}

// inferSeniorityFromExperience tries to infer seniority from experience-related text
func (s *CandidateStorageService) inferSeniorityFromExperience(text string) string {
	text = strings.ToLower(text)

	// Look for years of experience patterns
	if strings.Contains(text, "years") {
		// Extract number of years if possible
		// This is a simple implementation - could be enhanced with regex
		if strings.Contains(text, "5") || strings.Contains(text, "6") || strings.Contains(text, "7") ||
			strings.Contains(text, "8") || strings.Contains(text, "9") || strings.Contains(text, "10") {
			return "Senior"
		}
		if strings.Contains(text, "3") || strings.Contains(text, "4") {
			return "Mid"
		}
		if strings.Contains(text, "1") || strings.Contains(text, "2") {
			return "Junior"
		}
	}

	return "Mid" // Default fallback
}

// mapSeniorityToDepartment maps seniority level to department
func (s *CandidateStorageService) mapSeniorityToDepartment(seniorityLevel string) string {
	level := strings.ToLower(seniorityLevel)
	switch {
	case strings.Contains(level, "senior") || strings.Contains(level, "lead") || strings.Contains(level, "principal"):
		return "Engineering"
	case strings.Contains(level, "junior") || strings.Contains(level, "entry"):
		return "Engineering"
	case strings.Contains(level, "manager") || strings.Contains(level, "director"):
		return "Management"
	case strings.Contains(level, "architect"):
		return "Architecture"
	default:
		return "Engineering"
	}
}

// extractSkillsFromData extracts skills from extracted data (NER format)
func (s *CandidateStorageService) extractSkillsFromData(extractedData map[string]interface{}) []models.EmployeeSkillReq {
	var skills []models.EmployeeSkillReq

	// Extract skills from NER format - skills is directly the Categories map
	if categories, ok := extractedData["skills"].(map[string]interface{}); ok {
		// Handle the NER service format: map[string][]string
		for categoryName, skillList := range categories {
			if skillArray, ok := skillList.([]interface{}); ok {
				for _, skill := range skillArray {
					if skillName, ok := skill.(string); ok {
						skills = append(skills, models.EmployeeSkillReq{
							SkillName:        skillName,
							ProficiencyLevel: s.getProficiencyLevelForCategory(categoryName),
							YearsExperience:  s.getYearsExperienceForCategory(categoryName),
						})
					}
				}
			}
		}
	}

	return skills
}

// getProficiencyLevelForCategory returns appropriate proficiency level based on skill category
func (s *CandidateStorageService) getProficiencyLevelForCategory(categoryName string) int {
	category := strings.ToLower(categoryName)

	// Higher proficiency for core technical skills
	switch {
	case strings.Contains(category, "programming") || strings.Contains(category, "language"):
		return 4
	case strings.Contains(category, "framework") || strings.Contains(category, "library"):
		return 3
	case strings.Contains(category, "tool") || strings.Contains(category, "software"):
		return 3
	case strings.Contains(category, "database") || strings.Contains(category, "cloud"):
		return 3
	case strings.Contains(category, "methodology") || strings.Contains(category, "process"):
		return 2
	default:
		return 3 // Default proficiency level
	}
}

// getYearsExperienceForCategory returns appropriate years of experience based on skill category
func (s *CandidateStorageService) getYearsExperienceForCategory(categoryName string) float64 {
	category := strings.ToLower(categoryName)

	// More experience typically needed for advanced categories
	switch {
	case strings.Contains(category, "programming") || strings.Contains(category, "language"):
		return 2.0
	case strings.Contains(category, "framework") || strings.Contains(category, "library"):
		return 1.5
	case strings.Contains(category, "tool") || strings.Contains(category, "software"):
		return 1.0
	case strings.Contains(category, "database") || strings.Contains(category, "cloud"):
		return 1.5
	case strings.Contains(category, "methodology") || strings.Contains(category, "process"):
		return 1.0
	default:
		return 1.0 // Default years of experience
	}
}

// extractLastProject extracts the last project from work experience text
func (s *CandidateStorageService) extractLastProject(text string) string {
	// Look for work experience patterns
	lines := strings.Split(text, "\n")

	// Common patterns for project names in resumes
	projectPatterns := []string{
		"project",
		"developed",
		"built",
		"created",
		"implemented",
		"designed",
		"architected",
		"led",
		"managed",
		"delivered",
	}

	// Look for the most recent work experience section
	var lastProject string
	inWorkExperience := false

	fmt.Printf("DEBUG: Extracting project from text with %d lines\n", len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if we're entering work experience section
		if strings.Contains(strings.ToLower(line), "experience") ||
			strings.Contains(strings.ToLower(line), "employment") ||
			strings.Contains(strings.ToLower(line), "work history") {
			inWorkExperience = true
			fmt.Printf("DEBUG: Entering work experience section: %s\n", line)
			continue
		}

		// Check if we're leaving work experience section
		if inWorkExperience && (strings.Contains(strings.ToLower(line), "education") ||
			strings.Contains(strings.ToLower(line), "skills") ||
			strings.Contains(strings.ToLower(line), "certifications")) {
			fmt.Printf("DEBUG: Leaving work experience section: %s\n", line)
			break
		}

		if inWorkExperience {
			// Look for project-related keywords
			lowerLine := strings.ToLower(line)
			for _, pattern := range projectPatterns {
				if strings.Contains(lowerLine, pattern) {
					fmt.Printf("DEBUG: Found project pattern '%s' in line: %s\n", pattern, line)
					// Extract potential project name
					// Look for quoted text, capitalized words, or specific patterns
					project := s.extractProjectFromLine(line)
					fmt.Printf("DEBUG: Extracted project: '%s'\n", project)
					if project != "" {
						lastProject = project
					}
					break
				}
			}
		}
	}

	// If no project found in work experience, try to extract from summary or bio
	if lastProject == "" {
		fmt.Printf("DEBUG: No project found in work experience, trying summary\n")
		lastProject = s.extractProjectFromSummary(text)
	}

	fmt.Printf("DEBUG: Final extracted project: '%s'\n", lastProject)
	return lastProject
}

// extractProjectFromLine extracts project name from a specific line
func (s *CandidateStorageService) extractProjectFromLine(line string) string {
	// Look for quoted project names first (highest priority)
	if strings.Contains(line, "\"") {
		start := strings.Index(line, "\"")
		end := strings.LastIndex(line, "\"")
		if start != end && start != -1 {
			project := strings.TrimSpace(line[start+1 : end])
			if len(project) > 3 && len(project) < 100 { // Reasonable project name length
				return project
			}
		}
	}

	// Look for project patterns with common keywords
	words := strings.Fields(line)
	for i, word := range words {
		if len(word) > 2 && word[0] >= 'A' && word[0] <= 'Z' {
			// Check if this looks like a project name
			lowerWord := strings.ToLower(word)
			if strings.Contains(lowerWord, "platform") ||
				strings.Contains(lowerWord, "application") ||
				strings.Contains(lowerWord, "portal") ||
				strings.Contains(lowerWord, "banking") ||
				strings.Contains(lowerWord, "e-commerce") ||
				strings.Contains(lowerWord, "payment") {
				// Take this word and potentially the next few words
				project := word
				if i+1 < len(words) && len(words[i+1]) > 2 && words[i+1][0] >= 'A' && words[i+1][0] <= 'Z' {
					project += " " + words[i+1]
				}
				return project
			}
		}
	}

	// Look for system/service patterns as fallback
	for i, word := range words {
		if len(word) > 2 && word[0] >= 'A' && word[0] <= 'Z' {
			lowerWord := strings.ToLower(word)
			if strings.Contains(lowerWord, "system") ||
				strings.Contains(lowerWord, "service") ||
				strings.Contains(lowerWord, "project") {
				// Take this word and potentially the next few words
				project := word
				if i+1 < len(words) && len(words[i+1]) > 2 && words[i+1][0] >= 'A' && words[i+1][0] <= 'Z' {
					project += " " + words[i+1]
				}
				return project
			}
		}
	}

	return ""
}

// extractProjectFromSummary extracts project name from summary or bio section
func (s *CandidateStorageService) extractProjectFromSummary(text string) string {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "currently") ||
			strings.Contains(strings.ToLower(line), "recently") ||
			strings.Contains(strings.ToLower(line), "latest") {
			project := s.extractProjectFromLine(line)
			if project != "" {
				return project
			}
		}
	}

	return ""
}
