package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amannvl/freefileconverterz/internal/api"
	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/rs/zerolog/log"
)

// App represents the main application
type App struct {
	config *config.Config
	server *api.Server
}

// New creates a new App instance
func New(cfg *config.Config) (*App, error) {
	// Initialize logger
	logger := utils.NewLogger(utils.InfoLevel)
	if cfg.App.Env == "development" {
		logger = utils.NewLogger(utils.DebugLevel)
	}

	// Initialize file storage
	fileStorage, err := storage.NewFileStorage(cfg.Storage.UploadDir, cfg.Storage.MaxUploadSize, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file storage: %w", err)
	}

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(cfg.Storage.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Initialize converter factory
	converterFactory, err := converter.NewConverterFactory(cfg.Storage.TempDir, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize converter factory: %w", err)
	}

	// Initialize HTTP server
	server, err := api.NewServer(&api.Config{
		Port:         cfg.Server.Port,
		Environment:  cfg.App.Env,
		FileStorage:  fileStorage,
		Converter:    converterFactory,
		JWTSecret:    cfg.JWT.Secret,
		RateLimit:    cfg.Security.RateLimit,
		RateBurst:    cfg.Security.RateLimitBurst,
		CORSOrigins:  cfg.Security.CORSAllowedOrigins,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return &App{
		config: cfg,
		server: server,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Info().Str("port", a.config.Server.Port).Msg("Starting server")
		if err := a.server.Start(); err != nil {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		log.Info().Str("signal", sig.String()).Msg("Shutting down server...")
		return a.Shutdown()
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() error {
	// Shutdown HTTP server
	if err := a.server.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
		return err
	}

	// Perform any additional cleanup here
	log.Info().Msg("Application shutdown complete")
	return nil
}
