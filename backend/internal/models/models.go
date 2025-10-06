package models

import (
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID                  int                    `json:"id" db:"id"`
	Name                string                 `json:"name" db:"name"`
	Email               string                 `json:"email" db:"email"`
	Department          string                 `json:"department" db:"department"`
	Level               string                 `json:"level" db:"level"`
	Location            string                 `json:"location" db:"location"`
	Bio                 string                 `json:"bio" db:"bio"`
	CurrentProject      *string                `json:"current_project" db:"current_project"`
	ResumeUrl           *string                `json:"resume_url,omitempty" db:"resume_url"`
	OriginalText        *string                `json:"original_text,omitempty" db:"original_text"`
	ExtractedData       map[string]interface{} `json:"extracted_data,omitempty" db:"extracted_data"`
	ExtractionTimestamp *time.Time             `json:"extraction_timestamp,omitempty" db:"extraction_timestamp"`
	ExtractionSource    *string                `json:"extraction_source,omitempty" db:"extraction_source"`
	ExtractionStatus    *string                `json:"extraction_status,omitempty" db:"extraction_status"`
	Skills              []Skill                `json:"skills,omitempty"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

// Category represents a skill category
type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Skill represents a technical skill
type Skill struct {
	ID         int        `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Categories []Category `json:"categories,omitempty"`
}

// CreateSkillRequest represents a skill creation request
type CreateSkillRequest struct {
	Name       string `json:"name" db:"name"`
	Categories []int  `json:"categories,omitempty"`
}

// SkillWithCount represents a skill with employee count
type SkillWithCount struct {
	Skill
	EmployeeCount int `json:"employee_count" db:"employee_count"`
}

// SkillUpdate represents a skill update operation
type SkillUpdate struct {
	ID   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

// SkillStats represents statistics about skills
type SkillStats struct {
	TotalSkills      int             `json:"total_skills"`
	TotalCategories  int             `json:"total_categories"`
	MostPopularSkill *Skill          `json:"most_popular_skill,omitempty"`
	CategoryStats    []CategoryStats `json:"category_stats"`
}

// CategoryStats represents statistics for a skill category
type CategoryStats struct {
	Category      string `json:"category"`
	SkillCount    int    `json:"skill_count"`
	EmployeeCount int    `json:"employee_count"`
}

// SkillCategory represents the relationship between a skill and category
type SkillCategory struct {
	SkillID    int `json:"skill_id" db:"skill_id"`
	CategoryID int `json:"category_id" db:"category_id"`
}

// CategoryWithSkillCount represents a category with skill count
type CategoryWithSkillCount struct {
	Category
	SkillCount int `json:"skill_count" db:"skill_count"`
}

// EmployeeSkill represents the relationship between an employee and a skill
type EmployeeSkill struct {
	EmployeeID       int     `json:"employee_id" db:"employee_id"`
	SkillID          int     `json:"skill_id" db:"skill_id"`
	ProficiencyLevel int     `json:"proficiency_level" db:"proficiency_level"`
	YearsExperience  float64 `json:"years_experience" db:"years_experience"`
	Skill            Skill   `json:"skill,omitempty"`
}

// ProcessedResumeData represents data extracted from a resume by AI
type ProcessedResumeData struct {
	CandidateName       string           `json:"candidate_name"`
	ContactInfo         ContactInfo      `json:"contact_info"`
	SeniorityLevel      string           `json:"seniority_level"`
	YearsExperience     string           `json:"years_experience"`
	CurrentRole         string           `json:"current_role"`
	Skills              ResumeSkills     `json:"skills"`
	Experience          []WorkExperience `json:"experience"`
	Projects            []Project        `json:"projects"`
	Education           []Education      `json:"education"`
	Certifications      []string         `json:"certifications"`
	ProfessionalSummary string           `json:"professional_summary"`
	FileMetadata        FileMetadata     `json:"file_metadata"`
	ProcessingTimestamp string           `json:"processing_timestamp"`
}

// ContactInfo represents contact information
type ContactInfo struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
}

