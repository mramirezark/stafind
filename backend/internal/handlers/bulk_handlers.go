package handlers

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type BulkHandlers struct {
	bulkEmployeeService services.BulkEmployeeService
}

func NewBulkHandlers(bulkEmployeeService services.BulkEmployeeService) *BulkHandlers {
	return &BulkHandlers{
		bulkEmployeeService: bulkEmployeeService,
	}
}

// BulkCreateEmployees handles bulk creation of employees
func (h *BulkHandlers) BulkCreateEmployees(c *fiber.Ctx) error {
	var employees []interface{}
	if err := c.BodyParser(&employees); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Convert to Employee models
	employeeModels, err := h.convertToEmployeeModels(employees)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid employee data",
			"details": err.Error(),
		})
	}

	result, err := h.bulkEmployeeService.BulkCreateEmployees(employeeModels)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create employees",
			"details": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Bulk operation completed",
		"result":  result,
	})
}

// BulkUpdateEmployees handles bulk update of employees
func (h *BulkHandlers) BulkUpdateEmployees(c *fiber.Ctx) error {
	var employees []interface{}
	if err := c.BodyParser(&employees); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Convert to Employee models
	employeeModels, err := h.convertToEmployeeModels(employees)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid employee data",
			"details": err.Error(),
		})
	}

	result, err := h.bulkEmployeeService.BulkUpdateEmployees(employeeModels)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update employees",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Bulk operation completed",
		"result":  result,
	})
}

// BulkDeleteEmployees handles bulk deletion of employees
func (h *BulkHandlers) BulkDeleteEmployees(c *fiber.Ctx) error {
	var request struct {
		EmployeeIDs []int `json:"employee_ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if len(request.EmployeeIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No employee IDs provided",
		})
	}

	result, err := h.bulkEmployeeService.BulkDeleteEmployees(request.EmployeeIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete employees",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Bulk operation completed",
		"result":  result,
	})
}

// BulkUpsertEmployees handles bulk upsert of employees (create or update)
func (h *BulkHandlers) BulkUpsertEmployees(c *fiber.Ctx) error {
	var employees []interface{}
	if err := c.BodyParser(&employees); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Convert to Employee models
	employeeModels, err := h.convertToEmployeeModels(employees)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid employee data",
			"details": err.Error(),
		})
	}

	result, err := h.bulkEmployeeService.BulkUpsertEmployees(employeeModels)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to upsert employees",
			"details": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Bulk operation completed",
		"result":  result,
	})
}

// UploadResumes handles resume upload and parsing
func (h *BulkHandlers) UploadResumes(c *fiber.Ctx) error {
	// Get uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Failed to parse multipart form",
			"details": err.Error(),
		})
	}

	files := form.File["resumes"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	// Process files (this would integrate with resume parsing service)
	// For now, return a mock response
	processedFiles := make([]map[string]interface{}, len(files))
	for i, file := range files {
		processedFiles[i] = map[string]interface{}{
			"filename": file.Filename,
			"size":     file.Size,
			"status":   "processed",
			"employee_data": map[string]interface{}{
				"name":       "John Doe",
				"email":      "john.doe@example.com",
				"skills":     []string{"JavaScript", "React", "Node.js"},
				"experience": "Senior",
			},
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Files processed successfully",
		"files":   processedFiles,
	})
}

// GetBulkOperationStatus gets the status of a bulk operation
func (h *BulkHandlers) GetBulkOperationStatus(c *fiber.Ctx) error {
	operationID := c.Params("id")
	if operationID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Operation ID is required",
		})
	}

	// This would typically query a database for operation status
	// For now, return a mock response
	return c.Status(200).JSON(fiber.Map{
		"operation_id": operationID,
		"status":       "completed",
		"progress":     100,
		"result": map[string]interface{}{
			"total_processed": 10,
			"successful":      8,
			"failed":          2,
		},
	})
}

// convertToEmployeeModels converts interface{} slice to Employee models
func (h *BulkHandlers) convertToEmployeeModels(data []interface{}) ([]models.Employee, error) {
	// This is a simplified conversion
	// In a real implementation, you'd properly map the fields
	var employees []models.Employee

	for _, item := range data {
		employeeMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fiber.NewError(400, "Invalid employee data format")
		}

		// Convert map to Employee struct
		employee := models.Employee{
			Name:       getString(employeeMap, "name"),
			Email:      getString(employeeMap, "email"),
			Department: getString(employeeMap, "department"),
			Level:      getString(employeeMap, "level"),
			Location:   getString(employeeMap, "location"),
			Bio:        getString(employeeMap, "bio"),
		}

		// Convert skills if present
		if skillsData, ok := employeeMap["skills"].([]interface{}); ok {
			for _, skillData := range skillsData {
				if skillMap, ok := skillData.(map[string]interface{}); ok {
					skill := models.Skill{
						Name: getString(skillMap, "name"),
					}
					employee.Skills = append(employee.Skills, skill)
				}
			}
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

// Helper functions for type conversion
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

// RegisterBulkRoutes registers bulk operation routes
func (h *BulkHandlers) RegisterBulkRoutes(app *fiber.App) {
	api := app.Group("/api/v1/bulk")

	// Employee bulk operations
	api.Post("/employees/create", h.BulkCreateEmployees)
	api.Put("/employees/update", h.BulkUpdateEmployees)
	api.Delete("/employees/delete", h.BulkDeleteEmployees)
	api.Post("/employees/upsert", h.BulkUpsertEmployees)

	// Resume upload
	api.Post("/resumes/upload", h.UploadResumes)

	// Operation status
	api.Get("/operations/:id/status", h.GetBulkOperationStatus)
}
