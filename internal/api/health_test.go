package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	server := &Server{
		router: gin.New(),
	}
	server.setupRoutes()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health check endpoint",
			method:         "GET",
			url:            "/health",
		expectedStatus: http.StatusOK,
		expectedBody:   `{"status":"ok","service":"freefileconverterz","version":"1.0.0"}`,
		},
		{
			name:           "API v1 health check endpoint",
			method:         "GET",
			url:            "/api/v1/health",
		expectedStatus: http.StatusOK,
		expectedBody:   `{"status":"ok","service":"freefileconverterz","version":"1.0.0"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to inspect the response
			rr := httptest.NewRecorder()

			// Serve the request
			server.router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
