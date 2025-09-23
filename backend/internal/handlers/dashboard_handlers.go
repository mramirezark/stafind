package handlers

import (
	"stafind-backend/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandlers struct {
	dashboardService services.DashboardService
}

func NewDashboardHandlers(dashboardService services.DashboardService) *DashboardHandlers {
	return &DashboardHandlers{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats returns dashboard statistics
func (h *DashboardHandlers) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.dashboardService.GetDashboardStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

// GetRecentJobRequests returns recent job requests
func (h *DashboardHandlers) GetRecentJobRequests(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	jobRequests, err := h.dashboardService.GetRecentJobRequests(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(jobRequests)
}

// GetRecentEmployees returns recent employees
func (h *DashboardHandlers) GetRecentEmployees(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	employees, err := h.dashboardService.GetRecentEmployees(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(employees)
}

// GetDepartmentStats returns department statistics
func (h *DashboardHandlers) GetDepartmentStats(c *fiber.Ctx) error {
	stats, err := h.dashboardService.GetDepartmentStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

// GetSkillDemandStats returns skill demand statistics
func (h *DashboardHandlers) GetSkillDemandStats(c *fiber.Ctx) error {
	stats, err := h.dashboardService.GetSkillDemandStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}
