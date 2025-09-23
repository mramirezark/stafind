package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"stafind-backend/internal/services"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FileUploadHandlers struct {
	uploadDir   string
	fileService services.UploadedFileService
}

func NewFileUploadHandlers(fileService services.UploadedFileService) *FileUploadHandlers {
	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &FileUploadHandlers{
		uploadDir:   uploadDir,
		fileService: fileService,
	}
}

// UploadFile handles single file upload
func (h *FileUploadHandlers) UploadFile(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "No file provided",
			"details": err.Error(),
		})
	}

	// Validate file
	if err := h.validateFile(file); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "File validation failed",
			"details": err.Error(),
		})
	}

	// Generate unique filename
	filename := h.generateFilename(file.Filename)
	filePath := filepath.Join(h.uploadDir, filename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to save file",
			"details": err.Error(),
		})
	}

	// Calculate file hash
	fileHash, err := h.calculateFileHash(filePath)
	if err != nil {
		// Clean up the file if hash calculation fails
		os.Remove(filePath)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to calculate file hash",
			"details": err.Error(),
		})
	}

	// Create file record in database
	fileData := models.UploadedFileCreate{
		Filename:         filename,
		OriginalFilename: file.Filename,
		FilePath:         filePath,
		FileSize:         file.Size,
		ContentType:      file.Header.Get("Content-Type"),
		FileHash:         &fileHash,
		UploadType:       models.FileTypeResume,
		Status:           models.FileStatusActive,
		// TODO: Add uploaded_by from JWT token
		// TODO: Add employee_id if provided in request
	}

	uploadedFile, err := h.fileService.CreateFile(fileData)
	if err != nil {
		// Clean up the file if database save fails
		os.Remove(filePath)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to save file record",
			"details": err.Error(),
		})
	}

	// Return file info
	return c.Status(200).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"file": fiber.Map{
			"id":           uploadedFile.ID,
			"filename":     uploadedFile.Filename,
			"originalName": uploadedFile.OriginalFilename,
			"size":         uploadedFile.FileSize,
			"contentType":  uploadedFile.ContentType,
			"uploadType":   uploadedFile.UploadType,
			"status":       uploadedFile.Status,
			"url":          fmt.Sprintf("/uploads/%s", uploadedFile.Filename),
			"createdAt":    uploadedFile.CreatedAt,
		},
	})
}

// UploadMultipleFiles handles multiple file upload
func (h *FileUploadHandlers) UploadMultipleFiles(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Failed to parse multipart form",
			"details": err.Error(),
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No files provided",
		})
	}

	var uploadedFiles []fiber.Map
	var errors []string

	// Process each file
	for i, file := range files {
		// Validate file
		if err := h.validateFile(file); err != nil {
			errors = append(errors, fmt.Sprintf("File %d (%s): %v", i+1, file.Filename, err))
			continue
		}

		// Generate unique filename
		filename := h.generateFilename(file.Filename)
		filePath := filepath.Join(h.uploadDir, filename)

		// Save file
		if err := c.SaveFile(file, filePath); err != nil {
			errors = append(errors, fmt.Sprintf("File %d (%s): Failed to save - %v", i+1, file.Filename, err))
			continue
		}

		// Add to successful uploads
		uploadedFiles = append(uploadedFiles, fiber.Map{
			"filename":     filename,
			"originalName": file.Filename,
			"size":         file.Size,
			"contentType":  file.Header.Get("Content-Type"),
			"url":          fmt.Sprintf("/uploads/%s", filename),
			"key":          filename,
		})
	}

	// Return results
	response := fiber.Map{
		"message":       "Files processed",
		"uploadedFiles": uploadedFiles,
		"totalUploaded": len(uploadedFiles),
		"totalFiles":    len(files),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	return c.Status(200).JSON(response)
}

// DeleteFile handles file deletion
func (h *FileUploadHandlers) DeleteFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Filename is required",
		})
	}

	// Validate filename to prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete file",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}