// ResumeSkills represents skills extracted from resume
type ResumeSkills struct {
	Technical  []string `json:"technical"`
	Soft       []string `json:"soft"`
	Languages  []string `json:"languages"`
	Tools      []string `json:"tools"`
	Frameworks []string `json:"frameworks"`
}

// WorkExperience represents work experience entry
type WorkExperience struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	Duration    string `json:"duration"`
	Years       string `json:"years"`
	Description string `json:"description"`
}

// Project represents a project entry
type Project struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	Role         string   `json:"role"`
	Duration     string   `json:"duration"`
}

// Education represents education entry
type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Year        string `json:"year"`
}

// FileMetadata represents file metadata from Google Drive
type FileMetadata struct {
	GoogleDriveID string `json:"google_drive_id"`
	Filename      string `json:"filename"`
	FileSize      int64  `json:"file_size"`
	CreatedTime   string `json:"created_time"`
	ModifiedTime  string `json:"modified_time"`
	WebViewLink   string `json:"web_view_link"`
}

// ResumeProcessingResult represents the result of processing resume data
type ResumeProcessingResult struct {
	CandidateName       string                      `json:"candidate_name"`
	EmployeeID          int                         `json:"employee_id"`
	Action              string                      `json:"action"` // "created" or "updated"
	Status              string                      `json:"status"`
	SkillsProcessed     *SkillProcessingResult      `json:"skills_processed"`
	ExperienceProcessed *ExperienceProcessingResult `json:"experience_processed"`
	ProjectsProcessed   *ProjectProcessingResult    `json:"projects_processed"`
	Warnings            []string                    `json:"warnings,omitempty"`
	ProcessedAt         time.Time                   `json:"processed_at"`
}

// SkillProcessingResult represents the result of processing skills
type SkillProcessingResult struct {
	TotalSkills   int `json:"total_skills"`
	NewSkills     int `json:"new_skills"`
	UpdatedSkills int `json:"updated_skills"`
}

// ExperienceProcessingResult represents the result of processing experience
type ExperienceProcessingResult struct {
	TotalEntries     int `json:"total_entries"`
	ProcessedEntries int `json:"processed_entries"`
}

// ProjectProcessingResult represents the result of processing projects
type ProjectProcessingResult struct {
	TotalProjects     int    `json:"total_projects"`
	ProcessedProjects int    `json:"processed_projects"`
	ProjectsSummary   string `json:"projects_summary,omitempty"`
}

// TeamsSearchResults represents search results from Teams workflow
type TeamsSearchResults struct {
	SearchSummary TeamsSearchSummary `json:"search_summary"`
	Candidates    []CandidateMatch   `json:"candidates"`
}

// TeamsSearchSummary represents the summary of a Teams search
type TeamsSearchSummary struct {
	OriginalRequest     string         `json:"original_request"`
	TotalFilesProcessed int            `json:"total_files_processed"`
	MatchingCandidates  int            `json:"matching_candidates"`
	SearchCriteria      SearchCriteria `json:"search_criteria"`
	ProcessedAt         string         `json:"processed_at"`
}

// SearchCriteria represents search criteria from Teams message
type SearchCriteria struct {
	Intent             string           `json:"intent"`
	SkillsRequired     []string         `json:"skills_required"`
	ExperienceLevel    string           `json:"experience_level"`
	YearsExperienceMin string           `json:"years_experience_min"`
	Language           string           `json:"language"`
	OriginalMessage    string           `json:"original_message"`
	SearchCriteria     DetailedCriteria `json:"search_criteria"`
	ResponseSuggestion string           `json:"response_suggestion"`
	TeamsContext       TeamsContext     `json:"teams_context"`
}

// DetailedCriteria represents detailed search criteria
type DetailedCriteria struct {
	PrimarySkill    string   `json:"primary_skill"`
	SecondarySkills []string `json:"secondary_skills"`
	ExperienceFocus string   `json:"experience_focus"`
	PriorityLevel   string   `json:"priority_level"`
}

// TeamsContext represents Teams message context
type TeamsContext struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	MessageID   string `json:"message_id"`
	Timestamp   string `json:"timestamp"`
}

