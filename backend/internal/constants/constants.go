package constants

// HTTP Status Codes
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// Error Codes
const (
	// API Key related errors
	ErrorCodeMissingAPIKey       = "MISSING_API_KEY"
	ErrorCodeInvalidAPIKey       = "INVALID_API_KEY"
	ErrorCodeMissingAuthHeader   = "MISSING_AUTH_HEADER"
	ErrorCodeInvalidAuthFormat   = "INVALID_AUTH_FORMAT"
	ErrorCodeInvalidServiceToken = "INVALID_SERVICE_TOKEN"

	// Validation errors
	ErrorCodeInvalidID        = "INVALID_ID"
	ErrorCodeValidationFailed = "VALIDATION_FAILED"
	ErrorCodeRequiredField    = "REQUIRED_FIELD"

	// Resource errors
	ErrorCodeNotFound      = "NOT_FOUND"
	ErrorCodeAlreadyExists = "ALREADY_EXISTS"
	ErrorCodeConflict      = "CONFLICT"
)

// Error Messages
const (
	// API Key messages
	MsgAPIKeyRequired      = "API key required"
	MsgInvalidAPIKey       = "Invalid API key"
	MsgAuthHeaderRequired  = "Authorization header required"
	MsgInvalidAuthFormat   = "Invalid authorization format"
	MsgInvalidServiceToken = "Invalid service token"

	// Validation messages
	MsgInvalidID        = "Invalid ID"
	MsgInvalidAPIKeyID  = "Invalid API key ID"
	MsgValidationFailed = "Validation failed"
	MsgRequiredField    = "Required field missing"

	// Resource messages
	MsgNotFound      = "Resource not found"
	MsgAlreadyExists = "Resource already exists"
	MsgConflict      = "Resource conflict"
)

// Database Configuration
const (
	// Default database values
	DefaultDBHost     = "localhost"
	DefaultDBPort     = "5432"
	DefaultDBUser     = "postgres"
	DefaultDBPassword = "password"
	DefaultDBName     = "stafind"
	DefaultSSLMode    = "disable"

	// Database providers
	DBProviderPostgres = "postgres"
	DBProviderSupabase = "supabase"
	DefaultDBProvider  = DBProviderPostgres

	// Supabase-specific defaults
	SupabaseSSLMode             = "require"
	SupabaseDefaultPoolerMode   = "transaction" // transaction, session, or statement
	SupabaseDefaultMaxOpenConns = 25
	SupabaseDefaultMaxIdleConns = 5
	SupabaseDefaultConnMaxLife  = 300 // seconds (5 minutes)
	SupabaseDefaultConnMaxIdle  = 60  // seconds (1 minute)

	// Flyway configuration
	DefaultFlywayLocations = "./flyway_migrations"
)

// Environment Variables
const (
	EnvDBHost            = "DB_HOST"
	EnvDBPort            = "DB_PORT"
	EnvDBUser            = "DB_USER"
	EnvDBPassword        = "DB_PASSWORD"
	EnvDBName            = "DB_NAME"
	EnvDBSSLMode         = "DB_SSLMODE"
	EnvDBProvider        = "DB_PROVIDER"  // postgres or supabase
	EnvDBConnectionURL   = "DATABASE_URL" // Optional: Full connection string (Supabase format)
	EnvDBMaxOpenConns    = "DB_MAX_OPEN_CONNS"
	EnvDBMaxIdleConns    = "DB_MAX_IDLE_CONNS"
	EnvDBConnMaxLifetime = "DB_CONN_MAX_LIFETIME"
	EnvDBConnMaxIdleTime = "DB_CONN_MAX_IDLE_TIME"
	EnvSupabasePooler    = "SUPABASE_POOLER_MODE" // transaction, session, or statement
	EnvFlywayLocations   = "FLYWAY_LOCATIONS"
	EnvExternalAPIKey    = "EXTERNAL_API_KEY"
	EnvServiceToken      = "SERVICE_TOKEN"
	EnvTeamsWebhookURL   = "TEAMS_WEBHOOK_URL"
	EnvSMTPHost          = "SMTP_HOST"
	EnvSMTPPort          = "SMTP_PORT"
	EnvSMTPUser          = "SMTP_USER"
	EnvSMTPPass          = "SMTP_PASS"
	EnvAdminEmail        = "ADMIN_EMAIL"
	EnvHuggingFaceAPIKey = "HUGGINGFACE_API_KEY"
)

