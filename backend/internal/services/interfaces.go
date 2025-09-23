package services

import (
	"stafind-backend/internal/models"
)

// EmployeeService defines the interface for employee business logic
type EmployeeService interface {
	GetAllEmployees() ([]models.Employee, error)
	GetEmployeeByID(id int) (*models.Employee, error)
	CreateEmployee(req *models.CreateEmployeeRequest) (*models.Employee, error)
	UpdateEmployee(id int, req *models.CreateEmployeeRequest) (*models.Employee, error)
	DeleteEmployee(id int) error
}

// JobRequestService defines the interface for job request business logic
type JobRequestService interface {
	GetAllJobRequests() ([]models.JobRequest, error)
	GetJobRequestByID(id int) (*models.JobRequest, error)
	CreateJobRequest(req *models.CreateJobRequestRequest) (*models.JobRequest, error)
	UpdateJobRequest(id int, req *models.CreateJobRequestRequest) (*models.JobRequest, error)
	DeleteJobRequest(id int) error
}

// SearchService defines the interface for search business logic
type SearchService interface {
	SearchEmployees(searchReq *models.SearchRequest) ([]models.Match, error)
	FindMatchesForJobRequest(jobRequestID int) ([]models.Match, error)
}

// SkillService defines the interface for skill business logic
type SkillService interface {
	GetAllSkills() ([]models.Skill, error)
	GetSkillByID(id int) (*models.Skill, error)
	CreateSkill(skill *models.Skill) (*models.Skill, error)
	UpdateSkill(id int, skill *models.Skill) (*models.Skill, error)
	DeleteSkill(id int) error
}

// DashboardService defines the interface for dashboard business logic
type DashboardService interface {
	GetDashboardStats() (*models.DashboardStats, error)
	GetRecentJobRequests(limit int) ([]models.JobRequest, error)
	GetRecentEmployees(limit int) ([]models.Employee, error)
	GetDepartmentStats() ([]models.DepartmentStats, error)
	GetSkillDemandStats() ([]models.SkillDemandStats, error)
}
