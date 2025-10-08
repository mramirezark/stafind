package repositories

import (
	"database/sql"
	"fmt"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"
	"time"
)

// CVExtractRepository defines the interface for CV extraction operations
type CVExtractRepository interface {
	Create(tracking models.CVExtractCreate) (*models.CVExtract, error)
	GetByID(id int) (*models.CVExtract, error)
	GetByRequestID(requestID string) (*models.CVExtract, error)
	Upsert(tracking models.CVExtractCreate) (*models.CVExtract, error)
	Update(id int, tracking models.CVExtractUpdate) (*models.CVExtract, error)
	Delete(id int) error
	List(filters CVExtractFilters) ([]models.CVExtract, int64, error)
	GetStats(filters CVExtractStatsFilters) (*models.CVExtractStats, error)
	GetRecent(limit int) ([]models.CVExtract, error)
	GetByStatus(status string) ([]models.CVExtract, error)
}

// CVExtractFilters represents filters for listing CV extraction records
type CVExtractFilters struct {
	Status    *string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

// CVExtractStatsFilters represents filters for CV extraction statistics
type CVExtractStatsFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type cvExtractRepository struct {
	*BaseRepository
	db *sql.DB
}

// NewCVExtractRepository creates a new CV extraction repository
func NewCVExtractRepository(db *sql.DB) (CVExtractRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &cvExtractRepository{
		BaseRepository: baseRepo,
		db:             db,
	}, nil
}

// Create creates a new CV extraction record
func (r *cvExtractRepository) Create(tracking models.CVExtractCreate) (*models.CVExtract, error) {
	query := r.MustGetQuery("create_cv_extract")

	var cvExtract models.CVExtract
	err := r.db.QueryRow(
		query,
		tracking.ExtractRequestID,
		tracking.Status,
		tracking.NumFiles,
		0,          // files_processed
		0,          // files_failed
		nil,        // total_processing_time_ms
		nil,        // average_processing_time_ms
		time.Now(), // started_at
		nil,        // completed_at
		nil,        // error_message
		tracking.Metadata,
	).Scan(
		&cvExtract.ID,
		&cvExtract.CreatedAt,
		&cvExtract.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create CV extraction: %w", err)
	}

	// Set other fields
	cvExtract.ExtractRequestID = tracking.ExtractRequestID
	cvExtract.Status = tracking.Status
	cvExtract.NumFiles = tracking.NumFiles
	cvExtract.FilesProcessed = tracking.FileNumber
	cvExtract.FilesProcessed = 0
	cvExtract.FilesFailed = 0
	cvExtract.StartedAt = time.Now()
	cvExtract.Metadata = tracking.Metadata

	return &cvExtract, nil
}

// GetByID retrieves a CV extraction record by ID
func (r *cvExtractRepository) GetByID(id int) (*models.CVExtract, error) {
	query := r.MustGetQuery("get_cv_extract_by_id")

	var tracking models.CVExtract
	err := r.db.QueryRow(query, id).Scan(
		&tracking.ID,
		&tracking.ExtractRequestID,
		&tracking.Status,
		&tracking.NumFiles,
		&tracking.FilesProcessed,
		&tracking.FilesFailed,
		&tracking.TotalProcessingTimeMs,
		&tracking.AverageProcessingTimeMs,
		&tracking.StartedAt,
		&tracking.CompletedAt,
		&tracking.ErrorMessage,
		&tracking.Metadata,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("CV extraction not found")
		}
		return nil, fmt.Errorf("failed to get CV extraction: %w", err)
	}

	return &tracking, nil
}

// GetByRequestID retrieves a CV extraction record by extract request ID
func (r *cvExtractRepository) GetByRequestID(requestID string) (*models.CVExtract, error) {
	query := r.MustGetQuery("get_cv_extract_by_request_id")

	var tracking models.CVExtract
	err := r.db.QueryRow(query, requestID).Scan(
		&tracking.ID,
		&tracking.ExtractRequestID,
		&tracking.Status,
		&tracking.NumFiles,
		&tracking.FilesProcessed,
		&tracking.FilesFailed,
		&tracking.TotalProcessingTimeMs,
		&tracking.AverageProcessingTimeMs,
		&tracking.StartedAt,
		&tracking.CompletedAt,
		&tracking.ErrorMessage,
		&tracking.Metadata,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("CV extraction not found")
		}
		return nil, fmt.Errorf("failed to get CV extraction: %w", err)
	}

	return &tracking, nil
}