// Development defaults
const (
	DevAPIKey       = "dev-api-key-12345"
	DevServiceToken = "service-token-12345"
)

// Pagination defaults
const (
	DefaultPageSize = 10
	DefaultPage     = 1
	DefaultOffset   = 0
	MaxPageSize     = 1000
)

// File upload constants
const (
	DefaultFilePageSize = 20
	MaxFilePageSize     = 1000
)

// HTTP Headers
const (
	HeaderAPIKey        = "X-API-Key"
	HeaderAuthorization = "Authorization"
	HeaderBearer        = "Bearer"
)

// Context keys
const (
	ContextAPIKey       = "api_key"
	ContextAuthType     = "auth_type"
	ContextServiceToken = "service_token"
	ContextRequestID    = "request_id"
)

// NER (Named Entity Recognition) entity types
const (
	EntityPerson    = "PERSON"
	EntityOrg       = "ORG"
	EntityGPE       = "GPE"
	EntityWorkOfArt = "WORK_OF_ART"
	EntityEvent     = "EVENT"
)

// AI Agent skill categories
const (
	SkillCategoryBackend         = "Backend"
	SkillCategoryFrontend        = "Frontend"
	SkillCategoryFullStack       = "Full Stack"
	SkillCategoryMobile          = "Mobile"
	SkillCategoryDevOps          = "DevOps"
	SkillCategoryData            = "Data"
	SkillCategoryProduct         = "Product"
	SkillCategoryDesign          = "Design"
	SkillCategoryQA              = "QA"
	SkillCategoryTesting         = "Testing"
	SkillCategorySecurity        = "Security"
	SkillCategoryAI              = "AI"
	SkillCategoryMachineLearning = "Machine Learning"
	SkillCategoryCloud           = "Cloud"
)

// AI Agent skill mappings
var SkillCategoryMappings = map[string]string{
	"Engineering":      "Software Engineer",
	"Development":      "Developer",
	"DevOps":           "DevOps Engineer",
	"Data":             "Data Scientist",
	"Analytics":        "Data Analyst",
	"Product":          "Product Manager",
	"Design":           "UX/UI Designer",
	"QA":               "QA Engineer",
	"Testing":          "Test Engineer",
	"Security":         "Security Engineer",
	"Infrastructure":   "Infrastructure Engineer",
	"Backend":          "Backend Developer",
	"Frontend":         "Frontend Developer",
	"Full Stack":       "Full Stack Developer",
	"Mobile":           "Mobile Developer",
	"AI":               "AI Engineer",
	"Machine Learning": "ML Engineer",
	"Cloud":            "Cloud Engineer",
	"Architecture":     "Solution Architect",
	"Management":       "Engineering Manager",
	"Technical":        "Technical Lead",
}

// Database query placeholders
const (
	PostgresArrayFormat = "{%s}" // PostgreSQL array format
)

// File permissions
const (
	DefaultFilePermission = 0666
)

// Timeout values
const (
	DefaultTimeout = 30 // seconds
)

// Log levels
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Response messages
const (
	MsgSuccess            = "Success"
	MsgAPIKeyCreated      = "API key created successfully"
	MsgAPIKeyDeactivated  = "API key deactivated successfully"
	MsgAPIKeyRotated      = "API key rotated successfully. Save the new key - it won't be shown again!"
	MsgAPIKeyValid        = "API key is valid"
	MsgOperationCompleted = "Operation completed successfully"
)

// Query parameters
const (
	ParamLimit  = "limit"
	ParamOffset = "offset"
	ParamPage   = "page"
	ParamSize   = "size"
	ParamSort   = "sort"
	ParamOrder  = "order"
)

// Sort directions
const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

// Default sort fields
const (
	SortByCreatedAt = "created_at"
	SortByUpdatedAt = "updated_at"
	SortByName      = "name"
	SortByEmail     = "email"
)

