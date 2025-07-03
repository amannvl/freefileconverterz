package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/amannvl/freefileconverterz/internal/handlers"
	"github.com/amannvl/freefileconverterz/internal/service"
	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/factory"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

// getLogLevel converts a string log level to slog.Level
func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// setupBinaryPaths adds local binaries to PATH if they exist
func setupBinaryPaths() {
	if runtime.GOOS != "linux" {
		return
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		slog.Warn("Failed to get current working directory", "error", err)
		return
	}

	// Path to the local binaries
	binPath := filepath.Join(cwd, "bin", "linux", "amd64")

	// Check if the directory exists
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		slog.Info("Local binaries directory not found, using system PATH")
		return
	}

	// Add to PATH if not already present
	currentPath := os.Getenv("PATH")
	if !strings.Contains(currentPath, binPath) {
		os.Setenv("PATH", binPath+string(os.PathListSeparator)+currentPath)
	}

	// Set LD_LIBRARY_PATH if needed
	currentLibPath := os.Getenv("LD_LIBRARY_PATH")
	if !strings.Contains(currentLibPath, binPath) {
		os.Setenv("LD_LIBRARY_PATH", binPath+string(os.PathListSeparator)+currentLibPath)
	}

	slog.Info("Added local binaries to PATH", "path", binPath)
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup binary paths
	setupBinaryPaths()

	// Initialize logger
	log := slog.Default()
	if cfg.Logging.Level != "" {
		logLevel := getLogLevel(cfg.Logging.Level)
		handlerOpts := &slog.HandlerOptions{
			Level: logLevel,
		}

		var handler slog.Handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
		if cfg.App.Env == "development" {
			handler = slog.NewTextHandler(os.Stdout, handlerOpts)
		}

		slog.SetDefault(slog.New(handler))

		// Log binary paths
		for _, cmd := range []string{"unrar", "ffmpeg", "convert", "soffice", "7z"} {
			path, err := exec.LookPath(cmd)
			if err != nil {
				slog.Warn("Binary not found in PATH", "binary", cmd)
			} else {
				slog.Info("Using binary", "binary", cmd, "path", path)
			}
		}
	}

	// Ensure upload and temp directories exist
	if err := os.MkdirAll(cfg.Storage.UploadDir, 0755); err != nil {
		slog.Error("Failed to create upload directory", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.Storage.TempDir, 0755); err != nil {
		slog.Error("Failed to create temp directory", "error", err)
		os.Exit(1)
	}

	// Initialize template engine
	engine := html.New("./views", ".html")

	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		IdleTimeout:        30 * time.Second,
		BodyLimit:          50 * 1024 * 1024, // 50MB max request body size
		Views:             engine,
		ViewsLayout:       "layouts/main",
		DisableStartupMessage: true,
	})

	// Middleware
	app.Use(recover.New())
	// Configure CORS
	corsConfig := cors.Config{
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}

	// Only set AllowCredentials to true if we have specific origins
	if len(cfg.Security.CORSAllowedOrigins) > 0 && cfg.Security.CORSAllowedOrigins[0] != "*" {
		corsConfig.AllowOrigins = strings.Join(cfg.Security.CORSAllowedOrigins, ",")
		corsConfig.AllowCredentials = true
	} else {
		corsConfig.AllowOrigins = "*"
		corsConfig.AllowCredentials = false
	}

	app.Use(cors.New(corsConfig))

	// Initialize storage
	storageImpl, err := storage.NewStorage(cfg.Storage)
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Initialize and start cleanup service if needed
	if cfg.Storage.Provider == "local" {
		cleanupInterval := 1 * time.Hour // Run cleanup every hour
		maxFileAge := 24 * time.Hour    // Delete files older than 24 hours
		cleanupService := service.NewCleanupService(log, cfg.Storage.UploadDir, cleanupInterval, maxFileAge)
		cleanupService.Start()
		defer cleanupService.Stop()
	}

	// Initialize tool manager
	toolManager, err := tools.NewToolManager(filepath.Join(cfg.Storage.TempDir, "bin"), cfg.Storage.TempDir)
	if err != nil {
		slog.Error("Failed to create tool manager", "error", err)
		os.Exit(1)
	}

	// Log tool manager initialization
	slog.Info("Initializing tool manager...", "bin_dir", filepath.Join(cfg.Storage.TempDir, "bin"), "temp_dir", cfg.Storage.TempDir)

	// Ensure all required tools are available
	if err := toolManager.EnsureTools(); err != nil {
		slog.Error("Failed to ensure required tools are available", "error", err)
		slog.Info("Continuing with limited functionality...")
	}

	// Initialize converter factory with temp directory and tool manager
	converterFactory, err := factory.NewConverterFactory(cfg.Storage.TempDir, toolManager)
	if err != nil {
		slog.Error("Failed to create converter factory", "error", err)
		os.Exit(1)
	}

	// Setup API routes
	handlers.SetupRoutes(app, cfg, storageImpl, converterFactory, log)

	// Setup static files after API routes
	app.Static("/", "./static")
	app.Static("/assets", "./static/assets")

	// SPA fallback - serve index.html for all other routes (must be last)
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// Create a channel to listen for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	serverShutdown := make(chan error, 1)
	go func() {
		slog.Info("Starting server", "port", cfg.Server.Port, "environment", cfg.App.Env)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			serverShutdown <- err
		}
	}()

	// Wait for interrupt signal
	select {
	case err := <-serverShutdown:
		slog.Error("Server error:", "error", err)
	case sig := <-quit:
		slog.Info("Received signal, shutting down...", "signal", sig)
	}

	slog.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}

	// Clean up temp files on shutdown
	slog.Info("Cleaning up temporary files...")
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cleanupCancel()
	
	if err := storageImpl.Cleanup(cleanupCtx, 0); err != nil { // Clean up all temporary files
		slog.Error("Failed to clean up temp files", "error", err)
	}

	slog.Info("Server gracefully stopped")
}
