package handlers

import (
	"log/slog"

	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/amannvl/freefileconverterz/internal/middleware"
	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter/factory"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config, store storage.Storage, conv *factory.ConverterFactory, logger *slog.Logger) {
	// Initialize handler with dependencies
	handler := NewHandler(cfg, store, conv, logger)

	// Serve static files
	app.Static("/downloads", "./uploads")
	app.Static("/", "./static")

	// Public routes
	setupPublicRoutes(app)

	// API v1 group
	api := app.Group("/api/v1")
	setupAPIRoutes(api, handler)

	// Admin routes (require admin role)
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))
	setupAdminRoutes(admin, handler)

	// Catch-all route for SPA
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// 404 handler
	app.Use(notFoundHandler)
}

func setupPublicRoutes(app *fiber.App) {
	// Homepage
	app.Get("/", homeHandler)

	// Static pages
	app.Get("/about", aboutHandler)
	app.Get("/contact", contactHandler)
	app.Get("/privacy", privacyHandler)
	app.Get("/terms", termsHandler)
}

func setupAPIRoutes(api fiber.Router, h *Handler) {
	// Health check
	api.Get("/health", h.HealthCheck)

	// File conversion
	api.Post("/convert", h.ConvertFile)
	api.Get("/convert/:id/status", h.GetConversionStatus)
	api.Get("/convert/:id/download", h.DownloadFile)

	// User management (public)
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
	api.Post("/refresh", h.RefreshToken)

	// Protected routes (require authentication)
	authorized := api.Group("", middleware.Protected())
	{
		authorized.Get("/conversions", h.ListConversions)
		authorized.Get("/conversions/:id", h.GetConversion)
		authorized.Delete("/conversions/:id", h.DeleteConversion)
	}
}

func setupAdminRoutes(admin fiber.Router, h *Handler) {
	// Admin dashboard
	admin.Get("/dashboard", h.AdminDashboard)

	// User management
	admin.Get("/users", h.ListUsers)
	admin.Get("/users/:id", h.GetUser)
	admin.Put("/users/:id", h.UpdateUser)
	admin.Delete("/users/:id", h.DeleteUser)

	// System stats
	admin.Get("/stats", h.GetStats)

	// System settings
	admin.Get("/settings", h.GetSettings)
	admin.Put("/settings", h.UpdateSettings)
}

// Handler functions
func homeHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "FreeFileConverterZ - Convert Any File Format Online",
	}, "layouts/main")
}

func aboutHandler(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "About FreeFileConverterZ",
	}, "layouts/main")
}

func contactHandler(c *fiber.Ctx) error {
	return c.Render("contact", fiber.Map{
		"Title": "Contact Us",
	}, "layouts/main")
}

func privacyHandler(c *fiber.Ctx) error {
	return c.Render("privacy", fiber.Map{
		"Title": "Privacy Policy",
	}, "layouts/main")
}

func termsHandler(c *fiber.Ctx) error {
	return c.Render("terms", fiber.Map{
		"Title": "Terms of Service",
	}, "layouts/main")
}

func notFoundHandler(c *fiber.Ctx) error {
	return c.Status(404).Render("errors/404", fiber.Map{
		"Title": "Page Not Found",
	}, "layouts/error")
}

// API Handlers
func convertHandler(c *fiber.Ctx) error {
	// TODO: Implement file conversion logic
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "File conversion endpoint",
	})
}

func statusHandler(c *fiber.Ctx) error {
	// TODO: Implement status check logic
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Status check endpoint",
	})
}

func downloadHandler(c *fiber.Ctx) error {
	// TODO: Implement download logic
	return c.SendString("Download endpoint")
}

// Admin Handlers
func adminDashboardHandler(c *fiber.Ctx) error {
	// TODO: Implement admin dashboard logic
	return c.Render("admin/dashboard", fiber.Map{
		"Title": "Admin Dashboard",
	}, "layouts/admin")
}

// User Management Handlers
func registerHandler(c *fiber.Ctx) error {
	// TODO: Implement user registration
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User registration endpoint",
	})
}

func loginHandler(c *fiber.Ctx) error {
	// TODO: Implement user login
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User login endpoint",
	})
}

func refreshTokenHandler(c *fiber.Ctx) error {
	// TODO: Implement token refresh
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Token refresh endpoint",
	})
}

func listUsersHandler(c *fiber.Ctx) error {
	// TODO: Implement list users
	return c.JSON(fiber.Map{
		"status": "success",
		"users":   []string{},
	})
}

func getUserHandler(c *fiber.Ctx) error {
	// TODO: Implement get user
	return c.JSON(fiber.Map{
		"status": "success",
		"user":   nil,
	})
}

func updateUserHandler(c *fiber.Ctx) error {
	// TODO: Implement update user
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func deleteUserHandler(c *fiber.Ctx) error {
	// TODO: Implement delete user
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

// Settings Handlers
func getSettingsHandler(c *fiber.Ctx) error {
	// TODO: Implement get settings
	return c.JSON(fiber.Map{
		"status":   "success",
		"settings": map[string]interface{}{},
	})
}

func updateSettingsHandler(c *fiber.Ctx) error {
	// TODO: Implement update settings
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Settings updated successfully",
	})
}
