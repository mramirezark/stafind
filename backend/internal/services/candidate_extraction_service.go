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
			return s.createNewEmployee(originalText, extractedData, extractionSource, startTime)
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
	return s.updateExistingEmployee(existingEmployee, originalText, extractedData, extractionSource, startTime)
}

// createNewEmployee creates a new employee from extracted data
func (s *CandidateStorageService) createNewEmployee(
	originalText string,
	extractedData map[string]interface{},
	extractionSource string,
	startTime time.Time,
) (*models.CandidateExtractionResult, error) {
	// Extract basic information
	candidateName, _ := extractedData["candidate_name"].(string)
	contactInfo, _ := extractedData["contact_info"].(map[string]interface{})
	email, _ := contactInfo["email"].(string)
	location, _ := contactInfo["location"].(string)

	seniorityLevel, _ := extractedData["seniority_level"].(string)
	currentPosition, _ := extractedData["current_position"].(string)
	summary, _ := extractedData["summary"].(string)

	// Extract last project from work experience
	lastProject := s.extractLastProject(originalText)

	// Create employee request
	createReq := &models.CreateEmployeeRequest{
		Name:           candidateName,
		Email:          email,
		Department:     s.mapSeniorityToDepartment(seniorityLevel),
		Level:          currentPosition, // Job title/role (e.g., "Senior Software Engineer")
		Location:       location,
		Bio:            summary,
		CurrentProject: lastProject, // Last project they worked on
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
	seniorityLevel, _ := extractedData["seniority_level"].(string)
	summary, _ := extractedData["summary"].(string)
	currentPosition, _ := extractedData["current_position"].(string)

	// Extract last project from work experience
	lastProject := s.extractLastProject(originalText)

	// Update basic information with extraction data
	updateReq := &models.CreateEmployeeRequest{
		Name:       existingEmployee.Name,
		Email:      existingEmployee.Email,
		Department: s.mapSeniorityToDepartment(seniorityLevel),
		Level:      currentPosition, // Job title/role (e.g., "Senior Software Engineer")
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
		Skills: s.extractSkillsFromData(extractedData),
	}

	updatedEmployee, err := s.employeeRepo.UpdateWithExtraction(
		existingEmployee.ID,
		updateReq,
		originalText,
		extractedData,
		extractionSource,
		"completed",
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

// extractSkillsFromData extracts skills from extracted data
func (s *CandidateStorageService) extractSkillsFromData(extractedData map[string]interface{}) []models.EmployeeSkillReq {
	var skills []models.EmployeeSkillReq

	// Extract technical skills
	if skillsData, ok := extractedData["skills"].(map[string]interface{}); ok {
		// Technical skills
		if technical, ok := skillsData["technical"].([]interface{}); ok {
			for _, skill := range technical {
				if skillName, ok := skill.(string); ok {
					skills = append(skills, models.EmployeeSkillReq{
						SkillName:        skillName,
						ProficiencyLevel: 3, // Default proficiency level
						YearsExperience:  1.0,
					})
				}
			}
		}

		// Frameworks
		if frameworks, ok := skillsData["frameworks"].([]interface{}); ok {
			for _, skill := range frameworks {
				if skillName, ok := skill.(string); ok {
					skills = append(skills, models.EmployeeSkillReq{
						SkillName:        skillName,
						ProficiencyLevel: 3,
						YearsExperience:  1.0,
					})
				}
			}
		}

		// Tools
		if tools, ok := skillsData["tools"].([]interface{}); ok {
			for _, skill := range tools {
				if skillName, ok := skill.(string); ok {
					skills = append(skills, models.EmployeeSkillReq{
						SkillName:        skillName,
						ProficiencyLevel: 3,
						YearsExperience:  1.0,
					})
				}
			}
		}

		// Languages
		if languages, ok := skillsData["languages"].([]interface{}); ok {
			for _, skill := range languages {
				if skillName, ok := skill.(string); ok {
					skills = append(skills, models.EmployeeSkillReq{
						SkillName:        skillName,
						ProficiencyLevel: 3,
						YearsExperience:  1.0,
					})
				}
			}
		}
	}

	return skills
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
