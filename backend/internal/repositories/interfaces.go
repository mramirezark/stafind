package repositories

import (
	"stafind-backend/internal/models"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	GetAll() ([]models.Employee, error)
	GetByID(id int) (*models.Employee, error)
	GetByEmail(email string) (*models.Employee, error)
	Create(req *models.CreateEmployeeRequest) (*models.Employee, error)
	CreateWithExtraction(req *models.CreateEmployeeRequest, originalText string, extractedData map[string]interface{}, extractionSource, extractionStatus, resumeURL string) (*models.Employee, error)
	Update(id int, req *models.CreateEmployeeRequest) (*models.Employee, error)
	UpdateWithExtraction(id int, req *models.CreateEmployeeRequest, originalText string, extractedData map[string]interface{}, extractionSource, extractionStatus, resumeURL string) (*models.Employee, error)
	Delete(id int) error
	GetSkills(employeeID int) ([]models.Skill, error)
	AddSkill(employeeID int, skillReq *models.EmployeeSkillReq) error
	RemoveSkills(employeeID int) error
	GetEmployeesWithSkills(skillNames []string) ([]models.Employee, error)
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	GetByName(name string) (*models.Category, error)
	GetCategoriesWithSkillCount() ([]models.CategoryWithSkillCount, error)
	GetSkillsByCategoryID(categoryID int) ([]models.Skill, error)
	Create(category *models.Category) (*models.Category, error)
	CreateBatch(categories []models.Category) ([]models.Category, error)
	Update(id int, category *models.Category) (*models.Category, error)
	Delete(id int) error
	DeleteBatch(ids []int) error
	GetCategoryStats() (*models.SkillStats, error)
}

// SkillRepository defines the interface for skill data operations
type SkillRepository interface {
	GetAll() ([]models.Skill, error)
	GetByID(id int) (*models.Skill, error)
	GetByName(name string) (*models.Skill, error)
	GetByCategoryID(categoryID int) ([]models.Skill, error)
	GetSkillsWithCategories() ([]models.Skill, error)
	SearchSkills(query string) ([]models.Skill, error)
	GetPopularSkills(limit int) ([]models.Skill, error)
	GetSkillsByIDs(ids []int) ([]models.Skill, error)
	GetSkillsWithEmployeeCount() ([]models.SkillWithCount, error)
	Create(skill *models.Skill) (*models.Skill, error)
	CreateBatch(skills []models.Skill) ([]models.Skill, error)
	Update(id int, skill *models.Skill) (*models.Skill, error)
	UpdateBatch(updates []models.SkillUpdate) error
	Delete(id int) error
	DeleteBatch(ids []int) error
	GetSkillStats() (*models.SkillStats, error)
	AddSkillToCategory(skillID, categoryID int) error
	RemoveSkillFromCategory(skillID, categoryID int) error
	GetSkillCategories(skillID int) ([]models.Category, error)
	AssociateCategories(skillID int, categoryIDs []int) error
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
	UpdateStatus(id int, status string) error
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
