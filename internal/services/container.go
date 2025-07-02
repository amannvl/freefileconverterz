package services

import (
	"log"
	"os"

	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/amannvl/freefileconverterz/internal/storage"
)

// ServiceContainer holds all the service dependencies
type ServiceContainer struct {
	Config          *config.Config
	Storage         storage.Storage
	FileService     *FileService
	ConversionService *ConversionService
	LibreOfficePath string
}

// NewServiceContainer creates and initializes all services
func NewServiceContainer(cfg *config.Config) (*ServiceContainer, error) {
	// Initialize storage
	storageBackend, err := storage.NewStorage(cfg.Storage)
	if err != nil {
		return nil, err
	}

	// Create file service
	fileService := NewFileService(storageBackend, &cfg.Storage)

	// Set the full path to the LibreOffice binary
	libreOfficePath := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
	
	// Verify the binary exists and is executable
	if _, err := os.Stat(libreOfficePath); os.IsNotExist(err) {
		log.Printf("LibreOffice not found at %s, falling back to system PATH", libreOfficePath)
		libreOfficePath = "libreoffice" // Fall back to system PATH
	} else {
		// Make sure the binary is executable
		if err := os.Chmod(libreOfficePath, 0755); err != nil {
			log.Printf("Warning: Failed to set executable permissions on %s: %v", libreOfficePath, err)
		}
		log.Printf("Using LibreOffice at: %s", libreOfficePath)
	}

	// Create conversion service
	conversionService, err := NewConversionService(fileService, &cfg.Storage, libreOfficePath)
	if err != nil {
		return nil, err
	}

	return &ServiceContainer{
		Config:          cfg,
		Storage:         storageBackend,
		FileService:     fileService,
		ConversionService: conversionService,
		LibreOfficePath: libreOfficePath,
	}, nil
}

// Close cleans up any resources used by services
func (c *ServiceContainer) Close() error {
	// Close storage if it implements io.Closer
	if closer, ok := c.Storage.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

// EnsureTempDir ensures the temporary directory exists
func (c *ServiceContainer) EnsureTempDir() error {
	// Use the upload directory from config or default to ./uploads
	uploadDir := c.Config.App.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	return os.MkdirAll(uploadDir, 0755)
}
