package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Config holds the server configuration
type Config struct {
	Port         string
	Environment  string
	FileStorage  *storage.FileStorage
	Converter    *converter.ConverterFactory
	JWTSecret    string
	RateLimit    int
	RateBurst    int
	CORSOrigins []string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Server represents the HTTP server
type Server struct {
	config *Config
	router *gin.Engine
	server *http.Server
	logger zerolog.Logger
}

// NewServer creates a new HTTP server
func NewServer(cfg *Config) (*Server, error) {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize router
	router := gin.New()

	// Create server instance
	srv := &Server{
		config: cfg,
		router: router,
		logger: zerolog.New(os.Stdout).With().
			Str("component", "http_server").
			Str("env", cfg.Environment).
			Logger(),
	}

	// Configure middleware
	srv.setupMiddleware()

	// Setup routes
	srv.setupRoutes()

	// Create HTTP server
	srv.server = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return srv, nil
}

// setupMiddleware configures the server middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware recovers from any panics
	s.router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.router.Use(cors.New(corsConfig))

	// Logging middleware
	s.router.Use(s.loggingMiddleware())

	// Rate limiting
	if s.config.RateLimit > 0 {
		s.router.Use(s.rateLimitMiddleware())
	}
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", s.handleHealthCheck)

		// File conversion
		v1.POST("/convert", s.handleConvert)
		v1.GET("/formats", s.handleGetFormats)

		// User management (future)
		// v1.POST("/register", s.handleRegister)
		// v1.POST("/login", s.handleLogin)
	}
	
	// Add a root health check endpoint
	s.router.GET("/health", s.handleHealthCheck)

	// Serve static files from the frontend
	s.router.Static("/static", "./frontend/dist/static")

	// Serve frontend SPA
	s.router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info().
		Str("port", s.config.Port).
		Str("environment", s.config.Environment).
		Msg("Starting HTTP server")

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().
				Err(err).
				Msg("Failed to start server")
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info().Msg("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info().Msg("Server stopped")
	return nil
}

// handleConvert handles file conversion requests
func (s *Server) handleConvert(c *gin.Context) {
	// Get the file from the request
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get file from request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file provided",
		})
		return
	}
	defer file.Close()

	// Get the target format from the request
	targetFormat := c.PostForm("format")
	if targetFormat == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No target format provided",
		})
		return
	}

	// Get the file extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not determine file type from extension",
		})
		return
	}

	// Remove the leading dot from the extension
	ext = strings.TrimPrefix(ext, ".")

	// Get the appropriate converter
	converter, err := s.config.Converter.GetConverter(ext, targetFormat)
	if err != nil {
		s.logger.Error().Err(err).Str("source", ext).Str("target", targetFormat).Msg("Failed to get converter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Conversion from %s to %s is not supported", ext, targetFormat),
		})
		return
	}

	// Create a temporary directory for the conversion
	tempDir, err := os.MkdirTemp("", "convert_*")
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create temp directory")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary directory",
		})
		return
	}
	defer os.RemoveAll(tempDir)

	// Save the uploaded file to a temporary location
	srcPath := filepath.Join(tempDir, filepath.Base(fileHeader.Filename))
	dst, err := os.Create(srcPath)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create temp file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary file",
		})
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		s.logger.Error().Err(err).Msg("Failed to save uploaded file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save uploaded file",
		})
		return
	}
	dst.Close()

	// Create output path
	outputFile := fmt.Sprintf("%s.%s", filepath.Base(fileHeader.Filename), targetFormat)
	dstPath := filepath.Join(tempDir, outputFile)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	// Perform the conversion
	err = converter.Convert(ctx, srcPath, dstPath)
	if err != nil {
		s.logger.Error().Err(err).Str("source", ext).Str("target", targetFormat).Msg("Conversion failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Conversion failed: %v", err),
		})
		return
	}

	// Read the converted file
	result, err := os.ReadFile(dstPath)
	if err != nil {
		s.logger.Error().Err(err).Str("path", dstPath).Msg("Failed to read converted file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read converted file",
		})
		return
	}

	// Set the appropriate headers
	outputFilename := strings.TrimSuffix(filepath.Base(fileHeader.Filename), filepath.Ext(fileHeader.Filename)) + "." + targetFormat
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFilename))
	c.Header("Content-Type", c.GetHeader("Content-Type"))

	// Send the file content directly
	c.Data(http.StatusOK, http.DetectContentType(result), result)
}

// handleGetFormats returns the supported file formats
func (s *Server) handleGetFormats(c *gin.Context) {
	// TODO: Return supported formats from converter factory
	c.JSON(http.StatusOK, gin.H{
		"data": map[string]interface{}{
			"document": map[string][]string{
				"input":  {"pdf", "docx", "doc", "odt"},
				"output": {"pdf", "docx", "odt", "txt"},
			},
			"image": map[string][]string{
				"input":  {"jpg", "jpeg", "png", "webp", "gif", "bmp"},
				"output": {"jpg", "png", "webp", "gif"},
			},
			"audio": map[string][]string{
				"input":  {"mp3", "wav", "ogg", "flac", "aac"},
				"output": {"mp3", "wav", "ogg", "flac"},
			},
			"video": map[string][]string{
				"input":  {"mp4", "webm", "mov", "avi", "mkv"},
				"output": {"mp4", "webm", "mov"},
			},
		},
	})
}

// loggingMiddleware returns a Gin middleware that logs HTTP requests
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		event := s.logger.Info()
		if len(errMsg) > 0 {
			event = s.logger.Error().Str("error", errMsg)
		}

		event.Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", statusCode).
			Str("latency", latency.String()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request")
	}
}

// rateLimitMiddleware returns a Gin middleware that implements rate limiting
func (s *Server) rateLimitMiddleware() gin.HandlerFunc {
	// TODO: Implement rate limiting
	return func(c *gin.Context) {
		c.Next()
	}
}


