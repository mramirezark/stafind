package handlers

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"stafind-backend/internal/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CVExtractHandlers struct {
	extractService services.CVExtractService
}

func NewCVExtractHandlers(extractService services.CVExtractService) *CVExtractHandlers {
	return &CVExtractHandlers{
		extractService: extractService,
	}
}

// CreateOrUpdateExtract creates or updates a CV extraction record
func (h *CVExtractHandlers) CreateOrUpdateExtract(c *fiber.Ctx) error {
	var request struct {
		ExtractRequestID string  `json:"extract_request_id" validate:"required"`
		Status           string  `json:"status" validate:"required"`
		NumFiles         int     `json:"num_files" validate:"required,min=0"`
		FileNumber       int     `json:"file_number" validate:"required,min=0"`
		Metadata         *string `json:"metadata,omitempty"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	extract, err := h.extractService.CreateOrUpdateExtract(
		request.ExtractRequestID,
		request.Status,
		request.NumFiles,
		request.FileNumber,
		request.Metadata,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create or update CV extract",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "CV extract created or updated successfully",
		"data":    extract,
	})
}

// GetExtractByRequestID retrieves a CV extraction by request ID
func (h *CVExtractHandlers) GetExtractByRequestID(c *fiber.Ctx) error {
	requestID := c.Params("requestId")
	if requestID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Extract request ID is required",
		})
	}

	extract, err := h.extractService.GetExtractByRequestID(requestID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "CV extract not found",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": extract,
	})
}

// UpdateExtractStatus updates the status of a CV extraction
func (h *CVExtractHandlers) UpdateExtractStatus(c *fiber.Ctx) error {
	requestID := c.Params("requestId")
	if requestID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Extract request ID is required",
		})
	}

	var request struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	extract, err := h.extractService.UpdateExtractStatus(requestID, request.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update CV extract status",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "CV extract status updated successfully",
		"data":    extract,
	})
}

// UpdateExtractProgress updates the progress of a CV extraction
func (h *CVExtractHandlers) UpdateExtractProgress(c *fiber.Ctx) error {
	requestID := c.Params("requestId")
	if requestID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Extract request ID is required",
		})
	}

	var request struct {
		FilesProcessed int `json:"files_processed" validate:"required,min=0"`
		FilesFailed    int `json:"files_failed" validate:"required,min=0"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	extract, err := h.extractService.UpdateExtractProgress(
		requestID,
		request.FilesProcessed,
		request.FilesFailed,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update CV extract progress",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "CV extract progress updated successfully",
		"data":    extract,
	})
}

// CompleteExtract marks a CV extraction as completed
func (h *CVExtractHandlers) CompleteExtract(c *fiber.Ctx) error {
	requestID := c.Params("requestId")
	if requestID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Extract request ID is required",
		})
	}

	var request struct {
		TotalTimeMs  int64   `json:"total_time_ms" validate:"required,min=0"`
		ErrorMessage *string `json:"error_message,omitempty"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	extract, err := h.extractService.CompleteExtract(
		requestID,
		request.TotalTimeMs,
		request.ErrorMessage,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to complete CV extract",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "CV extract completed successfully",
		"data":    extract,
	})
}

// GetExtractStats retrieves CV extraction statistics
func (h *CVExtractHandlers) GetExtractStats(c *fiber.Ctx) error {
	// Parse query parameters for date filters
	var startDate, endDate *time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	filters := repositories.CVExtractStatsFilters{
		StartDate: startDate,
		EndDate:   endDate,
	}

	stats, err := h.extractService.GetExtractStats(filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get CV extract stats",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": stats,
	})
}

// ListExtracts retrieves a paginated list of CV extractions
func (h *CVExtractHandlers) ListExtracts(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	status := c.Query("status")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Parse date filters
	var startDate, endDate *time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	filters := repositories.CVExtractFilters{
		Status:    &status,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	extracts, total, err := h.extractService.ListExtracts(filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to list CV extracts",
			"details": err.Error(),
		})
	}

	// Calculate pagination info
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := models.CVExtractListResponse{
		Records:    extracts,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return c.Status(200).JSON(fiber.Map{
		"data": response,
	})
}

// GetRecentExtracts retrieves recent CV extraction records
func (h *CVExtractHandlers) GetRecentExtracts(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	extracts, err := h.extractService.GetRecentExtracts(limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get recent CV extracts",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": extracts,
	})
}

// RegisterCVExtractRoutes registers CV extract routes
func (h *CVExtractHandlers) RegisterCVExtractRoutes(app *fiber.App) {
	// CV extract routes
	api := app.Group("/api/v1/cv-extract")
	api.Post("/", h.CreateOrUpdateExtract)
	api.Get("/stats", h.GetExtractStats)
	api.Get("/", h.ListExtracts)
	api.Get("/recent", h.GetRecentExtracts)
	api.Get("/:requestId", h.GetExtractByRequestID)
	api.Put("/:requestId/status", h.UpdateExtractStatus)
	api.Put("/:requestId/progress", h.UpdateExtractProgress)
	api.Put("/:requestId/complete", h.CompleteExtract)
}
