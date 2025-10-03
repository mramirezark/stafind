package services

import (
	"fmt"
	"log"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"strings"
	"time"
)

// MatchingService handles candidate matching based on skills and requirements
type MatchingService struct {
	employeeRepo repositories.EmployeeRepository
	skillRepo    repositories.SkillRepository
}

// NewMatchingService creates a new matching service
func NewMatchingService(employeeRepo repositories.EmployeeRepository, skillRepo repositories.SkillRepository) *MatchingService {
	return &MatchingService{
		employeeRepo: employeeRepo,
		skillRepo:    skillRepo,
	}
}

// MatchCandidatesToRequirements matches candidates based on job requirements
func (s *MatchingService) MatchCandidatesToRequirements(requirements string, extractedSkills map[string]interface{}) (*models.MatchingResult, error) {
	startTime := time.Now()
	log.Printf("Starting candidate matching process for requirements: %s", requirements)

	// Extract required skills from the requirements text
	requiredSkills := s.extractRequiredSkills(requirements)
	log.Printf("Extracted %d required skills from job requirements", len(requiredSkills))

	// Get all employees with their skills
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	log.Printf("Found %d employees to match against", len(employees))

	// Match each employee against requirements
	var matches []models.MatchingCandidate
	for _, employee := range employees {
		match := s.calculateMatch(employee, requiredSkills, extractedSkills)
		if match.MatchScore > 0 {
			matches = append(matches, match)
		}
	}

	// Sort matches by score (highest first)
	s.sortMatchesByScore(matches)

	// Limit to top 10 matches
	if len(matches) > 10 {
		matches = matches[:10]
	}

	processingTime := time.Since(startTime)
	log.Printf("Matching completed in %v, found %d matches", processingTime, len(matches))

	return &models.MatchingResult{
		Requirements:    requirements,
		RequiredSkills:  requiredSkills,
		Matches:         matches,
		TotalCandidates: len(employees),
		ProcessingTime:  processingTime,
		Timestamp:       time.Now(),
	}, nil
}

// extractRequiredSkills extracts skills from job requirements text
func (s *MatchingService) extractRequiredSkills(requirements string) []string {
	var skills []string
	words := strings.Fields(strings.ToLower(requirements))

	// Common technical skills to look for
	skillKeywords := map[string][]string{
		"programming": {"python", "java", "javascript", "typescript", "go", "rust", "c++", "c#", "php", "ruby", "swift", "kotlin"},
		"web":         {"react", "angular", "vue", "html", "css", "nodejs", "express", "django", "flask", "spring", "laravel"},
		"database":    {"postgresql", "mysql", "mongodb", "redis", "elasticsearch", "sqlite", "oracle", "sql server"},
		"cloud":       {"aws", "azure", "gcp", "docker", "kubernetes", "terraform", "jenkins", "ci/cd"},
		"tools":       {"git", "github", "gitlab", "jira", "confluence", "slack", "figma", "postman"},
		"frameworks":  {"django", "flask", "spring", "express", "rails", "laravel", "symfony", "asp.net"},
		"mobile":      {"react native", "flutter", "ios", "android", "xamarin", "ionic"},
		"data":        {"pandas", "numpy", "tensorflow", "pytorch", "scikit-learn", "spark", "hadoop"},
	}

	for _, word := range words {
		// Clean the word
		cleanWord := strings.Trim(word, ".,!?;:")

		// Check against skill keywords
		for _, skillList := range skillKeywords {
			for _, skill := range skillList {
				if strings.Contains(cleanWord, skill) {
					// Add the skill if not already present
					if !s.containsSkill(skills, skill) {
						skills = append(skills, skill)
					}
				}
			}
		}
	}

	return skills
}

// calculateMatch calculates how well an employee matches the requirements
func (s *MatchingService) calculateMatch(employee models.Employee, requiredSkills []string, extractedSkills map[string]interface{}) models.MatchingCandidate {
	// Get employee skills
	employeeSkills := s.getEmployeeSkills(employee)

	// Calculate skill matches
	matchedSkills := s.findMatchingSkills(employeeSkills, requiredSkills)
	matchScore := s.calculateMatchScore(matchedSkills, len(requiredSkills))

	// Calculate experience bonus
	experienceBonus := s.calculateExperienceBonus(employee, extractedSkills)

	// Calculate final score
	finalScore := matchScore + experienceBonus
	if finalScore > 100 {
		finalScore = 100
	}

	// Determine match level
	matchLevel := s.determineMatchLevel(finalScore)

	return models.MatchingCandidate{
		EmployeeID:      employee.ID,
		EmployeeName:    employee.Name,
		EmployeeEmail:   employee.Email,
		EmployeeLevel:   employee.Level,
		EmployeeSkills:  employeeSkills,
		MatchedSkills:   matchedSkills,
		MatchScore:      finalScore,
		MatchLevel:      matchLevel,
		ExperienceMatch: experienceBonus > 0,
		SkillsMatch:     len(matchedSkills),
		TotalRequired:   len(requiredSkills),
	}
}

