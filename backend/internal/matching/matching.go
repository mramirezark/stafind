package matching

import (
	"sort"
	"stafind-backend/internal/models"
)

// MatchEngine handles the matching logic between job requests and employees
type MatchEngine struct{}

// NewMatchEngine creates a new match engine
func NewMatchEngine() *MatchEngine {
	return &MatchEngine{}
}

// FindMatches finds the best matching employees for a job request
func (me *MatchEngine) FindMatches(jobRequest *models.JobRequest, employees []models.Employee) []models.Match {
	var matches []models.Match

	for _, employee := range employees {
		score, matchingSkills := me.calculateMatchScore(jobRequest, &employee)

		if score > 0 {
			match := models.Match{
				JobRequestID:   jobRequest.ID,
				EmployeeID:     employee.ID,
				MatchScore:     score,
				MatchingSkills: matchingSkills,
				Employee:       employee,
			}
			matches = append(matches, match)
		}
	}

	// Sort matches by score (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchScore > matches[j].MatchScore
	})

	return matches
}

// SearchEmployees searches for employees based on search criteria
func (me *MatchEngine) SearchEmployees(searchReq *models.SearchRequest, employees []models.Employee) []models.Match {
	// Create a temporary job request for matching
	tempJobRequest := &models.JobRequest{
		RequiredSkills:  searchReq.RequiredSkills,
		PreferredSkills: searchReq.PreferredSkills,
		Department:      searchReq.Department,
		ExperienceLevel: searchReq.ExperienceLevel,
		Location:        searchReq.Location,
	}

	// Find matches using the existing logic
	matches := me.FindMatches(tempJobRequest, employees)

	// Filter by minimum match score if specified
	if searchReq.MinMatchScore > 0 {
		var filteredMatches []models.Match
		for _, match := range matches {
			if match.MatchScore >= searchReq.MinMatchScore {
				filteredMatches = append(filteredMatches, match)
			}
		}
		return filteredMatches
	}

	return matches
}

// calculateMatchScore calculates the match score between a job request and employee
func (me *MatchEngine) calculateMatchScore(jobRequest *models.JobRequest, employee *models.Employee) (float64, []string) {
	var totalScore float64
	var matchingSkills []string

	// Convert employee skills to a map for easy lookup
	employeeSkillMap := make(map[string]models.EmployeeSkill)
	for _, skill := range employee.Skills {
		employeeSkillMap[skill.Name] = models.EmployeeSkill{
			Skill:            skill,
			ProficiencyLevel: 3, // Default proficiency level
			YearsExperience:  2, // Default years of experience
		}
	}

	// Score required skills (higher weight)
	requiredScore := me.scoreSkills(jobRequest.RequiredSkills, employeeSkillMap, 3.0, &matchingSkills)

	// Score preferred skills (lower weight)
	preferredScore := me.scoreSkills(jobRequest.PreferredSkills, employeeSkillMap, 1.0, &matchingSkills)

	// Department match bonus
	departmentBonus := me.calculateDepartmentBonus(jobRequest.Department, employee.Department)

	// Experience level bonus
	experienceBonus := me.calculateExperienceBonus(jobRequest.ExperienceLevel, employee.Level)

	// Location bonus
	locationBonus := me.calculateLocationBonus(jobRequest.Location, employee.Location)

	totalScore = requiredScore + preferredScore + departmentBonus + experienceBonus + locationBonus

	return totalScore, matchingSkills
}

// scoreSkills scores a list of skills against employee's skills
func (me *MatchEngine) scoreSkills(skillNames []string, employeeSkillMap map[string]models.EmployeeSkill, weight float64, matchingSkills *[]string) float64 {
	if len(skillNames) == 0 {
		return 0
	}

	var totalScore float64
	var matchedCount int

	for _, skillName := range skillNames {
		if employeeSkill, exists := employeeSkillMap[skillName]; exists {
			// Base score for having the skill
			skillScore := weight * 2.0

			// Add proficiency bonus (1-5 scale)
			proficiencyBonus := float64(employeeSkill.ProficiencyLevel) * 0.5

			// Add experience bonus
			experienceBonus := employeeSkill.YearsExperience * 0.1

			totalScore += skillScore + proficiencyBonus + experienceBonus
			matchedCount++
			*matchingSkills = append(*matchingSkills, skillName)
		}
	}

	// Calculate coverage percentage
	coverage := float64(matchedCount) / float64(len(skillNames))

	// Apply coverage multiplier (encourage higher coverage)
	return totalScore * (0.5 + 0.5*coverage)
}

// calculateDepartmentBonus calculates bonus for department match
func (me *MatchEngine) calculateDepartmentBonus(jobDept, employeeDept string) float64 {
	if jobDept == "" || employeeDept == "" {
		return 0
	}
	if jobDept == employeeDept {
		return 2.0
	}
	return 0
}

// calculateExperienceBonus calculates bonus for experience level match
func (me *MatchEngine) calculateExperienceBonus(jobLevel, employeeLevel string) float64 {
	if jobLevel == "" || employeeLevel == "" {
		return 0
	}

	levelMap := map[string]int{
		"junior":    1,
		"mid":       2,
		"senior":    3,
		"staff":     4,
		"principal": 5,
	}

	jobLevelNum := levelMap[jobLevel]
	employeeLevelNum := levelMap[employeeLevel]

	if employeeLevelNum >= jobLevelNum {
		// Employee meets or exceeds required level
		return 1.5
	} else {
		// Employee is below required level, but still gets some points
		return float64(employeeLevelNum) / float64(jobLevelNum) * 1.0
	}
}

// calculateLocationBonus calculates bonus for location match
func (me *MatchEngine) calculateLocationBonus(jobLocation, employeeLocation string) float64 {
	if jobLocation == "" || employeeLocation == "" {
		return 0
	}
	if jobLocation == employeeLocation {
		return 1.0
	}
	return 0
}
