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

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id int) (*models.Category, error)
	GetCategoryByName(name string) (*models.Category, error)
	GetCategoriesWithSkillCount() ([]models.CategoryWithSkillCount, error)
	GetSkillsByCategoryID(categoryID int) ([]models.Skill, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	CreateCategoriesBatch(categories []models.Category) ([]models.Category, error)
	UpdateCategory(id int, category *models.Category) (*models.Category, error)
	DeleteCategory(id int) error
	DeleteCategoriesBatch(ids []int) error
	GetCategoryStats() (*models.SkillStats, error)
	ValidateCategory(category *models.Category) error
}

// SkillService defines the interface for skill business logic
type SkillService interface {
	GetAllSkills() ([]models.Skill, error)
	GetSkillByID(id int) (*models.Skill, error)
	GetSkillsByCategoryID(categoryID int) ([]models.Skill, error)
	GetSkillsWithCategories() ([]models.Skill, error)
	SearchSkills(query string) ([]models.Skill, error)
	GetPopularSkills(limit int) ([]models.Skill, error)
	GetSkillsByIDs(ids []int) ([]models.Skill, error)
	GetSkillsWithEmployeeCount() ([]models.SkillWithCount, error)
	CreateSkill(skill *models.Skill) (*models.Skill, error)
	CreateSkillsBatch(skills []models.Skill) ([]models.Skill, error)
	UpdateSkill(id int, skill *models.Skill) (*models.Skill, error)
	UpdateSkillsBatch(updates []models.SkillUpdate) error
	DeleteSkill(id int) error
	DeleteSkillsBatch(ids []int) error
	GetSkillStats() (*models.SkillStats, error)
	ValidateSkill(skill *models.Skill) error
	AddSkillToCategory(skillID, categoryID int) error
	RemoveSkillFromCategory(skillID, categoryID int) error
	GetSkillCategories(skillID int) ([]models.Category, error)
	GetSkillsByEmployeeID(employeeID int) ([]models.Skill, error)
	GetSkillsByEmployeeIDs(employeeIDs []int) (map[int][]models.Skill, error)
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
	GetAIAgentRequestByTeamsMessageID(teamsMessageID string) (*models.AIAgentRequest, error)
	UpdateAIAgentRequest(id int, req *models.AIAgentRequest) error
	UpdateAIAgentStatus(id int, status string) error
	ProcessAIAgentRequest(id int) (*models.AIAgentResponse, error)
	ExtractSkillsFromText(text string) (*models.SkillExtractResponse, error)
	GetAIAgentRequests(limit int, offset int) ([]models.AIAgentRequest, error)
	GetAIAgentResponse(requestID int) (*models.AIAgentResponse, error)
	FindMatchingEmployees(skills []string) ([]models.Match, error)
	GenerateMatchExplanations(matches []models.Match, extractedSkills []string) []models.AIAgentMatch
	GenerateMatchSummary(matches []models.AIAgentMatch, extractedSkills []string) string
	SaveMatch(match *models.Match) (*models.Match, error)
	SaveResponse(response *models.AIAgentResponse) error
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

// ExtractionService defines the interface for candidate and resume extraction using NER
type ExtractionService interface {
	ProcessText(request *models.ExtractAIRequest) (*models.ExtractAIResponse, error)
	GetHealthStatus() (string, error)
}

// ModernNERService defines the interface for modern NER-based skill extraction
type ModernNERService interface {
	ExtractSkillsFromText(text string) (*SkillExtractionResult, error)
	SetAPIKeys(openAIKey, huggingFaceToken, spacyEndpoint string)
	GetSupportedMethods() []string
	GetMethodConfidence(method string) float64
}
