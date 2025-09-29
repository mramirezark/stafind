package middleware

import (
	"os"
	"strings"

	"stafind-backend/internal/constants"

	"github.com/gofiber/fiber/v2"
)

// APIKeyMiddleware validates API key for external services
func APIKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get API key from header
		apiKey := c.Get(constants.HeaderAPIKey)
		if apiKey == "" {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": constants.MsgAPIKeyRequired,
				"code":  constants.ErrorCodeMissingAPIKey,
			})
		}

		// Validate API key
		validAPIKey := os.Getenv(constants.EnvExternalAPIKey)
		if validAPIKey == "" {
			// If no API key is set in environment, use a default for development
			validAPIKey = constants.DevAPIKey
		}

		if apiKey != validAPIKey {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": constants.MsgInvalidAPIKey,
				"code":  constants.ErrorCodeInvalidAPIKey,
			})
		}

		// Add API key info to context for logging
		c.Locals(constants.ContextAPIKey, apiKey)
		c.Locals(constants.ContextAuthType, "api_key")

		return c.Next()
	}
}

// OptionalAPIKeyMiddleware validates API key if present, but doesn't require it
func OptionalAPIKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get(constants.HeaderAPIKey)
		if apiKey != "" {
			validAPIKey := os.Getenv(constants.EnvExternalAPIKey)
			if validAPIKey == "" {
				validAPIKey = constants.DevAPIKey
			}

			if apiKey == validAPIKey {
				c.Locals(constants.ContextAPIKey, apiKey)
				c.Locals(constants.ContextAuthType, "api_key")
			}
		}

		return c.Next()
	}
}

// ServiceTokenMiddleware validates service account tokens
func ServiceTokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(constants.HeaderAuthorization)
		if authHeader == "" {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": constants.MsgAuthHeaderRequired,
				"code":  constants.ErrorCodeMissingAuthHeader,
			})
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != constants.HeaderBearer {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": constants.MsgInvalidAuthFormat,
				"code":  constants.ErrorCodeInvalidAuthFormat,
			})
		}

		token := tokenParts[1]

		// Validate service token (you would implement proper JWT validation here)
		validServiceToken := os.Getenv(constants.EnvServiceToken)
		if validServiceToken == "" {
			validServiceToken = constants.DevServiceToken
		}

		if token != validServiceToken {
			return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
				"error": constants.MsgInvalidServiceToken,
				"code":  constants.ErrorCodeInvalidServiceToken,
			})
		}

		c.Locals(constants.ContextServiceToken, token)
		c.Locals(constants.ContextAuthType, "service_token")

		return c.Next()
	}
}
