package services

import (
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"time"
)

// CVExtractService defines the interface for CV extraction operations
type CVExtractService interface {
	CreateOrUpdateExtract(requestID string, status string, numFiles int, fileNumber int, metadata *string) (*models.CVExtract, error)
	GetExtractByRequestID(requestID string) (*models.CVExtract, error)
	UpdateExtractStatus(requestID string, status string) (*models.CVExtract, error)
	UpdateExtractProgress(requestID string, filesProcessed int, filesFailed int) (*models.CVExtract, error)
	MarkExtractSuccess(requestID string, totalTimeMs int64) (*models.CVExtract, error)
	MarkExtractFailed(requestID string, errorMessage string) (*models.CVExtract, error)
	UpdateFileProgress(requestID string, currentFile int, filesProcessed int, filesFailed int) (*models.CVExtract, error)
	CompleteExtract(requestID string, totalTimeMs int64, errorMessage *string) (*models.CVExtract, error)
	GetExtractStats(filters repositories.CVExtractStatsFilters) (*models.CVExtractStats, error)
	ListExtracts(filters repositories.CVExtractFilters) ([]models.CVExtract, int64, error)
	GetRecentExtracts(limit int) ([]models.CVExtract, error)
}

type cvExtractService struct {
	extractRepo repositories.CVExtractRepository
}

// NewCVExtractService creates a new CV extraction service
func NewCVExtractService(extractRepo repositories.CVExtractRepository) CVExtractService {
	return &cvExtractService{
		extractRepo: extractRepo,
	}
}

// CreateOrUpdateExtract creates a new extract record or updates existing one
func (s *cvExtractService) CreateOrUpdateExtract(requestID string, status string, numFiles int, fileNumber int, metadata *string) (*models.CVExtract, error) {
	// Try to get existing extract by request ID
	existingExtract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		// If not found, create new extract
		createData := models.CVExtractCreate{
			ExtractRequestID: requestID,
			Status:           status,
			NumFiles:         numFiles,
			FileNumber:       fileNumber,
			Metadata:         metadata,
		}

		extract, err := s.extractRepo.Create(createData)
		if err != nil {
			return nil, fmt.Errorf("failed to create CV extract: %w", err)
		}

		return extract, nil
	}

	// If found, update the status
	updateData := models.CVExtractUpdate{
		Status: &status,
	}

	extract, err := s.extractRepo.Update(existingExtract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to update CV extract: %w", err)
	}

	return extract, nil
}

// GetExtractByRequestID retrieves an extract by request ID
func (s *cvExtractService) GetExtractByRequestID(requestID string) (*models.CVExtract, error) {
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get CV extract: %w", err)
	}

	return extract, nil
}

// UpdateExtractStatus updates the status of an extract
func (s *cvExtractService) UpdateExtractStatus(requestID string, status string) (*models.CVExtract, error) {
	// Get existing extract
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("CV extract not found: %w", err)
	}

	// Update status
	updateData := models.CVExtractUpdate{
		Status: &status,
	}

	updatedExtract, err := s.extractRepo.Update(extract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to update CV extract status: %w", err)
	}

	return updatedExtract, nil
}

// UpdateExtractProgress updates the progress of an extract
func (s *cvExtractService) UpdateExtractProgress(requestID string, filesProcessed int, filesFailed int) (*models.CVExtract, error) {
	// Get existing extract
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("CV extract not found: %w", err)
	}

	// Update progress
	updateData := models.CVExtractUpdate{
		FilesProcessed: &filesProcessed,
		FilesFailed:    &filesFailed,
	}

	updatedExtract, err := s.extractRepo.Update(extract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to update CV extract progress: %w", err)
	}

	return updatedExtract, nil
}

