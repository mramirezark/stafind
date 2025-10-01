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

// SearchService defines the interface for search business logic
type SearchService interface {
	SearchEmployees(searchReq *models.SearchRequest) ([]models.Match, error)
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
	GetRecentEmployees(limit int) ([]models.Employee, error)
	GetDepartmentStats() ([]models.DepartmentStats, error)
	GetSkillDemandStats() ([]models.SkillDemandStats, error)
	GetTopSuggestedEmployees(limit int) ([]models.TopSuggestedEmployee, error)
	GetDashboardMetrics() (*models.DashboardMetrics, error)
}

// AIAgentService defines the interface for AI agent business logic
type AIAgentService interface {
	CreateAIAgentRequest(req *models.CreateAIAgentRequest) (*models.AIAgentRequest, error)
	GetAIAgentRequest(id int) (*models.AIAgentRequest, error)
	UpdateAIAgentRequest(id int, req *models.AIAgentRequest) error
	ProcessAIAgentRequest(id int) (*models.AIAgentResponse, error)
	ExtractSkillsFromText(text string) (*models.SkillExtractionResponse, error)
	GetAIAgentRequests(limit int, offset int) ([]models.AIAgentRequest, error)
	GetAIAgentResponse(requestID int) (*models.AIAgentResponse, error)
}

// NotificationService defines the interface for notification business logic
type NotificationService interface {
	SendTeamsMessage(channelID string, message string) error
	SendAdminEmail(subject string, body string) error
	LogError(requestID int, error string) error
}

// APIKeyService defines the interface for API key business logic
type APIKeyService interface {
	CreateAPIKey(req *models.CreateAPIKeyRequest) (*models.APIKeyResponse, error)
	ValidateAPIKey(key string) (*models.APIKey, error)
	GetAPIKeys(limit, offset int) ([]models.APIKey, error)
	DeactivateAPIKey(id int) error
	UpdateLastUsed(key string) error
	RotateAPIKey(oldKeyID int) (*models.APIKeyResponse, error)
}

// LlamaAIService defines the interface for Llama AI processing
type LlamaAIService interface {
	ProcessText(request *models.LlamaAIRequest) (*models.LlamaAIResponse, error)
	GetHealthStatus() (string, error)
}
