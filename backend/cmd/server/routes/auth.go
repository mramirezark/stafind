package routes

import (
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(app *fiber.App, authHandlers *handlers.AuthHandlers) {
	// Public auth routes (no authentication required)
	auth := app.Group("/api/v1/auth")
	{
		auth.Post("/login", authHandlers.Login)
		auth.Post("/register", authHandlers.Register)
		auth.Post("/refresh", authHandlers.RefreshToken)
		auth.Post("/logout", authHandlers.Logout)
	}

	// Protected auth routes (authentication required)
	protectedAuth := app.Group("/api/v1/auth", middleware.AuthMiddleware())
	{
		protectedAuth.Get("/profile", authHandlers.GetProfile)
		protectedAuth.Put("/profile", authHandlers.UpdateProfile)
		protectedAuth.Post("/change-password", authHandlers.ChangePassword)
	}
}