// CandidateMatch represents a candidate match from search
type CandidateMatch struct {
	CandidateName   string           `json:"candidate_name"`
	ContactInfo     ContactInfo      `json:"contact_info"`
	SeniorityLevel  string           `json:"seniority_level"`
	YearsExperience string           `json:"years_experience"`
	CurrentRole     string           `json:"current_role"`
	Skills          CandidateSkills  `json:"skills"`
	Experience      []WorkExperience `json:"experience"`
	Projects        []Project        `json:"projects"`
	MatchAnalysis   MatchAnalysis    `json:"match_analysis"`
	FileMetadata    FileMetadata     `json:"file_metadata"`
	SearchContext   SearchContext    `json:"search_context"`
}

// CandidateSkills represents skills from candidate match
type CandidateSkills struct {
	Technical  []string `json:"technical"`
	Languages  []string `json:"languages"`
	Tools      []string `json:"tools"`
	Frameworks []string `json:"frameworks"`
}

// MatchAnalysis represents the analysis of candidate match
type MatchAnalysis struct {
	MatchesCriteria      bool     `json:"matches_criteria"`
	MatchScore           string   `json:"match_score"`
	PrimarySkillMatch    bool     `json:"primary_skill_match"`
	ExperienceLevelMatch bool     `json:"experience_level_match"`
	YearsExperienceMatch bool     `json:"years_experience_match"`
	MatchingSkills       []string `json:"matching_skills"`
	MissingSkills        []string `json:"missing_skills"`
	Strengths            []string `json:"strengths"`
	Recommendation       string   `json:"recommendation"`
}

// SearchContext represents the context of the search
type SearchContext struct {
	OriginalRequest string         `json:"original_request"`
	SearchCriteria  SearchCriteria `json:"search_criteria"`
	ProcessedAt     string         `json:"processed_at"`
}

// SearchProcessingResult represents the result of processing search results
type SearchProcessingResult struct {
	SearchID            int                         `json:"search_id"`
	OriginalRequest     string                      `json:"original_request"`
	ProcessedAt         time.Time                   `json:"processed_at"`
	Status              string                      `json:"status"`
	CandidatesFound     int                         `json:"candidates_found"`
	FilesProcessed      int                         `json:"files_processed"`
	CandidatesProcessed []CandidateProcessingResult `json:"candidates_processed"`
	Insights            SearchInsights              `json:"insights"`
}

// CandidateProcessingResult represents the result of processing a candidate
type CandidateProcessingResult struct {
	CandidateName   string                 `json:"candidate_name"`
	EmployeeID      int                    `json:"employee_id,omitempty"`
	Action          string                 `json:"action"` // "created", "found", "error"
	MatchScore      string                 `json:"match_score"`
	Status          string                 `json:"status"`
	SkillsProcessed *SkillProcessingResult `json:"skills_processed,omitempty"`
	Warnings        []string               `json:"warnings,omitempty"`
	Error           string                 `json:"error,omitempty"`
}

// SearchInsights represents insights from search results
type SearchInsights struct {
	TotalCandidates   int            `json:"total_candidates"`
	AverageMatchScore float64        `json:"average_match_score"`
	TopSkills         map[string]int `json:"top_skills"`
	ExperienceLevels  map[string]int `json:"experience_levels"`
	Recommendations   []string       `json:"recommendations"`
}

