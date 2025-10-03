package matching

import (
	"sort"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"
	"strings"
)

// MatchEngine handles the matching logic between job requests and employees
type MatchEngine struct{}

// NewMatchEngine creates a new match engine
func NewMatchEngine() *MatchEngine {
	return &MatchEngine{}
}

// SearchEmployees searches for employees based on search criteria
func (me *MatchEngine) SearchEmployees(searchReq *models.SearchRequest, employees []models.Employee) []models.Match {
	var matches []models.Match

	for _, employee := range employees {
		score, matchingSkills := me.calculateMatchScore(searchReq, &employee)

		if score > 0 {
			match := models.Match{
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

// calculateMatchScore calculates the match score between search criteria and employee
func (me *MatchEngine) calculateMatchScore(searchReq *models.SearchRequest, employee *models.Employee) (float64, []string) {
	var totalScore float64
	var matchingSkills []string

	// Convert employee skills to a map for easy lookup
	employeeSkillMap := make(map[string]models.EmployeeSkill)
	for _, skill := range employee.Skills {
		employeeSkillMap[skill.Name] = models.EmployeeSkill{
			Skill:            skill,
			ProficiencyLevel: constants.DefaultProficiencyLevel, // Default proficiency level
			YearsExperience:  constants.DefaultYearsExperience,  // Default years of experience
		}
	}

	// Score required skills (higher weight)
	requiredScore := me.scoreSkills(searchReq.RequiredSkills, employeeSkillMap, constants.RequiredSkillsWeight, &matchingSkills)

	// Score preferred skills (lower weight)
	preferredScore := me.scoreSkills(searchReq.PreferredSkills, employeeSkillMap, constants.PreferredSkillsWeight, &matchingSkills)

	// Department match bonus
	departmentBonus := me.calculateDepartmentBonus(searchReq.Department, employee.Department)

	// Experience level bonus
	experienceBonus := me.calculateExperienceBonus(searchReq.ExperienceLevel, employee.Level)

	// Location bonus
	locationBonus := me.calculateLocationBonus(searchReq.Location, employee.Location)

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
		if employeeSkill, matchedSkillName := me.findMatchingSkill(skillName, employeeSkillMap); employeeSkill != nil {
			// Base score for having the skill
			skillScore := weight * constants.BaseSkillScoreMultiplier

			// Add proficiency bonus (1-5 scale)
			proficiencyBonus := float64(employeeSkill.ProficiencyLevel) * constants.ProficiencyBonusMultiplier

			// Add experience bonus
			experienceBonus := employeeSkill.YearsExperience * constants.ExperienceBonusMultiplier

			totalScore += skillScore + proficiencyBonus + experienceBonus
			matchedCount++
			*matchingSkills = append(*matchingSkills, matchedSkillName)
		}
	}

	// Calculate coverage percentage
	coverage := float64(matchedCount) / float64(len(skillNames))

	// Apply coverage multiplier (encourage higher coverage)
	return totalScore * (constants.CoverageBaseMultiplier + constants.CoverageBonusMultiplier*coverage)
}

// findMatchingSkill finds a matching skill with case-insensitive and normalized matching
func (me *MatchEngine) findMatchingSkill(skillName string, employeeSkillMap map[string]models.EmployeeSkill) (*models.EmployeeSkill, string) {
	// Normalize the skill name for matching
	normalizedSkillName := me.normalizeSkillName(skillName)

	// Try exact match first (case-insensitive)
	for dbSkillName, employeeSkill := range employeeSkillMap {
		if strings.EqualFold(dbSkillName, skillName) {
			return &employeeSkill, dbSkillName
		}
	}

	// Try normalized match
	for dbSkillName, employeeSkill := range employeeSkillMap {
		if strings.EqualFold(me.normalizeSkillName(dbSkillName), normalizedSkillName) {
			return &employeeSkill, dbSkillName
		}
	}

	// Try partial match for common abbreviations
	for dbSkillName, employeeSkill := range employeeSkillMap {
		if me.isSkillAbbreviation(skillName, dbSkillName) {
			return &employeeSkill, dbSkillName
		}
	}

	// Try partial word matching for complex skill names (e.g., "Java Development" matches "Java")
	for dbSkillName, employeeSkill := range employeeSkillMap {
		if me.isPartialSkillMatch(skillName, dbSkillName) {
			return &employeeSkill, dbSkillName
		}
	}

	return nil, ""
}

// normalizeSkillName normalizes skill names for better matching
func (me *MatchEngine) normalizeSkillName(skillName string) string {
	normalized := strings.ToLower(strings.TrimSpace(skillName))

	// Common skill name normalizations
	skillMappings := constants.SkillNormalizationMap

	if mapped, exists := skillMappings[normalized]; exists {
		return mapped
	}

	return normalized
}

// isSkillAbbreviation checks if one skill is an abbreviation of another
func (me *MatchEngine) isSkillAbbreviation(abbr, fullName string) bool {
	abbr = strings.ToLower(strings.TrimSpace(abbr))
	fullName = strings.ToLower(strings.TrimSpace(fullName))

	// Common abbreviation mappings
	abbreviationMappings := constants.SkillAbbreviationMap

	if possibleMatches, exists := abbreviationMappings[abbr]; exists {
		for _, possible := range possibleMatches {
			if strings.Contains(fullName, possible) {
				return true
			}
		}
	}

	return false
}

// isPartialSkillMatch checks if a skill name contains another skill as a word
func (me *MatchEngine) isPartialSkillMatch(searchSkill, employeeSkill string) bool {
	searchSkill = strings.ToLower(strings.TrimSpace(searchSkill))
	employeeSkill = strings.ToLower(strings.TrimSpace(employeeSkill))

	// Split search skill into words
	searchWords := strings.Fields(searchSkill)
	employeeWords := strings.Fields(employeeSkill)

	// Check if any word from search skill matches any word from employee skill
	for _, searchWord := range searchWords {
		for _, employeeWord := range employeeWords {
			if searchWord == employeeWord {
				return true
			}
		}
	}

	return false
}

// calculateDepartmentBonus calculates bonus for department match
func (me *MatchEngine) calculateDepartmentBonus(jobDept, employeeDept string) float64 {
	if jobDept == "" || employeeDept == "" {
		return 0
	}
	if jobDept == employeeDept {
		return constants.DepartmentMatchBonus
	}
	return 0
}

// calculateExperienceBonus calculates bonus for experience level match
func (me *MatchEngine) calculateExperienceBonus(jobLevel, employeeLevel string) float64 {
	if jobLevel == "" || employeeLevel == "" {
		return 0
	}

	levelMap := constants.ExperienceLevelMap

	jobLevelNum := levelMap[jobLevel]
	employeeLevelNum := levelMap[employeeLevel]

	if employeeLevelNum >= jobLevelNum {
		// Employee meets or exceeds required level
		return constants.ExperienceLevelMatchBonus
	} else {
		// Employee is below required level, but still gets some points
		return float64(employeeLevelNum) / float64(jobLevelNum) * constants.ExperienceLevelPartialMultiplier
	}
}

// calculateLocationBonus calculates bonus for location match
func (me *MatchEngine) calculateLocationBonus(jobLocation, employeeLocation string) float64 {
	if jobLocation == "" || employeeLocation == "" {
		return 0
	}
	if jobLocation == employeeLocation {
		return constants.LocationMatchBonus
	}
	return 0
}
