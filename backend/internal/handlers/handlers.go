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
	AIAgentHandlers *AIAgentHandlers
	NERHandlers     *NERHandlers
}

func NewHandlers(
	employeeService services.EmployeeService,
	searchService services.SearchService,
	skillService services.SkillService,
	aiAgentService services.AIAgentService,
	nerService *services.NERService,
) *Handlers {
	return &Handlers{
		employeeService: employeeService,
		searchService:   searchService,
		skillService:    skillService,
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
	var skill models.Skill
	if err := c.BodyParser(&skill); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdSkill, err := h.skillService.CreateSkill(&skill)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Skill already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdSkill)
}