// ResourceSearch represents a resource search record
type ResourceSearch struct {
	ID                  int       `json:"id" db:"id"`
	SearchQuery         string    `json:"search_query" db:"search_query"`
	UserID              string    `json:"user_id" db:"user_id"`
	UserName            string    `json:"user_name" db:"user_name"`
	ChannelID           string    `json:"channel_id" db:"channel_id"`
	ChannelName         string    `json:"channel_name" db:"channel_name"`
	PrimarySkill        string    `json:"primary_skill" db:"primary_skill"`
	ExperienceLevel     string    `json:"experience_level" db:"experience_level"`
	YearsExperienceMin  int       `json:"years_experience_min" db:"years_experience_min"`
	TotalFilesProcessed int       `json:"total_files_processed" db:"total_files_processed"`
	CandidatesFound     int       `json:"candidates_found" db:"candidates_found"`
	SearchResults       string    `json:"search_results" db:"search_results"`
	Status              string    `json:"status" db:"status"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// Match represents a match between an employee and skills (for AI agent)
type Match struct {
	ID             int       `json:"id" db:"id"`
	EmployeeID     int       `json:"employee_id" db:"employee_id"`
	MatchScore     float64   `json:"match_score" db:"match_score"`
	MatchingSkills []string  `json:"matching_skills" db:"matching_skills"`
	Notes          string    `json:"notes" db:"notes"`
	Employee       Employee  `json:"employee,omitempty"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CreateEmployeeRequest represents the request to create a new employee
type CreateEmployeeRequest struct {
	Name           string             `json:"name" binding:"required"`
	Email          string             `json:"email" binding:"required,email"`
	Department     string             `json:"department"`
	Level          string             `json:"level"`
	Location       string             `json:"location"`
	Bio            string             `json:"bio"`
	CurrentProject string             `json:"current_project"`
	ResumeUrl      string             `json:"resume_url"`
	Skills         []EmployeeSkillReq `json:"skills"`
}

// EmployeeSkillReq represents a skill request for an employee
type EmployeeSkillReq struct {
	SkillName        string  `json:"skill_name" binding:"required"`
	ProficiencyLevel int     `json:"proficiency_level" binding:"min=1,max=5"`
	YearsExperience  float64 `json:"years_experience"`
}

// SearchRequest represents a request to search for employees
type SearchRequest struct {
	RequiredSkills  []string `json:"required_skills"`
	PreferredSkills []string `json:"preferred_skills"`
	Department      string   `json:"department"`
	ExperienceLevel string   `json:"experience_level"`
	Location        string   `json:"location"`
	MinMatchScore   float64  `json:"min_match_score"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalEmployees    int `json:"total_employees"`
	TotalRequests     int `json:"total_requests"`
	CompletedRequests int `json:"completed_requests"`
	PendingRequests   int `json:"pending_requests"`
}

// DepartmentStats represents department statistics
type DepartmentStats struct {
	Department string `json:"department"`
	Count      int    `json:"count"`
}

// SkillDemandStats represents skill demand statistics
type SkillDemandStats struct {
	SkillName string `json:"skill_name"`
	Count     int    `json:"count"`
	Category  string `json:"category"`
}

// TopSuggestedEmployee represents a top suggested employee
type TopSuggestedEmployee struct {
	EmployeeID     int     `json:"employee_id"`
	EmployeeName   string  `json:"employee_name"`
	EmployeeEmail  string  `json:"employee_email"`
	Department     string  `json:"department"`
	Level          string  `json:"level"`
	Location       string  `json:"location"`
	CurrentProject *string `json:"current_project"`
	MatchCount     int     `json:"match_count"`
	AvgMatchScore  float64 `json:"avg_match_score"`
}

// DashboardMetrics represents comprehensive dashboard metrics
type DashboardMetrics struct {
	Stats                 DashboardStats         `json:"stats"`
	MostRequestedSkills   []SkillDemandStats     `json:"most_requested_skills"`
	TopSuggestedEmployees []TopSuggestedEmployee `json:"top_suggested_employees"`
	RecentRequests        []AIAgentRequest       `json:"recent_requests"`
}

// AIAgentRequest represents a request from Microsoft Teams to the AI agent
type AIAgentRequest struct {
	ID              int        `json:"id" db:"id"`
	TeamsMessageID  string     `json:"teams_message_id" db:"teams_message_id"`
	ChannelID       string     `json:"channel_id" db:"channel_id"`
	UserID          string     `json:"user_id" db:"user_id"`
	UserName        string     `json:"user_name" db:"user_name"`
	MessageText     string     `json:"message_text" db:"message_text"`
	AttachmentURL   *string    `json:"attachment_url,omitempty" db:"attachment_url"`
	ExtractedText   *string    `json:"extracted_text,omitempty" db:"extracted_text"`
	ExtractedSkills []string   `json:"extracted_skills,omitempty" db:"extracted_skills"`
	Status          string     `json:"status" db:"status"` // pending, processing, completed, failed
	Error           *string    `json:"error,omitempty" db:"error"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// AIAgentResponse represents the response from the AI agent
type AIAgentResponse struct {
	RequestID      int            `json:"request_id"`
	Matches        []AIAgentMatch `json:"matches"`
	Summary        string         `json:"summary"`
	ProcessingTime int64          `json:"processing_time_ms"`
	Status         string         `json:"status"`
	Error          string         `json:"error,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

// AIAgentMatch represents a match result from the AI agent
type AIAgentMatch struct {
	EmployeeID     int      `json:"employee_id"`
	EmployeeName   string   `json:"employee_name"`
	EmployeeEmail  string   `json:"employee_email"`
	Position       string   `json:"position"`
	Seniority      string   `json:"seniority"`
	Location       string   `json:"location"`
	CurrentProject string   `json:"current_project"`
	ResumeLink     string   `json:"resume_link"`
	MatchScore     float64  `json:"match_score"`
	MatchingSkills []string `json:"matching_skills"`
	AISummary      string   `json:"ai_summary"`
	Bio            string   `json:"bio"`
}

// CreateAIAgentRequest represents the request to create a new AI agent request
type CreateAIAgentRequest struct {
	TeamsMessageID string  `json:"teams_message_id,omitempty"`
	ChannelID      string  `json:"channel_id,omitempty"`
	UserID         string  `json:"user_id,omitempty"`
	UserName       string  `json:"user_name,omitempty"`
	MessageText    string  `json:"message_text,omitempty"`
	AttachmentURL  *string `json:"attachment_url,omitempty"`
	// For public endpoint compatibility
	Description string   `json:"description,omitempty"`
	Skills      []string `json:"skills,omitempty"`
}

// SkillExtractRequest represents a request to extract skills from text
type SkillExtractRequest struct {
	Text string `json:"text" binding:"required"`
}

// SkillExtractResponse represents the response from skill extraction
type SkillExtractResponse struct {
	Skills []string `json:"skills"`
	Text   string   `json:"text"`
}

// ExtractAIRequest represents a request to process text with Llama AI
type ExtractAIRequest struct {
	TeamsMessageID string                 `json:"teams_message_id,omitempty"`
	ChannelID      string                 `json:"channel_id,omitempty"`
	UserID         string                 `json:"user_id,omitempty"`
	UserName       string                 `json:"user_name,omitempty"`
	MessageText    string                 `json:"message_text" binding:"required"`
	Text           string                 `json:"text" binding:"required"`
	FileName       string                 `json:"file_name" binding:"required"`
	FileURL        string                 `json:"file_url" binding:"required"`
	ProcessingType string                 `json:"processing_type,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ExtractAIResponse represents the response from Llama AI processing
type ExtractAIResponse struct {
	ProcessedContent string                 `json:"processed_content"`
	ProcessingTime   time.Duration          `json:"processing_time"`
	ModelUsed        string                 `json:"model_used"`
	TokensProcessed  int                    `json:"tokens_processed"`
	ProcessingType   string                 `json:"processing_type"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Timestamp        time.Time              `json:"timestamp"`
}

// CandidateInfo represents structured candidate information extracted from resume
type CandidateInfo struct {
	CandidateName   string          `json:"candidate_name"`
	ContactInfo     ContactInfo     `json:"contact_info"`
	SeniorityLevel  string          `json:"seniority_level"`
	YearsExperience string          `json:"years_experience"`
	CurrentPosition string          `json:"current_position"`
	Skills          CandidateSkills `json:"skills"`
	LastProject     LastProject     `json:"last_project"`
	Education       Education       `json:"education"`
	Certifications  []string        `json:"certifications"`
	Languages       []string        `json:"languages"`
	Summary         string          `json:"summary"`
}

// LastProject represents the last project information
type LastProject struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	Duration     string   `json:"duration"`
	Role         string   `json:"role"`
}

// LlamaSearchCriteria represents search criteria extracted from user request using Llama AI
type LlamaSearchCriteria struct {
	OriginalRequest    string                `json:"original_request"`
	Language           string                `json:"language"`
	SearchCriteria     SearchCriteriaDetails `json:"search_criteria"`
	PriorityLevel      string                `json:"priority_level"`
	Urgency            string                `json:"urgency"`
	ResponseSuggestion string                `json:"response_suggestion"`
}

// SearchCriteriaDetails represents detailed search criteria
type SearchCriteriaDetails struct {
	PrimarySkills      []string `json:"primary_skills"`
	SecondarySkills    []string `json:"secondary_skills"`
	Databases          []string `json:"databases"`
	Frameworks         []string `json:"frameworks"`
	ExperienceLevel    string   `json:"experience_level"`
	YearsExperienceMin string   `json:"years_experience_min"`
	PositionType       []string `json:"position_type"`
	ProjectFocus       []string `json:"project_focus"`
}

// CandidateMatchingRequest represents a request to match a candidate with search criteria
type CandidateMatchingRequest struct {
	CandidateName  string `json:"candidate_name"`
	CandidateInfo  string `json:"candidate_info" binding:"required"`
	SearchCriteria string `json:"search_criteria" binding:"required"`
	Language       string `json:"language,omitempty"`
}

// CandidateMatchResult represents the result of candidate matching
type CandidateMatchResult struct {
	MatchScore         int      `json:"match_score"`
	MatchPercentage    string   `json:"match_percentage"`
	MatchReasoning     string   `json:"match_reasoning"`
	Strengths          []string `json:"strengths"`
	Weaknesses         []string `json:"weaknesses"`
	MissingSkills      []string `json:"missing_skills"`
	Recommendation     string   `json:"recommendation"`
	InterviewQuestions []string `json:"interview_questions"`
	SalaryExpectation  string   `json:"salary_expectation"`
	Availability       string   `json:"availability"`
}

// ExtractProcessRequest represents a request for both AI agent and text extraction
type ExtractProcessRequest struct {
	ExtractRequestId string                 `json:"extract_request_id,omitempty"`
	Text             string                 `json:"text" binding:"required"` // The extracted text to process
	ResumeURL        string                 `json:"resume_url,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	FileNumber       int                    `json:"file_number,omitempty"`
	TotalFiles       int                    `json:"total_files,omitempty"`
	ExtractionSource string                 `json:"extraction_source,omitempty"` // Source of extraction (resume, linkedin, etc.)
	ProcessingType   string                 `json:"processing_type,omitempty"`
}

// CandidateExtractionResult represents the result of candidate extraction and storage
type CandidateExtractionResult struct {
	EmployeeID      int                    `json:"employee_id"`
	Action          string                 `json:"action"` // "created", "updated", "no_changes"
	Employee        *Employee              `json:"employee,omitempty"`
	ExtractedData   map[string]interface{} `json:"extracted_data"`
	ChangesDetected bool                   `json:"changes_detected"`
	ChangesSummary  []string               `json:"changes_summary,omitempty"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	Status          string                 `json:"status"`
	Message         string                 `json:"message"`
}

// MatchingResult represents the result of candidate matching
type MatchingResult struct {
	Requirements    string              `json:"requirements"`
	RequiredSkills  []string            `json:"required_skills"`
	Matches         []MatchingCandidate `json:"matches"`
	TotalCandidates int                 `json:"total_candidates"`
	ProcessingTime  time.Duration       `json:"processing_time"`
	Timestamp       time.Time           `json:"timestamp"`
}

// MatchingCandidate represents a single candidate match from the matching service
type MatchingCandidate struct {
	EmployeeID      int      `json:"employee_id"`
	EmployeeName    string   `json:"employee_name"`
	EmployeeEmail   string   `json:"employee_email"`
	EmployeeLevel   string   `json:"employee_level"`
	EmployeeSkills  []string `json:"employee_skills"`
	MatchedSkills   []string `json:"matched_skills"`
	MatchScore      int      `json:"match_score"`
	MatchLevel      string   `json:"match_level"`
	ExperienceMatch bool     `json:"experience_match"`
	SkillsMatch     int      `json:"skills_match"`
	TotalRequired   int      `json:"total_required"`
}
