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
	Create(match *models.Match) (*models.Match, error)
	GetByEmployeeID(employeeID int) ([]models.Match, error)
	GetAll() ([]models.Match, error)
	Delete(id int) error
}

// AIAgentRepository defines the interface for AI agent data operations
type AIAgentRepository interface {
	Create(req *models.AIAgentRequest) (*models.AIAgentRequest, error)
	GetByID(id int) (*models.AIAgentRequest, error)
	GetByTeamsMessageID(teamsMessageID string) (*models.AIAgentRequest, error)
	Update(id int, req *models.AIAgentRequest) error
	GetAll(limit int, offset int) ([]models.AIAgentRequest, error)
	SaveResponse(response *models.AIAgentResponse) error
	GetResponseByRequestID(requestID int) (*models.AIAgentResponse, error)
}

// APIKeyRepository defines the interface for API key data operations
type APIKeyRepository interface {
	Create(key *models.APIKey) (*models.APIKey, error)
	GetByID(id int) (*models.APIKey, error)
	GetByHash(hash string) (*models.APIKey, error)
	GetAll(limit, offset int) ([]models.APIKey, error)
	Update(id int, key *models.APIKey) error
	Deactivate(id int) error
	UpdateLastUsed(hash string) error
	Delete(id int) error
}
