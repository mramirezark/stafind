package services

import (
	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type dashboardService struct {
	employeeRepo repositories.EmployeeRepository
	skillRepo    repositories.SkillRepository
	aiAgentRepo  repositories.AIAgentRepository
	matchRepo    repositories.MatchRepository
}

func NewDashboardService(employeeRepo repositories.EmployeeRepository, skillRepo repositories.SkillRepository, aiAgentRepo repositories.AIAgentRepository, matchRepo repositories.MatchRepository) DashboardService {
	return &dashboardService{
		employeeRepo: employeeRepo,
		skillRepo:    skillRepo,
		aiAgentRepo:  aiAgentRepo,
		matchRepo:    matchRepo,
	}
}

// GetDashboardStats returns dashboard statistics
func (s *dashboardService) GetDashboardStats() (*models.DashboardStats, error) {
	// Get total employees
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	totalEmployees := len(employees)

	// Get AI agent requests
	requests, err := s.aiAgentRepo.GetAll(constants.MaxPageSize, constants.DefaultOffset) // Get all requests for stats
	if err != nil {
		return nil, err
	}

	// Count requests by status
	totalRequests := len(requests)
	completedRequests := 0
	pendingRequests := 0

	for _, req := range requests {
		switch req.Status {
		case "completed":
			completedRequests++
		case "pending", "processing":
			pendingRequests++
		}
	}

	return &models.DashboardStats{
		TotalEmployees:    totalEmployees,
		TotalRequests:     totalRequests,
		CompletedRequests: completedRequests,
		PendingRequests:   pendingRequests,
	}, nil
}

// GetRecentEmployees returns recent employees
func (s *dashboardService) GetRecentEmployees(limit int) ([]models.Employee, error) {
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Sort by created_at descending and limit
	for i := 0; i < len(employees)-1; i++ {
		for j := i + 1; j < len(employees); j++ {
			if employees[i].CreatedAt.Before(employees[j].CreatedAt) {
				employees[i], employees[j] = employees[j], employees[i]
			}
		}
	}

	if limit > len(employees) {
		limit = len(employees)
	}

	return employees[:limit], nil
}

// GetDepartmentStats returns department statistics
func (s *dashboardService) GetDepartmentStats() ([]models.DepartmentStats, error) {
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Count employees by department
	departmentCount := make(map[string]int)
	for _, emp := range employees {
		if emp.Department != "" {
			departmentCount[emp.Department]++
		}
	}

	// Convert to slice
	var stats []models.DepartmentStats
	for dept, count := range departmentCount {
		stats = append(stats, models.DepartmentStats{
			Department: dept,
			Count:      count,
		})
	}

	return stats, nil
}

// GetSkillDemandStats returns skill demand statistics
func (s *dashboardService) GetSkillDemandStats() ([]models.SkillDemandStats, error) {
	// Get all AI agent requests
	requests, err := s.aiAgentRepo.GetAll(constants.MaxPageSize, constants.DefaultOffset)
	if err != nil {
		return nil, err
	}

	// Count skills from extracted skills in requests
	skillCount := make(map[string]int)
	for _, req := range requests {
		if req.ExtractedSkills != nil {
			for _, skill := range req.ExtractedSkills {
				skillCount[skill]++
			}
		}
	}

	// Convert to slice and sort by count
	var stats []models.SkillDemandStats
	for skillName, count := range skillCount {
		// Get skill category
		skills, err := s.skillRepo.GetAll()
		if err != nil {
			continue
		}

		category := "Other"
		for _, skill := range skills {
			if skill.Name == skillName {
				category = skill.Category
				break
			}
		}

		stats = append(stats, models.SkillDemandStats{
			SkillName: skillName,
			Count:     count,
			Category:  category,
		})
	}

	// Sort by count (descending)
	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].Count < stats[j].Count {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	// Return top 10
	if len(stats) > 10 {
		stats = stats[:10]
	}

	return stats, nil
}

// GetTopSuggestedEmployees returns top suggested employees based on match frequency and scores
func (s *dashboardService) GetTopSuggestedEmployees(limit int) ([]models.TopSuggestedEmployee, error) {
	// Get all matches
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Group matches by employee and calculate stats
	employeeStats := make(map[int]*models.TopSuggestedEmployee)

	for _, match := range matches {
		if match.Employee.ID == 0 {
			continue // Skip if employee not loaded
		}

		if emp, exists := employeeStats[match.Employee.ID]; exists {
			emp.MatchCount++
			emp.AvgMatchScore = (emp.AvgMatchScore*float64(emp.MatchCount-1) + match.MatchScore) / float64(emp.MatchCount)
		} else {
			employeeStats[match.Employee.ID] = &models.TopSuggestedEmployee{
				EmployeeID:     match.Employee.ID,
				EmployeeName:   match.Employee.Name,
				EmployeeEmail:  match.Employee.Email,
				Department:     match.Employee.Department,
				Level:          match.Employee.Level,
				Location:       match.Employee.Location,
				CurrentProject: match.Employee.CurrentProject,
				MatchCount:     1,
				AvgMatchScore:  match.MatchScore,
			}
		}
	}

	// Convert to slice and sort by match count and average score
	var topEmployees []models.TopSuggestedEmployee
	for _, emp := range employeeStats {
		topEmployees = append(topEmployees, *emp)
	}

	// Sort by match count (descending), then by average score (descending)
	for i := 0; i < len(topEmployees)-1; i++ {
		for j := i + 1; j < len(topEmployees); j++ {
			if topEmployees[i].MatchCount < topEmployees[j].MatchCount ||
				(topEmployees[i].MatchCount == topEmployees[j].MatchCount &&
					topEmployees[i].AvgMatchScore < topEmployees[j].AvgMatchScore) {
				topEmployees[i], topEmployees[j] = topEmployees[j], topEmployees[i]
			}
		}
	}

	// Limit results
	if limit > len(topEmployees) {
		limit = len(topEmployees)
	}

	return topEmployees[:limit], nil
}

// GetDashboardMetrics returns comprehensive dashboard metrics
func (s *dashboardService) GetDashboardMetrics() (*models.DashboardMetrics, error) {
	// Get basic stats
	stats, err := s.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	// Get most requested skills
	mostRequestedSkills, err := s.GetSkillDemandStats()
	if err != nil {
		return nil, err
	}

	// Get top suggested employees
	topSuggestedEmployees, err := s.GetTopSuggestedEmployees(5)
	if err != nil {
		return nil, err
	}

	// Get recent requests
	recentRequests, err := s.aiAgentRepo.GetAll(10, 0)
	if err != nil {
		return nil, err
	}

	return &models.DashboardMetrics{
		Stats:                 *stats,
		MostRequestedSkills:   mostRequestedSkills,
		TopSuggestedEmployees: topSuggestedEmployees,
		RecentRequests:        recentRequests,
	}, nil
}
