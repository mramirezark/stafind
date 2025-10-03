package models

import "time"

// GoogleDriveFile represents a file from Google Drive
type GoogleDriveFile struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	MimeType     string    `json:"mimeType"`
	Size         int64     `json:"size"`
	DownloadURL  string    `json:"downloadUrl,omitempty"`
	CreatedTime  time.Time `json:"createdTime"`
	ModifiedTime time.Time `json:"modifiedTime"`
	Parents      []string  `json:"parents"`
	WebViewLink  string    `json:"webViewLink"`
}

// GoogleDriveFolder represents a folder from Google Drive
type GoogleDriveFolder struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	MimeType     string              `json:"mimeType"`
	CreatedTime  time.Time           `json:"createdTime"`
	ModifiedTime time.Time           `json:"modifiedTime"`
	Parents      []string            `json:"parents"`
	Files        []GoogleDriveFile   `json:"files,omitempty"`
	Folders      []GoogleDriveFolder `json:"folders,omitempty"`
}

// GoogleDriveScanRequest represents a request to scan a Google Drive folder
type GoogleDriveScanRequest struct {
	FolderID        string   `json:"folder_id" validate:"required"`
	Recursive       bool     `json:"recursive"`
	FileTypes       []string `json:"file_types"` // e.g., ["application/pdf", "application/msword"]
	ExtractSkills   bool     `json:"extract_skills"`
	ProcessExisting bool     `json:"process_existing"` // Whether to process files already in database
}

// GoogleDriveScanResponse represents the response from scanning a Google Drive folder
type GoogleDriveScanResponse struct {
	FolderID        string                   `json:"folder_id"`
	FolderName      string                   `json:"folder_name"`
	TotalFiles      int                      `json:"total_files"`
	ProcessedFiles  int                      `json:"processed_files"`
	SkippedFiles    int                      `json:"skipped_files"`
	FailedFiles     int                      `json:"failed_files"`
	Files           []GoogleDriveFileProcess `json:"files"`
	ExtractedSkills []Skill                  `json:"extracted_skills,omitempty"`
	Errors          []string                 `json:"errors,omitempty"`
	ProcessingTime  time.Duration            `json:"processing_time"`
}

// GoogleDriveFileProcess represents a file that was processed during scanning
type GoogleDriveFileProcess struct {
	FileID          string    `json:"file_id"`
	FileName        string    `json:"file_name"`
	FileSize        int64     `json:"file_size"`
	Status          string    `json:"status"` // "processed", "skipped", "failed"
	UploadedFileID  *int      `json:"uploaded_file_id,omitempty"`
	ExtractedSkills []Skill   `json:"extracted_skills,omitempty"`
	Error           string    `json:"error,omitempty"`
	ProcessedAt     time.Time `json:"processed_at"`
}

// GoogleDriveConfig represents the configuration for Google Drive API
type GoogleDriveConfig struct {
	CredentialsPath string   `json:"credentials_path"`
	TokenPath       string   `json:"token_path"`
	Scopes          []string `json:"scopes"`
}

// GoogleDriveAuthToken represents the authentication token
type GoogleDriveAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	TokenType    string    `json:"token_type"`
}
