package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Error: message,
	})
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Error: message,
	})
}

// InternalServerErrorWithDetails sends a 500 Internal Server Error response with details
func InternalServerErrorWithDetails(c *fiber.Ctx, message, details string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Error:   message,
		Details: details,
	})
}

// Success sends a 200 OK response
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
