package services

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"time"
)

// UploadedFileService defines the interface for uploaded file operations
type UploadedFileService interface {
	CreateFile(fileData models.UploadedFileCreate) (*models.UploadedFile, error)
	GetFile(id int) (*models.UploadedFile, error)
	GetFileByFilename(filename string) (*models.UploadedFile, error)
	UpdateFile(id int, fileData models.UploadedFileUpdate) (*models.UploadedFile, error)
	DeleteFile(id int) error
	SoftDeleteFile(id int) error
	ListFiles(filters repositories.FileListFilters) (*models.FileListResponse, error)
	GetFileStats() (*models.FileUploadStats, error)
	GetFilesByEmployee(employeeID int) ([]models.UploadedFile, error)
	GetFilesByType(uploadType string) ([]models.UploadedFile, error)
	CleanupOldFiles(daysOld int) (int64, error)
	CalculateFileHash(filePath string) (string, error)
	CheckDuplicateFile(hash string) (*models.UploadedFile, error)
}

type uploadedFileService struct {
	fileRepo repositories.UploadedFileRepository
}

// NewUploadedFileService creates a new uploaded file service
func NewUploadedFileService(fileRepo repositories.UploadedFileRepository) UploadedFileService {
	return &uploadedFileService{
		fileRepo: fileRepo,
	}
}

// CreateFile creates a new uploaded file record
func (s *uploadedFileService) CreateFile(fileData models.UploadedFileCreate) (*models.UploadedFile, error) {
	// Set default values
	if fileData.Status == "" {
		fileData.Status = models.FileStatusActive
	}
	if fileData.UploadType == "" {
		fileData.UploadType = models.FileTypeResume
	}

	// Calculate file hash if not provided
	if fileData.FileHash == nil {
		hash, err := s.CalculateFileHash(fileData.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate file hash: %w", err)
		}
		fileData.FileHash = &hash
	}

	// Check for duplicate files
	if fileData.FileHash != nil {
		existingFile, err := s.CheckDuplicateFile(*fileData.FileHash)
		if err == nil && existingFile != nil {
			return existingFile, nil // Return existing file if duplicate found
		}
	}

	// Create the file record
	file, err := s.fileRepo.Create(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

// GetFile retrieves an uploaded file by ID
func (s *uploadedFileService) GetFile(id int) (*models.UploadedFile, error) {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Check if file still exists on disk
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		// File doesn't exist on disk, mark as deleted
		file.Status = models.FileStatusDeleted
		now := time.Now()
		file.DeletedAt = &now
	}

	return file, nil
}

// GetFileByFilename retrieves an uploaded file by filename
func (s *uploadedFileService) GetFileByFilename(filename string) (*models.UploadedFile, error) {
	file, err := s.fileRepo.GetByFilename(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get file by filename: %w", err)
	}

	// Check if file still exists on disk
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		// File doesn't exist on disk, mark as deleted
		file.Status = models.FileStatusDeleted
		now := time.Now()
		file.DeletedAt = &now
	}

	return file, nil
}

// UpdateFile updates an uploaded file
func (s *uploadedFileService) UpdateFile(id int, fileData models.UploadedFileUpdate) (*models.UploadedFile, error) {
	file, err := s.fileRepo.Update(id, fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	return file, nil
}

// DeleteFile permanently deletes an uploaded file
func (s *uploadedFileService) DeleteFile(id int) error {
	// Get file info first
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get file for deletion: %w", err)
	}

	// Delete physical file
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete physical file: %w", err)
	}

	// Delete database record
	if err := s.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// SoftDeleteFile soft deletes an uploaded file
func (s *uploadedFileService) SoftDeleteFile(id int) error {
	// Get file info first
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get file for soft deletion: %w", err)
	}

	// Delete physical file
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete physical file: %w", err)
	}

	// Soft delete database record
	if err := s.fileRepo.SoftDelete(id); err != nil {
		return fmt.Errorf("failed to soft delete file record: %w", err)
	}

	return nil
}

// ListFiles retrieves a paginated list of uploaded files
func (s *uploadedFileService) ListFiles(filters repositories.FileListFilters) (*models.FileListResponse, error) {
	files, total, err := s.fileRepo.List(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// Calculate pagination info
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filters.Page
	if page <= 0 {
		page = 1
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := &models.FileListResponse{
		Files:      files,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// GetFileStats retrieves file upload statistics
func (s *uploadedFileService) GetFileStats() (*models.FileUploadStats, error) {
	stats, err := s.fileRepo.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	return stats, nil
}

// GetFilesByEmployee retrieves files associated with an employee
func (s *uploadedFileService) GetFilesByEmployee(employeeID int) ([]models.UploadedFile, error) {
	files, err := s.fileRepo.GetByEmployeeID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get files by employee: %w", err)
	}

	return files, nil
}

// GetFilesByType retrieves files by upload type
func (s *uploadedFileService) GetFilesByType(uploadType string) ([]models.UploadedFile, error) {
	files, err := s.fileRepo.GetByUploadType(uploadType)
	if err != nil {
		return nil, fmt.Errorf("failed to get files by type: %w", err)
	}

	return files, nil
}

// CleanupOldFiles removes old deleted files
func (s *uploadedFileService) CleanupOldFiles(daysOld int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	rowsAffected, err := s.fileRepo.CleanupOldFiles(cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old files: %w", err)
	}

	return rowsAffected, nil
}

// CalculateFileHash calculates SHA-256 hash of a file
func (s *uploadedFileService) CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CheckDuplicateFile checks if a file with the same hash already exists
func (s *uploadedFileService) CheckDuplicateFile(hash string) (*models.UploadedFile, error) {
	file, err := s.fileRepo.GetByHash(hash)
	if err != nil {
		return nil, err // File not found, not a duplicate
	}

	return file, nil
}
