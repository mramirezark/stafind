package handlers

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type Handlers struct {
	employeeService   services.EmployeeService
	jobRequestService services.JobRequestService
	searchService     services.SearchService
	skillService      services.SkillService
}

func NewHandlers(
	employeeService services.EmployeeService,
	jobRequestService services.JobRequestService,
	searchService services.SearchService,
	skillService services.SkillService,
) *Handlers {
	return &Handlers{
		employeeService:   employeeService,
		jobRequestService: jobRequestService,
		searchService:     searchService,
		skillService:      skillService,
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

// GetJobRequests returns all job requests
func (h *Handlers) GetJobRequests(c *fiber.Ctx) error {
	requests, err := h.jobRequestService.GetAllJobRequests()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(requests)
}

// GetJobRequest returns a specific job request by ID
func (h *Handlers) GetJobRequest(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job request ID"})
	}

	request, err := h.jobRequestService.GetJobRequestByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job request not found"})
	}
	return c.JSON(request)
}

// CreateJobRequest creates a new job request
func (h *Handlers) CreateJobRequest(c *fiber.Ctx) error {
	var req models.CreateJobRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	jobRequest, err := h.jobRequestService.CreateJobRequest(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(jobRequest)
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

// GetMatchesForJobRequest returns matches for a specific job request
func (h *Handlers) GetMatchesForJobRequest(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job request ID"})
	}

	matches, err := h.searchService.FindMatchesForJobRequest(id)
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
