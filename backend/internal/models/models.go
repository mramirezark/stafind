package models

import (
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID             int       `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Email          string    `json:"email" db:"email"`
	Department     string    `json:"department" db:"department"`
	Level          string    `json:"level" db:"level"`
	Location       string    `json:"location" db:"location"`
	Bio            string    `json:"bio" db:"bio"`
	CurrentProject *string   `json:"current_project" db:"current_project"`
	Skills         []Skill   `json:"skills,omitempty"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Skill represents a technical skill
type Skill struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Category string `json:"category" db:"category"`
}

// EmployeeSkill represents the relationship between an employee and a skill
type EmployeeSkill struct {
	EmployeeID       int     `json:"employee_id" db:"employee_id"`
	SkillID          int     `json:"skill_id" db:"skill_id"`
	ProficiencyLevel int     `json:"proficiency_level" db:"proficiency_level"`
	YearsExperience  float64 `json:"years_experience" db:"years_experience"`
	Skill            Skill   `json:"skill,omitempty"`
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

// SkillExtractionRequest represents a request to extract skills from text
type SkillExtractionRequest struct {
	Text string `json:"text" binding:"required"`
}

// SkillExtractionResponse represents the response from skill extraction
type SkillExtractionResponse struct {
	Skills []string `json:"skills"`
	Text   string   `json:"text"`
}
