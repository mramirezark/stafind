package routes

import (
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAPIRoutes configures protected API routes with authentication
func SetupAPIRoutes(
	app *fiber.App,
	h *handlers.Handlers,
	authHandlers *handlers.AuthHandlers,
	dashboardHandlers *handlers.DashboardHandlers,
	fileUploadHandlers *handlers.FileUploadHandlers,
	apiKeyHandlers *handlers.APIKeyHandlers,
) {
	// Protected API routes group
	api := app.Group("/api/v1", middleware.AuthMiddleware())
	{
		// Employee routes
		api.Get("/employees", h.GetEmployees)
		api.Get("/employees/:id", h.GetEmployee)
		api.Post("/employees", h.CreateEmployee)
		api.Put("/employees/:id", h.UpdateEmployee)
		api.Delete("/employees/:id", h.DeleteEmployee)

		// Search routes
		api.Post("/search", h.SearchEmployees)

		// Skill routes - specific routes first, then wildcard routes
		api.Get("/skills", h.GetSkills)
		api.Post("/skills", h.CreateSkill)
		api.Get("/skills/search", h.SearchSkills)
		api.Get("/skills/popular", h.GetPopularSkills)
		api.Get("/skills/with-count", h.GetSkillsWithCount)
		api.Get("/skills/stats", h.GetSkillStats)

		// Category routes - specific routes before wildcard
		api.Get("/skills/categories", h.GetCategories)
		api.Post("/skills/categories", h.CreateCategory)
		api.Get("/skills/categories/:id", h.GetCategory)
		api.Put("/skills/categories/:id", h.UpdateCategory)
		api.Delete("/skills/categories/:id", h.DeleteCategory)

		// Wildcard skill routes - must come after specific routes
		api.Get("/skills/:id", h.GetSkill)
		api.Put("/skills/:id", h.UpdateSkill)
		api.Delete("/skills/:id", h.DeleteSkill)

		// Role routes
		api.Get("/roles", authHandlers.ListRoles)

		// AI Agent routes
		api.Get("/ai-agent/requests", h.AIAgentHandlers.GetAIAgentRequests)
		api.Get("/ai-agent/requests/:id", h.AIAgentHandlers.GetAIAgentRequest)
		api.Get("/ai-agent/responses/:id", h.AIAgentHandlers.GetAIAgentResponse)
		api.Get("/ai-agents/:id", h.AIAgentHandlers.GetAIAgentRequest)
		api.Get("/ai-agents/:id/response", h.AIAgentHandlers.GetAIAgentResponse)
		api.Post("/ai-agents/:id/process", h.AIAgentHandlers.ProcessAIAgentRequest)
		api.Post("/ai-agent/process", h.AIAgentHandlers.ProcessAIAgentRequest)
		api.Post("/ai-agents/:id/process-by-id", h.AIAgentHandlers.ProcessAIAgentRequestByID)
		api.Post("/ai-agents/extract-skills", h.AIAgentHandlers.ExtractSkills)

		// Dashboard routes
		api.Get("/dashboard/stats", dashboardHandlers.GetDashboardStats)
		api.Get("/dashboard/metrics", dashboardHandlers.GetDashboardMetrics)
		api.Get("/dashboard/recent-employees", dashboardHandlers.GetRecentEmployees)
		api.Get("/dashboard/department-stats", dashboardHandlers.GetDepartmentStats)
		api.Get("/dashboard/skill-demand-stats", dashboardHandlers.GetSkillDemandStats)
		api.Get("/dashboard/skill-demand", dashboardHandlers.GetSkillDemandStats) // Alias for skill-demand-stats
		api.Get("/dashboard/top-suggested-employees", dashboardHandlers.GetTopSuggestedEmployees)
		api.Get("/dashboard/top-employees", dashboardHandlers.GetTopSuggestedEmployees) // Alias for top-suggested-employees

		// File upload routes
		api.Post("/upload", fileUploadHandlers.UploadFile)
		api.Get("/files/:id", fileUploadHandlers.GetFile)
		api.Delete("/files/:id", fileUploadHandlers.DeleteFile)
	}
}
