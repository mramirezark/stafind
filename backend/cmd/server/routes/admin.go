package routes

import (
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes configures admin-only routes with authentication and admin role requirement
func SetupAdminRoutes(app *fiber.App, authHandlers *handlers.AuthHandlers, apiKeyHandlers *handlers.APIKeyHandlers) {
	// Admin routes with authentication and admin role requirement
	admin := app.Group("/api/v1/admin", middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		// User management routes
		admin.Get("/users", authHandlers.ListUsers)
		admin.Get("/users/:id", authHandlers.GetUser)
		admin.Put("/users/:id", authHandlers.UpdateUser)
		admin.Delete("/users/:id", authHandlers.DeleteUser)

		// Role management routes
		admin.Get("/roles", authHandlers.ListRoles)
		admin.Get("/roles/:id", authHandlers.GetRole)
		admin.Post("/roles", authHandlers.CreateRole)
		admin.Put("/roles/:id", authHandlers.UpdateRole)
		admin.Delete("/roles/:id", authHandlers.DeleteRole)

		// API Key management routes
		admin.Get("/api-keys", apiKeyHandlers.GetAPIKeys)
		admin.Get("/api-keys/:id", apiKeyHandlers.GetAPIKey)
		admin.Post("/api-keys", apiKeyHandlers.CreateAPIKey)
		admin.Post("/api-keys/:id/rotate", apiKeyHandlers.RotateAPIKey)
		admin.Post("/api-keys/:id/deactivate", apiKeyHandlers.DeactivateAPIKey)
	}
}
