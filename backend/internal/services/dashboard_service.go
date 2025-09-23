package services

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"time"
)

type dashboardService struct {
	employeeRepo   repositories.EmployeeRepository
	jobRequestRepo repositories.JobRequestRepository
	skillRepo      repositories.SkillRepository
}

func NewDashboardService(employeeRepo repositories.EmployeeRepository, jobRequestRepo repositories.JobRequestRepository, skillRepo repositories.SkillRepository) DashboardService {
	return &dashboardService{
		employeeRepo:   employeeRepo,
		jobRequestRepo: jobRequestRepo,
		skillRepo:      skillRepo,
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

	// Get total job requests
	jobRequests, err := s.jobRequestRepo.GetAll()
	if err != nil {
		return nil, err
	}
	totalJobRequests := len(jobRequests)

	// Count active matches (job requests with status 'open')
	activeMatches := 0
	for _, jr := range jobRequests {
		if jr.Status == "open" {
			activeMatches++
		}
	}

	// Count recent requests (last 7 days)
	weekAgo := time.Now().AddDate(0, 0, -7)
	recentRequests := 0
	for _, jr := range jobRequests {
		if jr.CreatedAt.After(weekAgo) {
			recentRequests++
		}
	}

	return &models.DashboardStats{
		TotalEmployees:   totalEmployees,
		TotalJobRequests: totalJobRequests,
		ActiveMatches:    activeMatches,
		RecentRequests:   recentRequests,
	}, nil
}

// GetRecentJobRequests returns recent job requests
func (s *dashboardService) GetRecentJobRequests(limit int) ([]models.JobRequest, error) {
	jobRequests, err := s.jobRequestRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Sort by created_at descending and limit
	// Simple sort implementation
	for i := 0; i < len(jobRequests)-1; i++ {
		for j := i + 1; j < len(jobRequests); j++ {
			if jobRequests[i].CreatedAt.Before(jobRequests[j].CreatedAt) {
				jobRequests[i], jobRequests[j] = jobRequests[j], jobRequests[i]
			}
		}
	}

	if limit > len(jobRequests) {
		limit = len(jobRequests)
	}

	return jobRequests[:limit], nil
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
	jobRequests, err := s.jobRequestRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Count skill demand from job requests
	skillCount := make(map[string]int)
	skillCategory := make(map[string]string)

	for _, jr := range jobRequests {
		// Count required skills
		for _, skill := range jr.RequiredSkills {
			skillCount[skill]++
		}
		// Count preferred skills
		for _, skill := range jr.PreferredSkills {
			skillCount[skill]++
		}
	}

	// Get skill categories from skills table
	skills, err := s.skillRepo.GetAll()
	if err == nil {
		for _, skill := range skills {
			skillCategory[skill.Name] = skill.Category
		}
	}

	// Convert to slice
	var stats []models.SkillDemandStats
	for skill, count := range skillCount {
		category := skillCategory[skill]
		if category == "" {
			category = "Other"
		}
		stats = append(stats, models.SkillDemandStats{
			SkillName: skill,
			Count:     count,
			Category:  category,
		})
	}

	return stats, nil
}
