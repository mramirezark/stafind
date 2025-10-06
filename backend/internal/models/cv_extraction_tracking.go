package models

import (
	"time"
)

// CVExtractionTracking represents a CV extraction tracking record
type CVExtractionTracking struct {
	ID                      int        `json:"id" db:"id"`
	ExtractRequestId        string     `json:"extract_request_id" db:"extract_request_id"`
	Status                  string     `json:"status" db:"status"`
	NumFiles                int        `json:"num_files" db:"num_files"`
	FilesProcessed          int        `json:"files_processed" db:"files_processed"`
	FilesFailed             int        `json:"files_failed" db:"files_failed"`
	TotalProcessingTimeMs   *int64     `json:"total_processing_time_ms,omitempty" db:"total_processing_time_ms"`
	AverageProcessingTimeMs *int64     `json:"average_processing_time_ms,omitempty" db:"average_processing_time_ms"`
	StartedAt               time.Time  `json:"started_at" db:"started_at"`
	CompletedAt             *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	ErrorMessage            *string    `json:"error_message,omitempty" db:"error_message"`
	Metadata                *string    `json:"metadata,omitempty" db:"metadata"` // JSON string
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
}

// CVExtractionTrackingCreate represents data for creating a new CV extraction tracking record
type CVExtractionTrackingCreate struct {
	Status   string  `json:"status" validate:"required"`
	NumFiles int     `json:"num_files" validate:"required,min=0"`
	Metadata *string `json:"metadata,omitempty"`
}

// CVExtractionTrackingUpdate represents data for updating a CV extraction tracking record
type CVExtractionTrackingUpdate struct {
	Status                  *string    `json:"status,omitempty"`
	FilesProcessed          *int       `json:"files_processed,omitempty"`
	FilesFailed             *int       `json:"files_failed,omitempty"`
	TotalProcessingTimeMs   *int64     `json:"total_processing_time_ms,omitempty"`
	AverageProcessingTimeMs *int64     `json:"average_processing_time_ms,omitempty"`
	CompletedAt             *time.Time `json:"completed_at,omitempty"`
	ErrorMessage            *string    `json:"error_message,omitempty"`
	Metadata                *string    `json:"metadata,omitempty"`
}

// CVExtractionTrackingListResponse represents a paginated list of CV extraction tracking records
type CVExtractionTrackingListResponse struct {
	Records    []CVExtractionTracking `json:"records"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// CVExtractionStats represents statistics about CV extraction processes
type CVExtractionStats struct {
	TotalExtractions      int64   `json:"total_extractions"`
	SuccessfulExtractions int64   `json:"successful_extractions"`
	FailedExtractions     int64   `json:"failed_extractions"`
	AverageProcessingTime float64 `json:"average_processing_time_ms"`
	TotalFilesProcessed   int64   `json:"total_files_processed"`
	TotalFilesFailed      int64   `json:"total_files_failed"`
	SuccessRate           float64 `json:"success_rate"`
}

// Constants for CV extraction status
const (
	CVExtractionStatusPending    = "pending"
	CVExtractionStatusProcessing = "processing"
	CVExtractionStatusCompleted  = "completed"
	CVExtractionStatusFailed     = "failed"
)

// CVExtractionMetadata represents additional metadata for CV extraction tracking
type CVExtractionMetadata struct {
	Source           *string                `json:"source,omitempty"`          // e.g., "bulk_upload", "api", "scheduled"
	ProcessingMode   *string                `json:"processing_mode,omitempty"` // e.g., "parallel", "sequential"
	FileTypes        []string               `json:"file_types,omitempty"`      // e.g., ["pdf", "docx"]
	ExtractionConfig map[string]interface{} `json:"extraction_config,omitempty"`
	UserID           *int                   `json:"user_id,omitempty"`
	SessionID        *string                `json:"session_id,omitempty"`
	Notes            *string                `json:"notes,omitempty"`
}
