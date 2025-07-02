package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amannvl/freefileconverterz/internal/api/handlers"
	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/config"
	"github.com/amannvl/freefileconverterz/pkg/converter"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := utils.NewLogger(utils.InfoLevel)
	if cfg.Server.Environment == "development" {
		logger = utils.NewLogger(utils.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize storage
	if err := os.MkdirAll(cfg.Storage.UploadPath, 0755); err != nil {
		logger.Fatal("Failed to create upload directory: %v", err)
	}

	fileStorage, err := storage.NewFileStorage(cfg.Storage.UploadPath, cfg.Storage.MaxUploadSize)
	if err != nil {
		logger.Fatal("Failed to initialize storage: %v", err)
	}

	// Initialize converter
	conv := converter.NewConverter(logger)

	// Create cleanup ticker
	cleanupTicker := time.NewTicker(1 * time.Hour)
	defer cleanupTicker.Stop()

	// Start cleanup goroutine
	go func() {
		for range cleanupTicker.C {
			if err := fileStorage.Cleanup(24 * time.Hour); err != nil {
				logger.Error(fmt.Errorf("failed to clean up old files: %w", err))
			}
		}
	}()

	// Initialize router
	router := setupRouter(conv, fileStorage, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Errorf("server forced to shutdown: %w", err))
	}

	logger.Info("Server exiting")
}

func setupRouter(conv *converter.Converter, storage *storage.FileStorage, logger *utils.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// Request logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if query != "" {
			path = path + "?" + query
		}

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		fields := []interface{}{
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(fmt.Errorf("request error: %s", e), fields...)
			}
		} else {
			logger.Info("Request processed", fields...)
		}
	})

	// Routes
	api := router.Group("/api/v1")
	{
		convertHandler := handlers.NewConvertHandler(conv, storage, logger)
		
		// File conversion
		api.POST("/convert", convertHandler.ConvertFile)
		
		// Get supported formats
		api.GET("/formats", handlers.GetSupportedFormats)
		
		// Health check
		api.GET("/health", handlers.HealthCheck)
	}

	// Static files (for frontend if needed)
	router.Static("/static", "./static")

	// Frontend route (if serving frontend from the same server)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
			"message": fmt.Sprintf("The requested path %s was not found", c.Request.URL.Path),
		})
	})

	return router
}
