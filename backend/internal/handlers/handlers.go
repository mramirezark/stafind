package handlers

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type Handlers struct {
	employeeService services.EmployeeService
	searchService   services.SearchService
	skillService    services.SkillService
	categoryService services.CategoryService
	AIAgentHandlers *AIAgentHandlers
	NERHandlers     *NERHandlers
}

func NewHandlers(
	employeeService services.EmployeeService,
	searchService services.SearchService,
	skillService services.SkillService,
	categoryService services.CategoryService,
	aiAgentService services.AIAgentService,
	nerService *services.NERService,
) *Handlers {
	return &Handlers{
		employeeService: employeeService,
		searchService:   searchService,
		skillService:    skillService,
		categoryService: categoryService,
		AIAgentHandlers: NewAIAgentHandlers(aiAgentService),
		NERHandlers:     NewNERHandlers(nerService, searchService),
	}
}

// GetEmployees returns all employees
func (h *Handlers) GetEmployees(c *fiber.Ctx) error {
	employees, err := h.employeeService.GetAllEmployees()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(employees)
}

// GetEmployee returns a specific employee by ID
func (h *Handlers) GetEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	employee, err := h.employeeService.GetEmployeeByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employee not found"})
	}
	return c.JSON(employee)
}

// CreateEmployee creates a new employee
func (h *Handlers) CreateEmployee(c *fiber.Ctx) error {
	var req models.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	employee, err := h.employeeService.CreateEmployee(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(employee)
}

// UpdateEmployee updates an existing employee
func (h *Handlers) UpdateEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	var req models.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	employee, err := h.employeeService.UpdateEmployee(id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(employee)
}

// DeleteEmployee deletes an employee
func (h *Handlers) DeleteEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	err = h.employeeService.DeleteEmployee(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Employee deleted successfully"})
}

// SearchEmployees searches for employees based on criteria
func (h *Handlers) SearchEmployees(c *fiber.Ctx) error {
	var searchReq models.SearchRequest
	if err := c.BodyParser(&searchReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	matches, err := h.searchService.SearchEmployees(&searchReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(matches)
}

// GetSkills returns all available skills
func (h *Handlers) GetSkills(c *fiber.Ctx) error {
	skills, err := h.skillService.GetAllSkills()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(skills)
}

// CreateSkill creates a new skill
func (h *Handlers) CreateSkill(c *fiber.Ctx) error {
	var req models.CreateSkillRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdSkill, err := h.skillService.CreateSkillWithCategories(&req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Skill already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdSkill)
}

// GetSkill returns a specific skill by ID
func (h *Handlers) GetSkill(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	skill, err := h.skillService.GetSkillByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Skill not found"})
	}
	return c.JSON(skill)
}

// UpdateSkill updates an existing skill
func (h *Handlers) UpdateSkill(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	var req models.CreateSkillRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedSkill, err := h.skillService.UpdateSkillWithCategories(id, &req)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Skill already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedSkill)
}

// DeleteSkill deletes a skill
func (h *Handlers) DeleteSkill(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	err = h.skillService.DeleteSkill(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Skill deleted successfully"})
}

// SearchSkills searches for skills by name
func (h *Handlers) SearchSkills(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Search query is required"})
	}

	skills, err := h.skillService.SearchSkills(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(skills)
}

// GetPopularSkills returns the most popular skills
func (h *Handlers) GetPopularSkills(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	skills, err := h.skillService.GetPopularSkills(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(skills)
}

// GetSkillsWithCount returns skills with employee count
func (h *Handlers) GetSkillsWithCount(c *fiber.Ctx) error {
	skills, err := h.skillService.GetSkillsWithEmployeeCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(skills)
}

// GetSkillStats returns skill statistics
func (h *Handlers) GetSkillStats(c *fiber.Ctx) error {
	stats, err := h.skillService.GetSkillStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

// GetCategories returns all categories
func (h *Handlers) GetCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(categories)
}

// GetCategory returns a specific category by ID
func (h *Handlers) GetCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(category)
}

// CreateCategory creates a new category
func (h *Handlers) CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdCategory, err := h.categoryService.CreateCategory(&category)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Category already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdCategory)
}

// UpdateCategory updates an existing category
func (h *Handlers) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedCategory, err := h.categoryService.UpdateCategory(id, &category)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Category already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedCategory)
}

// DeleteCategory deletes a category
func (h *Handlers) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	err = h.categoryService.DeleteCategory(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Category deleted successfully"})
}
