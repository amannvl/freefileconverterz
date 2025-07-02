package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// Error represents an error that can be returned as JSON
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error,omitempty"`
}

// ErrorHandler is the global error handler middleware
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a fiber error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	// Log the error
	slog.Error("Request error",
		"method", c.Method(),
		"path", c.Path(),
		"status", code,
		"error", err,
	)

	// Return JSON response
	return c.Status(code).JSON(ErrorResponse{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

// NotFoundHandler handles 404 responses
func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Success: false,
		Error: &Error{
			Code:    fiber.StatusNotFound,
			Message: "Endpoint not found",
		},
	})
}

// MethodNotAllowedHandler handles 405 responses
func MethodNotAllowedHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusMethodNotAllowed).JSON(ErrorResponse{
		Success: false,
		Error: &Error{
			Code:    fiber.StatusMethodNotAllowed,
			Message: "Method not allowed",
		},
	})
}

// ErrorLogger logs errors that occur during request handling
func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue to next middleware/handler
		err := c.Next()

		// Log any errors that occurred
		if err != nil {
			slog.Error("Request failed",
				"method", c.Method(),
				"path", c.Path(),
				"status", c.Response().StatusCode(),
				"error", err,
			)
		}

		return err
	}
}

// RecoverHandler recovers from panics and logs the error
func RecoverHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fiber.NewError(fiber.StatusInternalServerError, http.StatusText(fiber.StatusInternalServerError))
				}

				slog.Error("Recovered from panic",
					"method", c.Method(),
					"path", c.Path(),
					"error", err,
					"stack", string(debug.Stack()),
				)

				_ = c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Success: false,
					Error: &Error{
						Code:    fiber.StatusInternalServerError,
						Message: "Internal Server Error",
					},
				})
			}
		}()

		return c.Next()
	}
}