// CompleteExtract marks an extract as completed with timing information
func (s *cvExtractService) CompleteExtract(requestID string, totalTimeMs int64, errorMessage *string) (*models.CVExtract, error) {
	// Get existing extract
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("CV extract not found: %w", err)
	}

	// Calculate average processing time
	var averageTimeMs *int64
	if extract.NumFiles > 0 {
		avg := totalTimeMs / int64(extract.NumFiles)
		averageTimeMs = &avg
	}

	// Determine final status
	status := models.CVExtractStatusCompleted
	if errorMessage != nil {
		status = models.CVExtractStatusFailed
	}

	now := time.Now()
	updateData := models.CVExtractUpdate{
		Status:                  &status,
		TotalProcessingTimeMs:   &totalTimeMs,
		AverageProcessingTimeMs: averageTimeMs,
		CompletedAt:             &now,
		ErrorMessage:            errorMessage,
	}

	updatedExtract, err := s.extractRepo.Update(extract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to complete CV extract: %w", err)
	}

	return updatedExtract, nil
}

// GetExtractStats retrieves CV extraction statistics
func (s *cvExtractService) GetExtractStats(filters repositories.CVExtractStatsFilters) (*models.CVExtractStats, error) {
	stats, err := s.extractRepo.GetStats(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get CV extract stats: %w", err)
	}

	return stats, nil
}

// ListExtracts retrieves a paginated list of extracts
func (s *cvExtractService) ListExtracts(filters repositories.CVExtractFilters) ([]models.CVExtract, int64, error) {
	extracts, total, err := s.extractRepo.List(filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list CV extracts: %w", err)
	}

	return extracts, total, nil
}

// GetRecentExtracts retrieves recent CV extraction records
func (s *cvExtractService) GetRecentExtracts(limit int) ([]models.CVExtract, error) {
	extracts, err := s.extractRepo.GetRecent(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent CV extracts: %w", err)
	}

	return extracts, nil
}

// MarkExtractSuccess marks an extract as successfully completed
func (s *cvExtractService) MarkExtractSuccess(requestID string, totalTimeMs int64) (*models.CVExtract, error) {
	// Get existing extract
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("CV extract not found: %w", err)
	}

	// Calculate average processing time
	var averageTimeMs *int64
	if extract.NumFiles > 0 {
		avg := totalTimeMs / int64(extract.NumFiles)
		averageTimeMs = &avg
	}

	now := time.Now()
	status := models.CVExtractStatusCompleted
	updateData := models.CVExtractUpdate{
		Status:                  &status,
		TotalProcessingTimeMs:   &totalTimeMs,
		AverageProcessingTimeMs: averageTimeMs,
		CompletedAt:             &now,
	}

	updatedExtract, err := s.extractRepo.Update(extract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to mark CV extract as successful: %w", err)
	}

	return updatedExtract, nil
}

// MarkExtractFailed marks an extract as failed
func (s *cvExtractService) MarkExtractFailed(requestID string, errorMessage string) (*models.CVExtract, error) {
	// Get existing extract
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("CV extract not found: %w", err)
	}

	now := time.Now()
	status := models.CVExtractStatusFailed
	updateData := models.CVExtractUpdate{
		Status:       &status,
		CompletedAt:  &now,
		ErrorMessage: &errorMessage,
	}

	updatedExtract, err := s.extractRepo.Update(extract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to mark CV extract as failed: %w", err)
	}

	return updatedExtract, nil
}

// UpdateFileProgress updates the progress of file processing
func (s *cvExtractService) UpdateFileProgress(requestID string, currentFile int, filesProcessed int, filesFailed int) (*models.CVExtract, error) {
	// Get existing extract
	extract, err := s.extractRepo.GetByRequestID(requestID)
	if err != nil {
		return nil, fmt.Errorf("CV extract not found: %w", err)
	}

	// Update progress
	updateData := models.CVExtractUpdate{
		FilesProcessed: &filesProcessed,
		FilesFailed:    &filesFailed,
	}

	updatedExtract, err := s.extractRepo.Update(extract.ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("failed to update CV extract progress: %w", err)
	}

	return updatedExtract, nil
}
