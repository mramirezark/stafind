package handlers

import (
	"strconv"
	"strings"

	"stafind-backend/internal/auth"
	"stafind-backend/internal/middleware"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// getValidationErrors returns formatted validation errors
func getValidationErrors(err error) []string {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, e.Error())
		}
	} else {
		errors = append(errors, err.Error())
	}
	return errors
}

// handleServiceError handles service errors and returns appropriate HTTP responses
func handleServiceError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *services.ValidationError:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": e.Error(),
		})
	case *services.NotFoundError:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": e.Error(),
		})
	case *services.ConflictError:
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": e.Error(),
		})
	default:
		// Check if it's a string error
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "validation") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
}

// AuthHandlers handles authentication-related HTTP requests
type AuthHandlers struct {
	userService services.UserService
	roleService services.RoleService
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(userService services.UserService, roleService services.RoleService) *AuthHandlers {
	return &AuthHandlers{
		userService: userService,
		roleService: roleService,
	}
}

// Login handles user login
func (h *AuthHandlers) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	// Authenticate user
	response, err := h.userService.Login(&req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(response)
}

// Logout handles user logout
func (h *AuthHandlers) Logout(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Extract token
	tokenParts := authHeader[7:] // Remove "Bearer " prefix
	tokenHash := auth.GenerateTokenHash(tokenParts)

	// Revoke session
	err := h.userService.Logout(tokenHash)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

// Register handles user registration
func (h *AuthHandlers) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	// Create user
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// RefreshToken handles token refresh
func (h *AuthHandlers) RefreshToken(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Extract token
	tokenParts := authHeader[7:] // Remove "Bearer " prefix

	// Refresh token
	newToken, err := h.userService.RefreshToken(tokenParts)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"token": newToken,
	})
}

// GetProfile retrieves the current user's profile
func (h *AuthHandlers) GetProfile(c *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		return handleServiceError(c, err)
	}

	profile, err := h.userService.GetUserProfile(user.ID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(profile)
}

// UpdateProfile updates the current user's profile
func (h *AuthHandlers) UpdateProfile(c *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		return handleServiceError(c, err)
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	updatedUser, err := h.userService.UpdateUserProfile(user.ID, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(updatedUser)
}

// ChangePassword handles password change
func (h *AuthHandlers) ChangePassword(c *fiber.Ctx) error {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		return handleServiceError(c, err)
	}

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	err = h.userService.ChangePassword(user.ID, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// ListUsers handles listing users (admin only)
func (h *AuthHandlers) ListUsers(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.ListUsers(page, limit)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetUser handles getting a specific user (admin only)
func (h *AuthHandlers) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(user)
}

// UpdateUser handles updating a specific user (admin only)
func (h *AuthHandlers) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	updatedUser, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(updatedUser)
}

// DeleteUser handles deleting a user (admin only)
func (h *AuthHandlers) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// ListRoles handles listing all roles
func (h *AuthHandlers) ListRoles(c *fiber.Ctx) error {
	roles, err := h.roleService.ListRoles()
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"roles": roles,
	})
}

// GetRole handles getting a specific role
func (h *AuthHandlers) GetRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(role)
}

// CreateRole handles creating a new role (admin only)
func (h *AuthHandlers) CreateRole(c *fiber.Ctx) error {
	var req models.Role
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

// UpdateRole handles updating a role (admin only)
func (h *AuthHandlers) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	var req models.Role
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": getValidationErrors(err),
		})
	}

	role, err := h.roleService.UpdateRole(id, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(role)
}

// DeleteRole handles deleting a role (admin only)
func (h *AuthHandlers) DeleteRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	err = h.roleService.DeleteRole(id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}
