package repositories

import (
	"database/sql"
	"fmt"
	"stafind-backend/internal/models"
	"time"
)

// UploadedFileRepository defines the interface for uploaded file operations
type UploadedFileRepository interface {
	Create(file models.UploadedFileCreate) (*models.UploadedFile, error)
	GetByID(id int) (*models.UploadedFile, error)
	GetByFilename(filename string) (*models.UploadedFile, error)
	GetByHash(hash string) (*models.UploadedFile, error)
	Update(id int, file models.UploadedFileUpdate) (*models.UploadedFile, error)
	Delete(id int) error
	SoftDelete(id int) error
	List(filters FileListFilters) ([]models.UploadedFile, int64, error)
	GetStats() (*models.FileUploadStats, error)
	GetByEmployeeID(employeeID int) ([]models.UploadedFile, error)
	GetByUploadType(uploadType string) ([]models.UploadedFile, error)
	CleanupOldFiles(olderThan time.Time) (int64, error)
}

// FileListFilters represents filters for listing files
type FileListFilters struct {
	UploadType *string
	Status     *string
	UploadedBy *int
	EmployeeID *int
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
	Search     *string
}

type uploadedFileRepository struct {
	db *sql.DB
}

// NewUploadedFileRepository creates a new uploaded file repository
func NewUploadedFileRepository(db *sql.DB) UploadedFileRepository {
	return &uploadedFileRepository{db: db}
}

