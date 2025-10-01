package routes

import (
	"stafind-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRouteMiddleware configures route-specific middleware
func SetupRouteMiddleware(app *fiber.App) {
	// Add any route-specific middleware here
	// For example, rate limiting, logging, etc.

	// Example: Rate limiting for specific routes
	// app.Use("/api/v1/auth", limiter.New(limiter.Config{
	//     Max: 5,
	//     Expiration: 1 * time.Minute,
	// }))
}

// CreateAPIGroup creates a new API group with common middleware
func CreateAPIGroup(app *fiber.App, prefix string) fiber.Router {
	return app.Group(prefix, middleware.AuthMiddleware())
}

// CreateAdminGroup creates a new admin group with admin role requirement
func CreateAdminGroup(app *fiber.App, prefix string) fiber.Router {
	return app.Group(prefix, middleware.AuthMiddleware(), middleware.RequireRole("admin"))
}