// Matching Engine Constants
const (
	// Default proficiency and experience values
	DefaultProficiencyLevel = 3
	DefaultYearsExperience  = 2

	// Scoring weights
	RequiredSkillsWeight  = 3.0
	PreferredSkillsWeight = 1.0

	// Base skill score multiplier
	BaseSkillScoreMultiplier = 2.0

	// Proficiency bonus multiplier
	ProficiencyBonusMultiplier = 0.5

	// Experience bonus multiplier
	ExperienceBonusMultiplier = 0.1

	// Coverage calculation multipliers
	CoverageBaseMultiplier  = 0.5
	CoverageBonusMultiplier = 0.5

	// Department match bonus
	DepartmentMatchBonus = 2.0

	// Experience level bonuses
	ExperienceLevelMatchBonus        = 1.5
	ExperienceLevelPartialMultiplier = 1.0

	// Location match bonus
	LocationMatchBonus = 1.0
)

// Experience level mappings
var ExperienceLevelMap = map[string]int{
	"junior":    1,
	"mid":       2,
	"senior":    3,
	"staff":     4,
	"principal": 5,
}

// Skill normalization mappings
var SkillNormalizationMap = map[string]string{
	"js":            "javascript",
	"javascript":    "javascript",
	"typescript":    "typescript",
	"ts":            "typescript",
	"react":         "react",
	"reactjs":       "react",
	"node":          "node.js",
	"nodejs":        "node.js",
	"node.js":       "node.js",
	"vue":           "vue.js",
	"vuejs":         "vue.js",
	"vue.js":        "vue.js",
	"angular":       "angular",
	"angularjs":     "angular",
	"java":          "java",
	"python":        "python",
	"py":            "python",
	"go":            "go",
	"golang":        "go",
	"c#":            "c#",
	"csharp":        "c#",
	"c sharp":       "c#",
	"f#":            "f#",
	"fsharp":        "f#",
	"f sharp":       "f#",
	"c++":           "c++",
	"cpp":           "c++",
	"c plus plus":   "c++",
	"vb.net":        "vb.net",
	"vbnet":         "vb.net",
	"vb .net":       "vb.net",
	"visual basic":  "visual basic",
	"vb":            "vb",
	"vb6":           "vb",
	"vb 6":          "vb",
	"docker":        "docker",
	"kubernetes":    "kubernetes",
	"k8s":           "kubernetes",
	"aws":           "aws",
	"postgresql":    "postgresql",
	"postgres":      "postgresql",
	"mongodb":       "mongodb",
	"mongo":         "mongodb",
	"redis":         "redis",
	"git":           "git",
	"graphql":       "graphql",
	"microservices": "microservices",
	"cicd":          "ci/cd",
	"ci/cd":         "ci/cd",
	"spring boot":   "spring boot",
	"spring":        "spring boot",
}

// Skill abbreviation mappings
var SkillAbbreviationMap = map[string][]string{
	"r":      {"react"},
	"js":     {"javascript"},
	"ts":     {"typescript"},
	"py":     {"python"},
	"go":     {"golang"},
	"c#":     {"csharp", "c sharp"},
	"f#":     {"fsharp", "f sharp"},
	"cpp":    {"c++", "c plus plus"},
	"vb":     {"vb.net", "vbnet", "visual basic"},
	"vb.net": {"vb", "vbnet", "visual basic"},
	"k8s":    {"kubernetes"},
	"aws":    {"amazon web services"},
	"db":     {"database"},
	"sql":    {"postgresql", "mysql", "database"},
}

// Query Management Constants
const (
	// File extensions
	SQLExtension = ".sql"

	// Error messages
	MsgQueryNotFound          = "query '%s' not found"
	MsgRequiredQueryNotFound  = "Required query '%s' not found: %v"
	MsgFailedToReadQueryDir   = "failed to read query directory"
	MsgFailedToReadQueryFile  = "failed to read query file %s"
	MsgMissingRequiredQueries = "missing required queries: %s"
)

// Required query names for validation
var RequiredQueries = []string{
	"get_all_employees",
	"get_employee_by_id",
	"create_employee",
	"update_employee",
	"delete_employee",
	"get_employee_skills",
	"get_all_skills",
	"get_skill_by_id",
	"get_skill_by_name",
	"create_skill",
	"update_skill",
	"delete_skill",
}
