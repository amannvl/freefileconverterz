package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HealthResponse represents the health check response structure
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime,omitempty"`
}

// HealthCheck handles the health check endpoint
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", HealthCheck)
	
	// Readiness check (includes dependencies)
	app.Get("/ready", func(c *fiber.Ctx) error {
		// TODO: Add dependency checks (database, redis, etc.)
		return c.JSON(HealthResponse{
			Status:  "ready",
			Version: "1.0.0",
		})
	})
	
	// Liveness check (basic process health)
	app.Get("/live", func(c *fiber.Ctx) error {
		return c.JSON(HealthResponse{
			Status:  "alive",
			Version: "1.0.0",
		})
	})
}
