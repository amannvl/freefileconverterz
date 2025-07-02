package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected protects routes with JWT authentication
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// TODO: Replace with your actual secret key from config
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Add user info to context
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("userID", claims["user_id"])
		c.Locals("userRole", claims["role"])

		return c.Next()
	}
}

// RequireRole requires the user to have a specific role
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}
		return c.Next()
	}
}

// ErrorHandler is defined in error_handler.go

// CORSMiddleware handles CORS headers
func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.Next()
	}
}

// RateLimitMiddleware limits the number of requests
func RateLimitMiddleware() fiber.Handler {
	// TODO: Implement rate limiting using Redis
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// FileUploadLimiter returns a middleware that sets the Content-Length header
// Note: The actual body limit should be set in the Fiber app configuration
func FileUploadLimiter(maxSize int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set the Content-Length header
		c.Set("Content-Length", strconv.FormatInt(maxSize, 10))
		return c.Next()
	}
}
