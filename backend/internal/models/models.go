package models

import (
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
	Department string    `json:"department" db:"department"`
	Level      string    `json:"level" db:"level"`
	Location   string    `json:"location" db:"location"`
	Bio        string    `json:"bio" db:"bio"`
	Skills     []Skill   `json:"skills,omitempty"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
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

// JobRequest represents a job description or tech stack request
type JobRequest struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	Department      string    `json:"department" db:"department"`
	RequiredSkills  []string  `json:"required_skills" db:"required_skills"`
	PreferredSkills []string  `json:"preferred_skills" db:"preferred_skills"`
	ExperienceLevel string    `json:"experience_level" db:"experience_level"`
	Location        string    `json:"location" db:"location"`
	Priority        string    `json:"priority" db:"priority"`
	Status          string    `json:"status" db:"status"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Match represents a match between a job request and an employee
type Match struct {
	ID             int       `json:"id" db:"id"`
	JobRequestID   int       `json:"job_request_id" db:"job_request_id"`
	EmployeeID     int       `json:"employee_id" db:"employee_id"`
	MatchScore     float64   `json:"match_score" db:"match_score"`
	MatchingSkills []string  `json:"matching_skills" db:"matching_skills"`
	Notes          string    `json:"notes" db:"notes"`
	Employee       Employee  `json:"employee,omitempty"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CreateEmployeeRequest represents the request to create a new employee
type CreateEmployeeRequest struct {
	Name       string             `json:"name" binding:"required"`
	Email      string             `json:"email" binding:"required,email"`
	Department string             `json:"department"`
	Level      string             `json:"level"`
	Location   string             `json:"location"`
	Bio        string             `json:"bio"`
	Skills     []EmployeeSkillReq `json:"skills"`
}

// EmployeeSkillReq represents a skill request for an employee
type EmployeeSkillReq struct {
	SkillName        string  `json:"skill_name" binding:"required"`
	ProficiencyLevel int     `json:"proficiency_level" binding:"min=1,max=5"`
	YearsExperience  float64 `json:"years_experience"`
}

// CreateJobRequestRequest represents the request to create a new job request
type CreateJobRequestRequest struct {
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description"`
	Department      string   `json:"department"`
	RequiredSkills  []string `json:"required_skills"`
	PreferredSkills []string `json:"preferred_skills"`
	ExperienceLevel string   `json:"experience_level"`
	Location        string   `json:"location"`
	Priority        string   `json:"priority"`
	CreatedBy       string   `json:"created_by"`
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
	TotalEmployees   int `json:"total_employees"`
	TotalJobRequests int `json:"total_job_requests"`
	ActiveMatches    int `json:"active_matches"`
	RecentRequests   int `json:"recent_requests"`
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
