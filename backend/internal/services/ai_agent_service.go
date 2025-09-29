package services

import (
	"fmt"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/matching"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"strings"
	"time"
)

type aiAgentService struct {
	aiAgentRepo         repositories.AIAgentRepository
	employeeRepo        repositories.EmployeeRepository
	skillRepo           repositories.SkillRepository
	matchRepo           repositories.MatchRepository
	matchEngine         *matching.MatchEngine
	notificationService NotificationService
	nerService          *NERService
}

// NewAIAgentService creates a new AI agent service
func NewAIAgentService(
	aiAgentRepo repositories.AIAgentRepository,
	employeeRepo repositories.EmployeeRepository,
	skillRepo repositories.SkillRepository,
	matchRepo repositories.MatchRepository,
	notificationService NotificationService,
) AIAgentService {
	return &aiAgentService{
		aiAgentRepo:         aiAgentRepo,
		employeeRepo:        employeeRepo,
		skillRepo:           skillRepo,
		matchRepo:           matchRepo,
		matchEngine:         matching.NewMatchEngine(),
		notificationService: notificationService,
		nerService:          NewNERService(),
	}
}

func (s *aiAgentService) CreateAIAgentRequest(req *models.CreateAIAgentRequest) (*models.AIAgentRequest, error) {
	aiRequest := &models.AIAgentRequest{
		TeamsMessageID: req.TeamsMessageID,
		ChannelID:      req.ChannelID,
		UserID:         req.UserID,
		UserName:       req.UserName,
		MessageText:    req.MessageText,
		AttachmentURL:  req.AttachmentURL,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	result, err := s.aiAgentRepo.Create(aiRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *aiAgentService) GetAIAgentRequest(id int) (*models.AIAgentRequest, error) {
	return s.aiAgentRepo.GetByID(id)
}

func (s *aiAgentService) UpdateAIAgentRequest(id int, req *models.AIAgentRequest) error {
	return s.aiAgentRepo.Update(id, req)
}

func (s *aiAgentService) ProcessAIAgentRequest(id int) (*models.AIAgentResponse, error) {
	startTime := time.Now()

	// Get the request
	request, err := s.aiAgentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update status to processing
	request.Status = "processing"
	if err := s.aiAgentRepo.Update(id, request); err != nil {
		return nil, err
	}

	// Extract text from message or attachment
	extractedText, err := s.extractText(request)
	if err != nil {
		request.Status = "failed"
		errorMsg := fmt.Sprintf("Text extraction failed: %v", err)
		request.Error = &errorMsg
		if updateErr := s.aiAgentRepo.Update(id, request); updateErr != nil {
			s.notificationService.LogError(id, fmt.Sprintf("Failed to update request status after text extraction error: %v", updateErr))
		}
		s.notificationService.LogError(id, *request.Error)
		return nil, err
	}

	// Extract skills from text using NER service
	nerResult, err := s.nerService.ExtractSkillsFromText(extractedText)
	if err != nil {
		request.Status = "failed"
		errorMsg := fmt.Sprintf("NER skill extraction failed: %v", err)
		request.Error = &errorMsg
		if updateErr := s.aiAgentRepo.Update(id, request); updateErr != nil {
			s.notificationService.LogError(id, fmt.Sprintf("Failed to update request status after NER error: %v", updateErr))
		}
		s.notificationService.LogError(id, *request.Error)
		return nil, err
	}

	// Extract skills from NER result
	var allSkills []string
	if nerResult != nil {
		allSkills = append(allSkills, nerResult.Skills.ProgrammingLanguages...)
		allSkills = append(allSkills, nerResult.Skills.WebTechnologies...)
		allSkills = append(allSkills, nerResult.Skills.Databases...)
		allSkills = append(allSkills, nerResult.Skills.CloudDevOps...)
		allSkills = append(allSkills, nerResult.Skills.SoftSkills...)
		allSkills = append(allSkills, nerResult.Skills.ToolsFrameworks...)
	}

	// Deduplicate skills
	skills := s.deduplicateSkills(allSkills)

	// Update request with extracted data
	request.ExtractedText = &extractedText
	request.ExtractedSkills = skills
	request.Status = "completed"
	now := time.Now()
	request.ProcessedAt = &now
	if err := s.aiAgentRepo.Update(id, request); err != nil {
		// If update fails, try to set status to failed
		request.Status = "failed"
		errorMsg := fmt.Sprintf("Failed to update request status: %v", err)
		request.Error = &errorMsg
		s.aiAgentRepo.Update(id, request) // Try to update with error status
		s.notificationService.LogError(id, errorMsg)
		return nil, err
	}

	// Find matching employees
	matches, err := s.findMatchingEmployees(skills)
	if err != nil {
		request.Status = "failed"
		errorMsg := fmt.Sprintf("Employee matching failed: %v", err)
		request.Error = &errorMsg
		if updateErr := s.aiAgentRepo.Update(id, request); updateErr != nil {
			s.notificationService.LogError(id, fmt.Sprintf("Failed to update request status after matching error: %v", updateErr))
		}
		s.notificationService.LogError(id, *request.Error)
		return nil, err
	}

	// Save matches to database
	for _, match := range matches {
		_, err := s.matchRepo.Create(&match)
		if err != nil {
			// Log error but don't fail the request
			s.notificationService.LogError(id, fmt.Sprintf("Failed to save match for employee %d: %v", match.EmployeeID, err))
		}
	}

	// Generate explanations for matches
	aiMatches := s.generateMatchExplanations(matches, skills)

	// Create response
	response := &models.AIAgentResponse{
		RequestID:      id,
		Matches:        aiMatches,
		Summary:        s.generateSummary(aiMatches, skills),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Status:         "completed",
	}

	// Save response to database
	if err := s.aiAgentRepo.SaveResponse(response); err != nil {
		// Log error but don't fail the request
		s.notificationService.LogError(id, fmt.Sprintf("Failed to save response: %v", err))
	}

	return response, nil
}

func (s *aiAgentService) ExtractSkillsFromText(text string) (*models.SkillExtractionResponse, error) {
	// Use NER service for skill extraction
	nerResult, err := s.nerService.ExtractSkillsFromText(text)
	if err != nil {
		return nil, fmt.Errorf("NER skill extraction failed: %v", err)
	}

	// Extract skills from NER result
	var allSkills []string
	if nerResult != nil {
		allSkills = append(allSkills, nerResult.Skills.ProgrammingLanguages...)
		allSkills = append(allSkills, nerResult.Skills.WebTechnologies...)
		allSkills = append(allSkills, nerResult.Skills.Databases...)
		allSkills = append(allSkills, nerResult.Skills.CloudDevOps...)
		allSkills = append(allSkills, nerResult.Skills.SoftSkills...)
		allSkills = append(allSkills, nerResult.Skills.ToolsFrameworks...)
	}

	// Deduplicate and return
	skills := s.deduplicateSkills(allSkills)

	return &models.SkillExtractionResponse{
		Skills: skills,
		Text:   text,
	}, nil
}

func (s *aiAgentService) GetAIAgentRequests(limit int, offset int) ([]models.AIAgentRequest, error) {
	return s.aiAgentRepo.GetAll(limit, offset)
}

func (s *aiAgentService) GetAIAgentResponse(requestID int) (*models.AIAgentResponse, error) {
	return s.aiAgentRepo.GetResponseByRequestID(requestID)
}

// extractText extracts text from message or attachment
func (s *aiAgentService) extractText(request *models.AIAgentRequest) (string, error) {
	// If there's an attachment, extract text from it
	if request.AttachmentURL != nil && *request.AttachmentURL != "" {
		return s.extractTextFromAttachment(*request.AttachmentURL)
	}

	// Otherwise use the message text
	return request.MessageText, nil
}

// extractTextFromAttachment extracts text from an attachment URL
func (s *aiAgentService) extractTextFromAttachment(url string) (string, error) {
	// For now, we'll implement a simple text extraction
	// In a real implementation, you would use libraries like:
	// - For PDFs: github.com/ledongthuc/pdf
	// - For Word docs: github.com/unidoc/unioffice
	// - For other formats: appropriate libraries

	// This is a placeholder implementation
	return "Extracted text from attachment", nil
}

// findMatchingEmployees finds employees matching the extracted skills
func (s *aiAgentService) findMatchingEmployees(skills []string) ([]models.Match, error) {
	// Get all employees
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Create a search request
	searchReq := &models.SearchRequest{
		RequiredSkills: skills,
		MinMatchScore:  0.1, // Low threshold to get more results
	}

	// Use the existing matching engine
	matches := s.matchEngine.SearchEmployees(searchReq, employees)

	// Return top 5 matches
	if len(matches) > 5 {
		matches = matches[:5]
	}

	return matches, nil
}

// generateMatchExplanations generates explanations for each match
func (s *aiAgentService) generateMatchExplanations(matches []models.Match, extractedSkills []string) []models.AIAgentMatch {
	var aiMatches []models.AIAgentMatch

	for _, match := range matches {
		// Get resume link for the employee
		resumeLink := s.getResumeLink(match.EmployeeID)

		aiMatch := models.AIAgentMatch{
			EmployeeID:     match.EmployeeID,
			EmployeeName:   match.Employee.Name,
			EmployeeEmail:  match.Employee.Email,
			Position:       s.getPositionFromDepartment(match.Employee.Department),
			Seniority:      match.Employee.Level,
			Location:       match.Employee.Location,
			CurrentProject: s.getCurrentProject(match.Employee.CurrentProject),
			ResumeLink:     resumeLink,
			MatchScore:     match.MatchScore,
			MatchingSkills: match.MatchingSkills,
			AISummary:      s.generateAISummary(match, extractedSkills),
			Bio:            match.Employee.Bio,
		}
		aiMatches = append(aiMatches, aiMatch)
	}

	return aiMatches
}

// generateExplanation generates a human-readable explanation for a match
func (s *aiAgentService) generateExplanation(match models.Match, extractedSkills []string) string {
	employee := match.Employee

	// Count matching skills
	matchingCount := len(match.MatchingSkills)
	totalCount := len(extractedSkills)

	explanation := fmt.Sprintf("%s is a %s %s in %s",
		employee.Name,
		strings.ToLower(employee.Level),
		employee.Department,
		employee.Location)

	if matchingCount > 0 {
		explanation += fmt.Sprintf(" with %d matching skills: %s",
			matchingCount,
			strings.Join(match.MatchingSkills, ", "))
	}

	if matchingCount < totalCount {
		missingSkills := s.findMissingSkills(extractedSkills, match.MatchingSkills)
		explanation += fmt.Sprintf(". Missing skills: %s",
			strings.Join(missingSkills, ", "))
	}

	if employee.Bio != "" {
		explanation += fmt.Sprintf(". Bio: %s", employee.Bio)
	}

	return explanation
}

// findMissingSkills finds skills that were requested but not matched
func (s *aiAgentService) findMissingSkills(requested, matched []string) []string {
	matchedMap := make(map[string]bool)
	for _, skill := range matched {
		matchedMap[strings.ToLower(skill)] = true
	}

	var missing []string
	for _, skill := range requested {
		if !matchedMap[strings.ToLower(skill)] {
			missing = append(missing, skill)
		}
	}

	return missing
}

// generateSummary generates a summary of all matches
func (s *aiAgentService) generateSummary(matches []models.AIAgentMatch, extractedSkills []string) string {
	if len(matches) == 0 {
		return fmt.Sprintf("No employees found matching the skills: %s", strings.Join(extractedSkills, ", "))
	}

	summary := fmt.Sprintf("Found %d employees matching skills: %s\n\n",
		len(matches),
		strings.Join(extractedSkills, ", "))

	for i, match := range matches {
		summary += fmt.Sprintf("%d. %s (%s) - Score: %.1f - %s\n",
			i+1,
			match.EmployeeName,
			match.Position,
			match.MatchScore,
			strings.Join(match.MatchingSkills, ", "))
	}

	return summary
}

// deduplicateSkills removes duplicate skills from the list
func (s *aiAgentService) deduplicateSkills(skills []string) []string {
	skillMap := make(map[string]bool)
	var uniqueSkills []string

	for _, skill := range skills {
		skill = strings.TrimSpace(skill)
		if skill != "" && !skillMap[strings.ToLower(skill)] {
			skillMap[strings.ToLower(skill)] = true
			uniqueSkills = append(uniqueSkills, skill)
		}
	}

	return uniqueSkills
}

// getResumeLink retrieves the resume link for an employee
func (s *aiAgentService) getResumeLink(employeeID int) string {
	// For now, return a placeholder URL
	// In a real implementation, you would query the uploaded_files table
	// to find the most recent resume for this employee
	return fmt.Sprintf("/api/v1/files/resume/%d", employeeID)
}

// getPositionFromDepartment maps department to position title
func (s *aiAgentService) getPositionFromDepartment(department string) string {
	positionMap := constants.SkillCategoryMappings

	if position, exists := positionMap[department]; exists {
		return position
	}
	return department + " Specialist"
}

// generateAISummary generates an AI-powered summary for the match
func (s *aiAgentService) generateAISummary(match models.Match, extractedSkills []string) string {
	employee := match.Employee
	matchingCount := len(match.MatchingSkills)
	totalCount := len(extractedSkills)

	// Calculate match percentage
	matchPercentage := float64(matchingCount) / float64(totalCount) * 100

	// Generate summary based on match quality
	var summary string

	if matchPercentage >= 80 {
		summary = fmt.Sprintf("Excellent match! %s is a %s %s with %d out of %d required skills (%.0f%% match). ",
			employee.Name, strings.ToLower(employee.Level), s.getPositionFromDepartment(employee.Department),
			matchingCount, totalCount, matchPercentage)
	} else if matchPercentage >= 60 {
		summary = fmt.Sprintf("Strong match! %s is a %s %s with %d out of %d required skills (%.0f%% match). ",
			employee.Name, strings.ToLower(employee.Level), s.getPositionFromDepartment(employee.Department),
			matchingCount, totalCount, matchPercentage)
	} else if matchPercentage >= 40 {
		summary = fmt.Sprintf("Good match! %s is a %s %s with %d out of %d required skills (%.0f%% match). ",
			employee.Name, strings.ToLower(employee.Level), s.getPositionFromDepartment(employee.Department),
			matchingCount, totalCount, matchPercentage)
	} else {
		summary = fmt.Sprintf("Partial match! %s is a %s %s with %d out of %d required skills (%.0f%% match). ",
			employee.Name, strings.ToLower(employee.Level), s.getPositionFromDepartment(employee.Department),
			matchingCount, totalCount, matchPercentage)
	}

	// Add location and project information
	if employee.Location != "" {
		summary += fmt.Sprintf("Located in %s. ", employee.Location)
	}

	if employee.CurrentProject != nil && *employee.CurrentProject != "" {
		summary += fmt.Sprintf("Currently working on: %s. ", *employee.CurrentProject)
	}

	// Add matching skills
	if matchingCount > 0 {
		summary += fmt.Sprintf("Key matching skills: %s.", strings.Join(match.MatchingSkills, ", "))
	}

	// Add missing skills if any
	if matchingCount < totalCount {
		missingSkills := s.findMissingSkills(extractedSkills, match.MatchingSkills)
		if len(missingSkills) > 0 {
			summary += fmt.Sprintf(" Missing skills: %s.", strings.Join(missingSkills, ", "))
		}
	}

	return summary
}

// getCurrentProject safely extracts the current project string from a pointer
func (s *aiAgentService) getCurrentProject(project *string) string {
	if project == nil {
		return ""
	}
	return *project
}
