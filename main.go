package main

import (
	"log"
	"os"

	"github.com/amannvl/freefileconverterz/internal/app"
	"github.com/amannvl/freefileconverterz/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run the application
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	os.Exit(0)
}
