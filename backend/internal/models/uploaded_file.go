package models

import (
	"time"
)

// UploadedFile represents an uploaded file in the system
type UploadedFile struct {
	ID               int        `json:"id" db:"id"`
	Filename         string     `json:"filename" db:"filename"`
	OriginalFilename string     `json:"original_filename" db:"original_filename"`
	FilePath         string     `json:"file_path" db:"file_path"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	ContentType      string     `json:"content_type" db:"content_type"`
	FileHash         *string    `json:"file_hash,omitempty" db:"file_hash"`
	UploadType       string     `json:"upload_type" db:"upload_type"`
	Status           string     `json:"status" db:"status"`
	UploadedBy       *int       `json:"uploaded_by,omitempty" db:"uploaded_by"`
	EmployeeID       *int       `json:"employee_id,omitempty" db:"employee_id"`
	Metadata         *string    `json:"metadata,omitempty" db:"metadata"` // JSON string
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// UploadedFileCreate represents data for creating a new uploaded file
type UploadedFileCreate struct {
	Filename         string  `json:"filename" validate:"required"`
	OriginalFilename string  `json:"original_filename" validate:"required"`
	FilePath         string  `json:"file_path" validate:"required"`
	FileSize         int64   `json:"file_size" validate:"required,min=1"`
	ContentType      string  `json:"content_type" validate:"required"`
	FileHash         *string `json:"file_hash,omitempty"`
	UploadType       string  `json:"upload_type" validate:"required"`
	Status           string  `json:"status"`
	UploadedBy       *int    `json:"uploaded_by,omitempty"`
	EmployeeID       *int    `json:"employee_id,omitempty"`
	Metadata         *string `json:"metadata,omitempty"`
}

// UploadedFileUpdate represents data for updating an uploaded file
type UploadedFileUpdate struct {
	Status     *string `json:"status,omitempty"`
	EmployeeID *int    `json:"employee_id,omitempty"`
	Metadata   *string `json:"metadata,omitempty"`
}

// FileMetadata represents additional metadata for uploaded files
type FileMetadata struct {
	ParsedData       *ParsedResumeData `json:"parsed_data,omitempty"`
	ProcessingStatus string            `json:"processing_status,omitempty"`
	Error            *string           `json:"error,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Description      *string           `json:"description,omitempty"`
}

// ParsedResumeData represents parsed data from resume files
type ParsedResumeData struct {
	Name        string           `json:"name"`
	Email       string           `json:"email"`
	Phone       *string          `json:"phone,omitempty"`
	Skills      []string         `json:"skills"`
	Experience  string           `json:"experience"`
	Location    *string          `json:"location,omitempty"`
	Bio         *string          `json:"bio,omitempty"`
	Education   []string         `json:"education,omitempty"`
	WorkHistory []WorkExperience `json:"work_history,omitempty"`
}

// WorkExperience represents work experience from parsed resume
type WorkExperience struct {
	Company     string  `json:"company"`
	Position    string  `json:"position"`
	Duration    string  `json:"duration"`
	Description *string `json:"description,omitempty"`
}

// FileUploadStats represents statistics about file uploads
type FileUploadStats struct {
	TotalFiles    int64            `json:"total_files"`
	TotalSize     int64            `json:"total_size"`
	FilesByType   map[string]int64 `json:"files_by_type"`
	FilesByStatus map[string]int64 `json:"files_by_status"`
	RecentUploads []UploadedFile   `json:"recent_uploads"`
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	UploadType string  `json:"upload_type" validate:"required"`
	EmployeeID *int    `json:"employee_id,omitempty"`
	Metadata   *string `json:"metadata,omitempty"`
}

// FileUploadResponse represents the response after file upload
type FileUploadResponse struct {
	File    UploadedFile `json:"file"`
	URL     string       `json:"url"`
	Message string       `json:"message"`
}

// FileListResponse represents a paginated list of files
type FileListResponse struct {
	Files      []UploadedFile `json:"files"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// Constants for file status and types
const (
	FileStatusActive   = "active"
	FileStatusDeleted  = "deleted"
	FileStatusArchived = "archived"

	FileTypeResume   = "resume"
	FileTypeDocument = "document"
	FileTypeImage    = "image"
	FileTypeOther    = "other"

	ProcessingStatusPending    = "pending"
	ProcessingStatusProcessing = "processing"
	ProcessingStatusCompleted  = "completed"
	ProcessingStatusFailed     = "failed"
)
