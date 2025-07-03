package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// handleHealthCheck handles the health check endpoint
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health [get]
func (s *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:  "ok",
		Service: "freefileconverterz",
		Version: "1.0.0",
	})
}