// Create creates a new uploaded file record
func (r *uploadedFileRepository) Create(file models.UploadedFileCreate) (*models.UploadedFile, error) {
	query := `
		INSERT INTO uploaded_files (
			filename, original_filename, file_path, file_size, content_type,
			file_hash, upload_type, status, uploaded_by, employee_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	var uploadedFile models.UploadedFile
	err := r.db.QueryRow(
		query,
		file.Filename,
		file.OriginalFilename,
		file.FilePath,
		file.FileSize,
		file.ContentType,
		file.FileHash,
		file.UploadType,
		file.Status,
		file.UploadedBy,
		file.EmployeeID,
		file.Metadata,
	).Scan(
		&uploadedFile.ID,
		&uploadedFile.CreatedAt,
		&uploadedFile.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create uploaded file: %w", err)
	}

	// Set other fields
	uploadedFile.Filename = file.Filename
	uploadedFile.OriginalFilename = file.OriginalFilename
	uploadedFile.FilePath = file.FilePath
	uploadedFile.FileSize = file.FileSize
	uploadedFile.ContentType = file.ContentType
	uploadedFile.FileHash = file.FileHash
	uploadedFile.UploadType = file.UploadType
	uploadedFile.Status = file.Status
	uploadedFile.UploadedBy = file.UploadedBy
	uploadedFile.EmployeeID = file.EmployeeID
	uploadedFile.Metadata = file.Metadata

	return &uploadedFile, nil
}

// GetByID retrieves an uploaded file by ID
func (r *uploadedFileRepository) GetByID(id int) (*models.UploadedFile, error) {
	query := `
		SELECT id, filename, original_filename, file_path, file_size, content_type,
		       file_hash, upload_type, status, uploaded_by, employee_id, metadata,
		       created_at, updated_at, deleted_at
		FROM uploaded_files
		WHERE id = $1 AND deleted_at IS NULL`

	var file models.UploadedFile
	err := r.db.QueryRow(query, id).Scan(
		&file.ID,
		&file.Filename,
		&file.OriginalFilename,
		&file.FilePath,
		&file.FileSize,
		&file.ContentType,
		&file.FileHash,
		&file.UploadType,
		&file.Status,
		&file.UploadedBy,
		&file.EmployeeID,
		&file.Metadata,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("uploaded file not found")
		}
		return nil, fmt.Errorf("failed to get uploaded file: %w", err)
	}

	return &file, nil
}

// GetByFilename retrieves an uploaded file by filename
func (r *uploadedFileRepository) GetByFilename(filename string) (*models.UploadedFile, error) {
	query := `
		SELECT id, filename, original_filename, file_path, file_size, content_type,
		       file_hash, upload_type, status, uploaded_by, employee_id, metadata,
		       created_at, updated_at, deleted_at
		FROM uploaded_files
		WHERE filename = $1 AND deleted_at IS NULL`

	var file models.UploadedFile
	err := r.db.QueryRow(query, filename).Scan(
		&file.ID,
		&file.Filename,
		&file.OriginalFilename,
		&file.FilePath,
		&file.FileSize,
		&file.ContentType,
		&file.FileHash,
		&file.UploadType,
		&file.Status,
		&file.UploadedBy,
		&file.EmployeeID,
		&file.Metadata,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("uploaded file not found")
		}
		return nil, fmt.Errorf("failed to get uploaded file: %w", err)
	}

	return &file, nil
}

// GetByHash retrieves an uploaded file by file hash
func (r *uploadedFileRepository) GetByHash(hash string) (*models.UploadedFile, error) {
	query := `
		SELECT id, filename, original_filename, file_path, file_size, content_type,
		       file_hash, upload_type, status, uploaded_by, employee_id, metadata,
		       created_at, updated_at, deleted_at
		FROM uploaded_files
		WHERE file_hash = $1 AND deleted_at IS NULL`

	var file models.UploadedFile
	err := r.db.QueryRow(query, hash).Scan(
		&file.ID,
		&file.Filename,
		&file.OriginalFilename,
		&file.FilePath,
		&file.FileSize,
		&file.ContentType,
		&file.FileHash,
		&file.UploadType,
		&file.Status,
		&file.UploadedBy,
		&file.EmployeeID,
		&file.Metadata,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("uploaded file not found")
		}
		return nil, fmt.Errorf("failed to get uploaded file: %w", err)
	}

	return &file, nil
}

// Update updates an uploaded file
func (r *uploadedFileRepository) Update(id int, file models.UploadedFileUpdate) (*models.UploadedFile, error) {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if file.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *file.Status)
		argIndex++
	}

	if file.EmployeeID != nil {
		setParts = append(setParts, fmt.Sprintf("employee_id = $%d", argIndex))
		args = append(args, *file.EmployeeID)
		argIndex++
	}

	if file.Metadata != nil {
		setParts = append(setParts, fmt.Sprintf("metadata = $%d", argIndex))
		args = append(args, *file.Metadata)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(id)
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE uploaded_files 
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, filename, original_filename, file_path, file_size, content_type,
		          file_hash, upload_type, status, uploaded_by, employee_id, metadata,
		          created_at, updated_at, deleted_at`,
		fmt.Sprintf("%s", setParts[0]),
		argIndex,
	)

	// Fix the query building
	query = fmt.Sprintf(`
		UPDATE uploaded_files 
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, filename, original_filename, file_path, file_size, content_type,
		          file_hash, upload_type, status, uploaded_by, employee_id, metadata,
		          created_at, updated_at, deleted_at`,
		fmt.Sprintf("%s", setParts[0]),
		argIndex,
	)

	var uploadedFile models.UploadedFile
	err := r.db.QueryRow(query, args...).Scan(
		&uploadedFile.ID,
		&uploadedFile.Filename,
		&uploadedFile.OriginalFilename,
		&uploadedFile.FilePath,
		&uploadedFile.FileSize,
		&uploadedFile.ContentType,
		&uploadedFile.FileHash,
		&uploadedFile.UploadType,
		&uploadedFile.Status,
		&uploadedFile.UploadedBy,
		&uploadedFile.EmployeeID,
		&uploadedFile.Metadata,
		&uploadedFile.CreatedAt,
		&uploadedFile.UpdatedAt,
		&uploadedFile.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("uploaded file not found")
		}
		return nil, fmt.Errorf("failed to update uploaded file: %w", err)
	}

	return &uploadedFile, nil
}

// Delete permanently deletes an uploaded file record
func (r *uploadedFileRepository) Delete(id int) error {
	query := `DELETE FROM uploaded_files WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete uploaded file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("uploaded file not found")
	}

	return nil
}

// SoftDelete soft deletes an uploaded file
func (r *uploadedFileRepository) SoftDelete(id int) error {
	query := `
		UPDATE uploaded_files 
		SET status = 'deleted', deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete uploaded file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("uploaded file not found")
	}

	return nil
}

