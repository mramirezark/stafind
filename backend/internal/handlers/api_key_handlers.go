package handlers

import (
	"stafind-backend/internal/constants"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type APIKeyHandlers struct {
	apiKeyService services.APIKeyService
}

func NewAPIKeyHandlers(apiKeyService services.APIKeyService) *APIKeyHandlers {
	return &APIKeyHandlers{
		apiKeyService: apiKeyService,
	}
}

// CreateAPIKey creates a new API key
func (h *APIKeyHandlers) CreateAPIKey(c *fiber.Ctx) error {
	var req models.CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	apiKey, err := h.apiKeyService.CreateAPIKey(&req)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(constants.StatusCreated).JSON(apiKey)
}

// GetAPIKeys returns all API keys
func (h *APIKeyHandlers) GetAPIKeys(c *fiber.Ctx) error {
	limitStr := c.Query(constants.ParamLimit, "10")
	offsetStr := c.Query(constants.ParamOffset, "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = constants.DefaultPageSize
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = constants.DefaultOffset
	}

	apiKeys, err := h.apiKeyService.GetAPIKeys(limit, offset)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(apiKeys)
}

// GetAPIKey returns a specific API key by ID
func (h *APIKeyHandlers) GetAPIKey(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{"error": constants.MsgInvalidAPIKeyID})
	}

	// Note: This would need to be implemented in the service
	// For now, return a placeholder
	return c.JSON(fiber.Map{
		"id":      id,
		"message": "API key details endpoint - to be implemented",
	})
}

// DeactivateAPIKey deactivates an API key
func (h *APIKeyHandlers) DeactivateAPIKey(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{"error": constants.MsgInvalidAPIKeyID})
	}

	err = h.apiKeyService.DeactivateAPIKey(id)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": constants.MsgAPIKeyDeactivated})
}

// RotateAPIKey rotates an API key (creates new, deactivates old)
func (h *APIKeyHandlers) RotateAPIKey(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{"error": constants.MsgInvalidAPIKeyID})
	}

	response, err := h.apiKeyService.RotateAPIKey(id)
	if err != nil {
		return c.Status(constants.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": constants.MsgAPIKeyRotated,
		"data":    response,
	})
}

// ValidateAPIKey validates an API key
func (h *APIKeyHandlers) ValidateAPIKey(c *fiber.Ctx) error {
	apiKey := c.Get(constants.HeaderAPIKey)
	if apiKey == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error": constants.MsgAPIKeyRequired,
			"code":  constants.ErrorCodeMissingAPIKey,
		})
	}

	keyInfo, err := h.apiKeyService.ValidateAPIKey(apiKey)
	if err != nil {
		return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
			"error": constants.MsgInvalidAPIKey,
			"code":  constants.ErrorCodeInvalidAPIKey,
		})
	}

	return c.JSON(fiber.Map{
		"valid":        true,
		"service_name": keyInfo.ServiceName,
		"description":  keyInfo.Description,
		"is_active":    keyInfo.IsActive,
		"created_at":   keyInfo.CreatedAt,
		"expires_at":   keyInfo.ExpiresAt,
		"last_used_at": keyInfo.LastUsedAt,
	})
}

// TestAPIKey tests API key authentication
func (h *APIKeyHandlers) TestAPIKey(c *fiber.Ctx) error {
	apiKey := c.Get(constants.HeaderAPIKey)
	if apiKey == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"error": constants.MsgAPIKeyRequired,
			"code":  constants.ErrorCodeMissingAPIKey,
		})
	}

	keyInfo, err := h.apiKeyService.ValidateAPIKey(apiKey)
	if err != nil {
		return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
			"error": constants.MsgInvalidAPIKey,
			"code":  constants.ErrorCodeInvalidAPIKey,
		})
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"message":      constants.MsgAPIKeyValid,
		"service_name": keyInfo.ServiceName,
		"description":  keyInfo.Description,
		"is_active":    keyInfo.IsActive,
	})
}
