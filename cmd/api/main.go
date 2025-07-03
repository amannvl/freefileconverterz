package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/amannvl/freefileconverterz/internal/api/handlers"
	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/config"
	"github.com/amannvl/freefileconverterz/pkg/converter/factory"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := utils.NewLogger(utils.DebugLevel)
	if cfg.Server.Environment == "development" {
		logger = utils.NewLogger(utils.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize storage
	if err := os.MkdirAll(cfg.Storage.UploadPath, 0755); err != nil {
		logger.Fatal("Failed to create upload directory", "error", err)
	}

	fileStorage, err := storage.NewFileStorage(cfg.Storage.UploadPath, cfg.Storage.MaxUploadSize, logger)
	if err != nil {
		logger.Fatal("Failed to initialize storage: %v", err)
	}

	// Create bin and temp directories
	binDir := "./bin"
	tempDir := "./temp"

	// Initialize tool manager
	toolManager, err := tools.NewToolManager(binDir, tempDir)
	if err != nil {
		logger.Fatal("Failed to create tool manager", "error", err)
	}

	// Ensure all required tools are available
	if err := toolManager.EnsureTools(); err != nil {
		logger.Fatal("Failed to ensure required tools", "error", err)
	}

	// Initialize converter factory
	convFactory, err := factory.NewConverterFactory(tempDir, toolManager)
	if err != nil {
		logger.Fatal("Failed to create converter factory", "error", err)
	}

	// Create a converter instance that implements the Converter interface
	conv := &converterAdapter{
		factory: convFactory,
		logger:  logger,
	}

	// Create cleanup ticker
	cleanupTicker := time.NewTicker(1 * time.Hour)
	defer cleanupTicker.Stop()

	// Start cleanup goroutine
	go func() {
		for range cleanupTicker.C {
			if err := fileStorage.Cleanup(24 * time.Hour); err != nil {
				logger.Error("Failed to clean up old files", "error", err)
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
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Close the logger when the application exits
	if closer, ok := logger.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			logger.Error("Failed to close logger", "error", err)
		}
	}

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}

// converterAdapter adapts ConverterFactory to the iface.Converter interface
type converterAdapter struct {
	factory *factory.ConverterFactory
	logger  utils.Logger
}

// NewConverterAdapter creates a new converter adapter
func NewConverterAdapter(factory *factory.ConverterFactory, logger utils.Logger) *converterAdapter {
	return &converterAdapter{
		factory: factory,
		logger:  logger,
	}
}

// Convert implements the iface.Converter interface
func (c *converterAdapter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Get the appropriate converter
	sourceExt := strings.TrimPrefix(filepath.Ext(inputPath), ".")
	targetExt := strings.TrimPrefix(filepath.Ext(outputPath), ".")
	if targetExt == "" {
		return fmt.Errorf("output format must be specified")
	}

	// Get the converter for the source and target formats
	conv, err := c.factory.GetConverter(sourceExt, targetExt)
	if err != nil {
		return fmt.Errorf("failed to get converter: %w", err)
	}

	// Perform the conversion
	return conv.Convert(ctx, inputPath, outputPath)
}

// Cleanup implements the iface.Converter interface
func (c *converterAdapter) Cleanup(files ...string) error {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			c.logger.Error("Failed to clean up file", "error", err, "file", file)
		}
	}
	return nil
}

// SupportsConversion implements the iface.Converter interface
func (c *converterAdapter) SupportsConversion(sourceFormat, targetFormat string) bool {
	// Get the converter for the source and target formats
	_, err := c.factory.GetConverter(sourceFormat, targetFormat)
	return err == nil
}

// SupportedFormats returns a map of supported source formats to target formats
func (c *converterAdapter) SupportedFormats() map[string][]string {
	// TODO: Implement this by aggregating formats from all available converters
	return map[string][]string{
		"jpg":  {"png", "webp"},
		"jpeg": {"png", "webp"},
		"png":  {"jpg", "webp"},
		"webp": {"jpg", "png"},
	}
}

// ValidateOptions validates the conversion options
func (c *converterAdapter) ValidateOptions(options map[string]interface{}) error {
	// TODO: Implement proper validation based on the converter type
	if format, ok := options["format"]; !ok || format == "" {
		return fmt.Errorf("output format is required")
	}
	return nil
}

// setupRouter configures the HTTP routes
func setupRouter(conv iface.Converter, storage *storage.FileStorage, logger utils.Logger) *gin.Engine {
	router := gin.New()

	// Create handlers
	handler := handlers.NewConvertHandler(conv, storage, logger)

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

		logger.Info("Request processed",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		}
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// File conversion
		v1.POST("/convert", handler.ConvertFile)

		// TODO: Add other endpoints as they are implemented
	}

	// Frontend route (if serving frontend from the same server)
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": fmt.Sprintf("The requested path %s was not found", c.Request.URL.Path),
		})
	})

	return router
}