// List retrieves a paginated list of uploaded files
func (r *uploadedFileRepository) List(filters FileListFilters) ([]models.UploadedFile, int64, error) {
	// Build WHERE clause
	whereConditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	if filters.UploadType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("upload_type = $%d", argIndex))
		args = append(args, *filters.UploadType)
		argIndex++
	}

	if filters.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.UploadedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("uploaded_by = $%d", argIndex))
		args = append(args, *filters.UploadedBy)
		argIndex++
	}

	if filters.EmployeeID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("employee_id = $%d", argIndex))
		args = append(args, *filters.EmployeeID)
		argIndex++
	}

	if filters.Search != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("(original_filename ILIKE $%d OR filename ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+*filters.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("%s", whereConditions[0])
		for i := 1; i < len(whereConditions); i++ {
			whereClause += " AND " + whereConditions[i]
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM uploaded_files %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count uploaded files: %w", err)
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
	limit := 20 // default page size
	if filters.PageSize > 0 {
		limit = filters.PageSize
	}
	offset := 0
	if filters.Page > 0 {
		offset = (filters.Page - 1) * limit
	}

	// Get files
	query := fmt.Sprintf(`
		SELECT id, filename, original_filename, file_path, file_size, content_type,
		       file_hash, upload_type, status, uploaded_by, employee_id, metadata,
		       created_at, updated_at, deleted_at
		FROM uploaded_files
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list uploaded files: %w", err)
	}
	defer rows.Close()

	var files []models.UploadedFile
	for rows.Next() {
		var file models.UploadedFile
		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.OriginalFilename,
			&file.FilePath,
			&file.FileSize,
			&file.ContentType,
			&file.FileHash,
			&file.UploadType,
			&file.Status,
			&file.UploadedBy,
			&file.EmployeeID,
			&file.Metadata,
			&file.CreatedAt,
			&file.UpdatedAt,
			&file.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan uploaded file: %w", err)
		}
		files = append(files, file)
	}

	return files, total, nil
}

// GetStats retrieves file upload statistics
func (r *uploadedFileRepository) GetStats() (*models.FileUploadStats, error) {
	stats := &models.FileUploadStats{
		FilesByType:   make(map[string]int64),
		FilesByStatus: make(map[string]int64),
	}

	// Get total files and size
	err := r.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(file_size), 0)
		FROM uploaded_files
		WHERE deleted_at IS NULL`).Scan(&stats.TotalFiles, &stats.TotalSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}

	// Get files by type
	rows, err := r.db.Query(`
		SELECT upload_type, COUNT(*)
		FROM uploaded_files
		WHERE deleted_at IS NULL
		GROUP BY upload_type`)
	if err != nil {
		return nil, fmt.Errorf("failed to get files by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var uploadType string
		var count int64
		if err := rows.Scan(&uploadType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan files by type: %w", err)
		}
		stats.FilesByType[uploadType] = count
	}

	// Get files by status
	rows, err = r.db.Query(`
		SELECT status, COUNT(*)
		FROM uploaded_files
		WHERE deleted_at IS NULL
		GROUP BY status`)
	if err != nil {
		return nil, fmt.Errorf("failed to get files by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan files by status: %w", err)
		}
		stats.FilesByStatus[status] = count
	}

	// Get recent uploads (last 10)
	recentFiles, _, err := r.List(FileListFilters{
		PageSize:  10,
		SortBy:    "created_at",
		SortOrder: "desc",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recent uploads: %w", err)
	}
	stats.RecentUploads = recentFiles

	return stats, nil
}

// GetByEmployeeID retrieves files associated with an employee
func (r *uploadedFileRepository) GetByEmployeeID(employeeID int) ([]models.UploadedFile, error) {
	files, _, err := r.List(FileListFilters{
		EmployeeID: &employeeID,
		PageSize:   1000, // Get all files for the employee
	})
	return files, err
}

// GetByUploadType retrieves files by upload type
func (r *uploadedFileRepository) GetByUploadType(uploadType string) ([]models.UploadedFile, error) {
	files, _, err := r.List(FileListFilters{
		UploadType: &uploadType,
		PageSize:   1000, // Get all files of this type
	})
	return files, err
}

// CleanupOldFiles removes old deleted files
func (r *uploadedFileRepository) CleanupOldFiles(olderThan time.Time) (int64, error) {
	query := `DELETE FROM uploaded_files WHERE deleted_at IS NOT NULL AND deleted_at < $1`
	result, err := r.db.Exec(query, olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old files: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