// getEmployeeSkills extracts skills from employee data
func (s *MatchingService) getEmployeeSkills(employee models.Employee) []string {
	var skills []string

	// Add skills from the employee's skills relationship
	for _, skill := range employee.Skills {
		skills = append(skills, strings.ToLower(skill.Name))
	}

	// Add skills from extracted data if available
	if employee.ExtractedData != nil {
		if skillsData, ok := employee.ExtractedData["skills"].(map[string]interface{}); ok {
			// Programming languages
			if pl, ok := skillsData["programming_languages"].([]interface{}); ok {
				for _, skill := range pl {
					if skillStr, ok := skill.(string); ok {
						skills = append(skills, strings.ToLower(skillStr))
					}
				}
			}
			// Web technologies
			if wt, ok := skillsData["web_technologies"].([]interface{}); ok {
				for _, skill := range wt {
					if skillStr, ok := skill.(string); ok {
						skills = append(skills, strings.ToLower(skillStr))
					}
				}
			}
			// Databases
			if db, ok := skillsData["databases"].([]interface{}); ok {
				for _, skill := range db {
					if skillStr, ok := skill.(string); ok {
						skills = append(skills, strings.ToLower(skillStr))
					}
				}
			}
			// Cloud/DevOps
			if cd, ok := skillsData["cloud_devops"].([]interface{}); ok {
				for _, skill := range cd {
					if skillStr, ok := skill.(string); ok {
						skills = append(skills, strings.ToLower(skillStr))
					}
				}
			}
		}
	}

	return s.removeDuplicates(skills)
}

// findMatchingSkills finds which employee skills match the required skills
func (s *MatchingService) findMatchingSkills(employeeSkills []string, requiredSkills []string) []string {
	var matches []string

	for _, required := range requiredSkills {
		for _, employee := range employeeSkills {
			if strings.Contains(employee, required) || strings.Contains(required, employee) {
				if !s.containsSkill(matches, required) {
					matches = append(matches, required)
				}
				break
			}
		}
	}

	return matches
}

// calculateMatchScore calculates the match score based on matched skills
func (s *MatchingService) calculateMatchScore(matchedSkills []string, totalRequired int) int {
	if totalRequired == 0 {
		return 0
	}

	// Base score is percentage of skills matched
	baseScore := (len(matchedSkills) * 100) / totalRequired

	// Bonus for high match percentage
	if len(matchedSkills) == totalRequired {
		baseScore += 20 // Perfect match bonus
	} else if len(matchedSkills) >= totalRequired/2 {
		baseScore += 10 // Good match bonus
	}

	return baseScore
}

// calculateExperienceBonus calculates bonus points for experience match
func (s *MatchingService) calculateExperienceBonus(employee models.Employee, extractedSkills map[string]interface{}) int {
	bonus := 0

	// Check if employee level matches requirements
	if strings.Contains(strings.ToLower(employee.Level), "senior") {
		bonus += 10
	} else if strings.Contains(strings.ToLower(employee.Level), "lead") {
		bonus += 15
	}

	// Check years of experience from extracted data
	if employee.ExtractedData != nil {
		if years, ok := employee.ExtractedData["years_experience"].(float64); ok {
			if years >= 5 {
				bonus += 10
			} else if years >= 3 {
				bonus += 5
			}
		}
	}

	return bonus
}

// determineMatchLevel determines the match level based on score
func (s *MatchingService) determineMatchLevel(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 75:
		return "Very Good"
	case score >= 60:
		return "Good"
	case score >= 40:
		return "Fair"
	default:
		return "Poor"
	}
}

// Helper functions
func (s *MatchingService) containsSkill(skills []string, skill string) bool {
	for _, s := range skills {
		if s == skill {
			return true
		}
	}
	return false
}

func (s *MatchingService) removeDuplicates(skills []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, skill := range skills {
		if !keys[skill] {
			keys[skill] = true
			result = append(result, skill)
		}
	}

	return result
}

func (s *MatchingService) sortMatchesByScore(matches []models.MatchingCandidate) {
	// Simple bubble sort for matches by score (highest first)
	n := len(matches)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if matches[j].MatchScore < matches[j+1].MatchScore {
				matches[j], matches[j+1] = matches[j+1], matches[j]
			}
		}
	}
}