// Upsert creates or updates a CV extraction record
func (r *cvExtractRepository) Upsert(tracking models.CVExtractCreate) (*models.CVExtract, error) {
	query := r.MustGetQuery("upsert_cv_extract")

	var cvExtract models.CVExtract
	err := r.db.QueryRow(
		query,
		tracking.ExtractRequestID,
		tracking.Status,
		tracking.NumFiles,
		0,          // files_processed
		0,          // files_failed
		nil,        // total_processing_time_ms
		nil,        // average_processing_time_ms
		time.Now(), // started_at
		nil,        // completed_at
		nil,        // error_message
		tracking.Metadata,
	).Scan(
		&cvExtract.ID,
		&cvExtract.ExtractRequestID,
		&cvExtract.Status,
		&cvExtract.NumFiles,
		&cvExtract.FilesProcessed,
		&cvExtract.FilesFailed,
		&cvExtract.TotalProcessingTimeMs,
		&cvExtract.AverageProcessingTimeMs,
		&cvExtract.StartedAt,
		&cvExtract.CompletedAt,
		&cvExtract.ErrorMessage,
		&cvExtract.Metadata,
		&cvExtract.CreatedAt,
		&cvExtract.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert CV extraction: %w", err)
	}

	return &cvExtract, nil
}

// Update updates a CV extraction record
func (r *cvExtractRepository) Update(id int, tracking models.CVExtractUpdate) (*models.CVExtract, error) {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{id}
	argIndex := 2

	if tracking.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *tracking.Status)
		argIndex++
	}

	if tracking.FilesProcessed != nil {
		setParts = append(setParts, fmt.Sprintf("files_processed = $%d", argIndex))
		args = append(args, *tracking.FilesProcessed)
		argIndex++
	}

	if tracking.FilesFailed != nil {
		setParts = append(setParts, fmt.Sprintf("files_failed = $%d", argIndex))
		args = append(args, *tracking.FilesFailed)
		argIndex++
	}

	if tracking.TotalProcessingTimeMs != nil {
		setParts = append(setParts, fmt.Sprintf("total_processing_time_ms = $%d", argIndex))
		args = append(args, *tracking.TotalProcessingTimeMs)
		argIndex++
	}

	if tracking.AverageProcessingTimeMs != nil {
		setParts = append(setParts, fmt.Sprintf("average_processing_time_ms = $%d", argIndex))
		args = append(args, *tracking.AverageProcessingTimeMs)
		argIndex++
	}

	if tracking.CompletedAt != nil {
		setParts = append(setParts, fmt.Sprintf("completed_at = $%d", argIndex))
		args = append(args, *tracking.CompletedAt)
		argIndex++
	}

	if tracking.ErrorMessage != nil {
		setParts = append(setParts, fmt.Sprintf("error_message = $%d", argIndex))
		args = append(args, *tracking.ErrorMessage)
		argIndex++
	}

	if tracking.Metadata != nil {
		setParts = append(setParts, fmt.Sprintf("metadata = $%d", argIndex))
		args = append(args, *tracking.Metadata)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(id)
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())

	// Build the query - join all setParts with commas
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query := fmt.Sprintf(`
		UPDATE cv_extract 
		SET %s
		WHERE id = $1
		RETURNING id, extract_request_id, status, num_files, files_processed, files_failed,
		          total_processing_time_ms, average_processing_time_ms,
		          started_at, completed_at, error_message, metadata,
		          created_at, updated_at`,
		setClause,
	)

	var cvExtract models.CVExtract
	err := r.db.QueryRow(query, args...).Scan(
		&cvExtract.ID,
		&cvExtract.ExtractRequestID,
		&cvExtract.Status,
		&cvExtract.NumFiles,
		&cvExtract.FilesProcessed,
		&cvExtract.FilesFailed,
		&cvExtract.TotalProcessingTimeMs,
		&cvExtract.AverageProcessingTimeMs,
		&cvExtract.StartedAt,
		&cvExtract.CompletedAt,
		&cvExtract.ErrorMessage,
		&cvExtract.Metadata,
		&cvExtract.CreatedAt,
		&cvExtract.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("CV extraction not found")
		}
		return nil, fmt.Errorf("failed to update CV extraction: %w", err)
	}

	return &cvExtract, nil
}

// Delete permanently deletes a CV extraction record
func (r *cvExtractRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_cv_extract")
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete CV extraction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("CV extraction not found")
	}

	return nil
}

