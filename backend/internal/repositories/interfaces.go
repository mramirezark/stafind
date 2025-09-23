package repositories

import (
	"stafind-backend/internal/models"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	GetAll() ([]models.Employee, error)
	GetByID(id int) (*models.Employee, error)
	Create(req *models.CreateEmployeeRequest) (*models.Employee, error)
	Update(id int, req *models.CreateEmployeeRequest) (*models.Employee, error)
	Delete(id int) error
	GetSkills(employeeID int) ([]models.Skill, error)
	AddSkill(employeeID int, skillReq *models.EmployeeSkillReq) error
	RemoveSkills(employeeID int) error
}

// JobRequestRepository defines the interface for job request data operations
type JobRequestRepository interface {
	GetAll() ([]models.JobRequest, error)
	GetByID(id int) (*models.JobRequest, error)
	Create(req *models.CreateJobRequestRequest) (*models.JobRequest, error)
	Update(id int, req *models.CreateJobRequestRequest) (*models.JobRequest, error)
	Delete(id int) error
}

// SkillRepository defines the interface for skill data operations
type SkillRepository interface {
	GetAll() ([]models.Skill, error)
	GetByID(id int) (*models.Skill, error)
	GetByName(name string) (*models.Skill, error)
	Create(skill *models.Skill) (*models.Skill, error)
	Update(id int, skill *models.Skill) (*models.Skill, error)
	Delete(id int) error
}

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	GetByJobRequestID(jobRequestID int) ([]models.Match, error)
	Create(match *models.Match) (*models.Match, error)
	DeleteByJobRequestID(jobRequestID int) error
}
