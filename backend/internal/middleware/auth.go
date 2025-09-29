package middleware

import (
	"strings"

	"stafind-backend/internal/auth"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if token starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set user context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_first_name", claims.FirstName)
		c.Locals("user_last_name", claims.LastName)
		c.Locals("user_roles", claims.Roles)

		return c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user is authenticated
		userRoles, ok := c.Locals("user_roles").([]string)
		if !ok {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		// Check if user has required role
		hasRole := false
		for _, role := range userRoles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(constants.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAnyRole middleware checks if user has any of the required roles
func RequireAnyRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user is authenticated
		userRoles, ok := c.Locals("user_roles").([]string)
		if !ok {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return c.Status(constants.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// OptionalAuth middleware validates JWT tokens if present but doesn't require them
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Check if token starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Next()
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return c.Next() // Continue without authentication
		}

		// Set user context if token is valid
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_first_name", claims.FirstName)
		c.Locals("user_last_name", claims.LastName)
		c.Locals("user_roles", claims.Roles)

		return c.Next()
	}
}

// GetCurrentUser extracts the current user from context
func GetCurrentUser(c *fiber.Ctx) (*models.User, error) {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return nil, fiber.NewError(constants.StatusUnauthorized, "User not authenticated")
	}

	email, ok := c.Locals("user_email").(string)
	if !ok {
		return nil, fiber.NewError(constants.StatusUnauthorized, "User email not found")
	}

	firstName, ok := c.Locals("user_first_name").(string)
	if !ok {
		return nil, fiber.NewError(constants.StatusUnauthorized, "User first name not found")
	}

	lastName, ok := c.Locals("user_last_name").(string)
	if !ok {
		return nil, fiber.NewError(constants.StatusUnauthorized, "User last name not found")
	}

	roles, ok := c.Locals("user_roles").([]string)
	if !ok {
		roles = []string{}
	}

	// Convert string roles to Role structs
	var roleStructs []models.Role
	for _, roleName := range roles {
		roleStructs = append(roleStructs, models.Role{Name: roleName})
	}

	return &models.User{
		ID:        userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Roles:     roleStructs,
	}, nil
}

// HasRole checks if the current user has a specific role
func HasRole(c *fiber.Ctx, roleName string) bool {
	userRoles, ok := c.Locals("user_roles").([]string)
	if !ok {
		return false
	}

	for _, role := range userRoles {
		if role == roleName {
			return true
		}
	}
	return false
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *fiber.Ctx) bool {
	return HasRole(c, "admin")
}

// IsHRManager checks if the current user is an HR manager
func IsHRManager(c *fiber.Ctx) bool {
	return HasRole(c, "hr_manager")
}

// IsHiringManager checks if the current user is a hiring manager
func IsHiringManager(c *fiber.Ctx) bool {
	return HasRole(c, "hiring_manager")
}