// List retrieves a paginated list of CV extraction records
func (r *cvExtractRepository) List(filters CVExtractFilters) ([]models.CVExtract, int64, error) {
	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filters.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("started_at >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("started_at <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + whereConditions[0]
		for i := 1; i < len(whereConditions); i++ {
			whereClause += " AND " + whereConditions[i]
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM cv_extract %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count CV extraction records: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		orderBy = filters.SortBy
		if filters.SortOrder == "asc" {
			orderBy += " ASC"
		} else {
			orderBy += " DESC"
		}
	}

	// Build LIMIT and OFFSET
	limit := constants.DefaultPageSize
	if filters.PageSize > 0 {
		limit = filters.PageSize
	}
	offset := 0
	if filters.Page > 0 {
		offset = (filters.Page - 1) * limit
	}

	// Get records
	query := fmt.Sprintf(`
		SELECT id, status, num_files, files_processed, files_failed,
		       total_processing_time_ms, average_processing_time_ms,
		       started_at, completed_at, error_message, metadata,
		       created_at, updated_at
		FROM cv_extract
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list CV extraction records: %w", err)
	}
	defer rows.Close()

	var records []models.CVExtract
	for rows.Next() {
		var tracking models.CVExtract
		err := rows.Scan(
			&tracking.ID,
			&tracking.Status,
			&tracking.NumFiles,
			&tracking.FilesProcessed,
			&tracking.FilesFailed,
			&tracking.TotalProcessingTimeMs,
			&tracking.AverageProcessingTimeMs,
			&tracking.StartedAt,
			&tracking.CompletedAt,
			&tracking.ErrorMessage,
			&tracking.Metadata,
			&tracking.CreatedAt,
			&tracking.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan CV extraction record: %w", err)
		}
		records = append(records, tracking)
	}

	return records, total, nil
}

// GetStats retrieves CV extraction statistics
func (r *cvExtractRepository) GetStats(filters CVExtractStatsFilters) (*models.CVExtractStats, error) {
	query := r.MustGetQuery("get_cv_extract_stats")

	var stats models.CVExtractStats
	err := r.db.QueryRow(query, filters.StartDate, filters.EndDate).Scan(
		&stats.TotalExtractions,
		&stats.SuccessfulExtractions,
		&stats.FailedExtractions,
		&stats.AverageProcessingTime,
		&stats.TotalFilesProcessed,
		&stats.TotalFilesFailed,
		&stats.SuccessRate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get CV extraction stats: %w", err)
	}

	return &stats, nil
}

// GetRecent retrieves recent CV extraction records
func (r *cvExtractRepository) GetRecent(limit int) ([]models.CVExtract, error) {
	query := r.MustGetQuery("get_recent_cv_extract")

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent CV extraction records: %w", err)
	}
	defer rows.Close()

	var records []models.CVExtract
	for rows.Next() {
		var tracking models.CVExtract
		err := rows.Scan(
			&tracking.ID,
			&tracking.Status,
			&tracking.NumFiles,
			&tracking.FilesProcessed,
			&tracking.FilesFailed,
			&tracking.TotalProcessingTimeMs,
			&tracking.AverageProcessingTimeMs,
			&tracking.StartedAt,
			&tracking.CompletedAt,
			&tracking.ErrorMessage,
			&tracking.Metadata,
			&tracking.CreatedAt,
			&tracking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent CV extraction record: %w", err)
		}
		records = append(records, tracking)
	}

	return records, nil
}

// GetByStatus retrieves CV extraction records by status
func (r *cvExtractRepository) GetByStatus(status string) ([]models.CVExtract, error) {
	query := r.MustGetQuery("get_cv_extract_by_status")

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get CV extraction records by status: %w", err)
	}
	defer rows.Close()

	var records []models.CVExtract
	for rows.Next() {
		var tracking models.CVExtract
		err := rows.Scan(
			&tracking.ID,
			&tracking.Status,
			&tracking.NumFiles,
			&tracking.FilesProcessed,
			&tracking.FilesFailed,
			&tracking.TotalProcessingTimeMs,
			&tracking.AverageProcessingTimeMs,
			&tracking.StartedAt,
			&tracking.CompletedAt,
			&tracking.ErrorMessage,
			&tracking.Metadata,
			&tracking.CreatedAt,
			&tracking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan CV extraction record by status: %w", err)
		}
		records = append(records, tracking)
	}

	return records, nil
}