// GetFile serves uploaded files
func (h *FileUploadHandlers) GetFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Filename is required",
		})
	}

	// Validate filename to prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Serve file
	return c.SendFile(filePath)
}

// ListFiles lists all uploaded files
func (h *FileUploadHandlers) ListFiles(c *fiber.Ctx) error {
	// Read directory
	entries, err := os.ReadDir(h.uploadDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to read upload directory",
			"details": err.Error(),
		})
	}

	var files []fiber.Map
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, fiber.Map{
				"filename":   entry.Name(),
				"size":       info.Size(),
				"modifiedAt": info.ModTime(),
				"url":        fmt.Sprintf("/uploads/%s", entry.Name()),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"files": files,
		"count": len(files),
	})
}

// validateFile validates uploaded file
func (h *FileUploadHandlers) validateFile(file *multipart.FileHeader) error {
	// Check file size (10MB limit)
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	// Check file type
	allowedTypes := []string{
		"application/pdf",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/msword",
	}

	contentType := file.Header.Get("Content-Type")
	allowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type %s is not allowed", contentType)
	}

	return nil
}

// generateFilename generates a unique filename
func (h *FileUploadHandlers) generateFilename(originalName string) string {
	// Get file extension
	ext := filepath.Ext(originalName)

	// Generate unique name with timestamp and random string
	timestamp := time.Now().Unix()
	randomStr := strconv.FormatInt(time.Now().UnixNano()%1000000, 36)

	return fmt.Sprintf("%d_%s%s", timestamp, randomStr, ext)
}

// calculateFileHash calculates SHA-256 hash of a file
func (h *FileUploadHandlers) calculateFileHash(filePath string) (string, error) {
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

// ListFilesWithDB retrieves files from database with pagination
func (h *FileUploadHandlers) ListFilesWithDB(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	uploadType := c.Query("upload_type")
	status := c.Query("status")
	search := c.Query("search")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Build filters
	filters := repositories.FileListFilters{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	if uploadType != "" {
		filters.UploadType = &uploadType
	}
	if status != "" {
		filters.Status = &status
	}
	if search != "" {
		filters.Search = &search
	}

	// Get files from service
	response, err := h.fileService.ListFiles(filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to list files",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(response)
}

// GetFileStats retrieves file upload statistics
func (h *FileUploadHandlers) GetFileStats(c *fiber.Ctx) error {
	stats, err := h.fileService.GetFileStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get file stats",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(stats)
}

// GetFileByID retrieves a file by ID
func (h *FileUploadHandlers) GetFileByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	file, err := h.fileService.GetFile(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "File not found",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(file)
}

// UpdateFile updates file metadata
func (h *FileUploadHandlers) UpdateFile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	var updateData models.UploadedFileUpdate
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	file, err := h.fileService.UpdateFile(id, updateData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update file",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(file)
}

// SoftDeleteFile soft deletes a file
func (h *FileUploadHandlers) SoftDeleteFile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	if err := h.fileService.SoftDeleteFile(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete file",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}

// RegisterFileUploadRoutes registers file upload routes
func (h *FileUploadHandlers) RegisterFileUploadRoutes(app *fiber.App) {
	// File upload routes
	api := app.Group("/api/v1/upload")
	api.Post("/", h.UploadFile)
	api.Post("/multiple", h.UploadMultipleFiles)
	api.Delete("/:filename", h.DeleteFile)
	api.Get("/", h.ListFiles) // Legacy endpoint for file listing

	// Database-backed file management routes
	dbApi := app.Group("/api/v1/files")
	dbApi.Get("/", h.ListFilesWithDB)
	dbApi.Get("/stats", h.GetFileStats)
	dbApi.Get("/:id", h.GetFileByID)
	dbApi.Put("/:id", h.UpdateFile)
	dbApi.Delete("/:id", h.SoftDeleteFile)

	// Static file serving
	app.Static("/uploads", "./uploads")
}
