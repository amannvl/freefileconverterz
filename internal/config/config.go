package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name      string
	Env       string
	Version   string
	BaseURL   string
	UploadDir string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	ExpiresIn  time.Duration
	RefreshIn  time.Duration
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	Provider           string        `mapstructure:"PROVIDER"`
	MaxUploadSize      int64         `mapstructure:"MAX_UPLOAD_SIZE"`
	AllowedFileTypes   []string      `mapstructure:"ALLOWED_FILE_TYPES"`
	TempDir           string        `mapstructure:"TEMP_DIR"`
	UploadDir         string        `mapstructure:"UPLOAD_DIR"`
	S3Region          string        `mapstructure:"S3_REGION"`
	S3Bucket          string        `mapstructure:"S3_BUCKET"`
	S3AccessKeyID     string        `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string        `mapstructure:"S3_SECRET_ACCESS_KEY"`
	S3Endpoint        string        `mapstructure:"S3_ENDPOINT"`
	S3UseSSL          bool          `mapstructure:"S3_USE_SSL"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimit          int
	RateLimitBurst     int
	CORSAllowedOrigins []string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
	File   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Print current working directory
	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)

	// Print environment variables
	fmt.Println("Environment variables:")
	fmt.Printf("UPLOAD_DIR: %s\n", os.Getenv("UPLOAD_DIR"))
	fmt.Printf("TEMP_DIR: %s\n", os.Getenv("TEMP_DIR"))

	// Load .env file if it exists
	envPath := filepath.Join(cwd, ".env")
	fmt.Printf("Looking for .env file at: %s\n", envPath)
	
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	} else {
		fmt.Println("Successfully loaded .env file")
		// Print environment variables again after loading .env
		fmt.Println("Environment variables after loading .env:")
		fmt.Printf("UPLOAD_DIR: %s\n", os.Getenv("UPLOAD_DIR"))
		fmt.Printf("TEMP_DIR: %s\n", os.Getenv("TEMP_DIR"))
	}

	// Set default values
	cfg := &Config{
		App: AppConfig{
			Name:      getEnv("APP_NAME", "FreeFileConverterZ"),
			Env:       getEnv("APP_ENV", "development"),
			Version:   getEnv("APP_VERSION", "1.0.0"),
			BaseURL:   getEnv("APP_URL", "http://localhost:3000"),
			UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		},
		Server: ServerConfig{
			Port:            getEnv("PORT", "3000"),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "freefileconverterz"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 300*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn: getEnvAsDuration("JWT_EXPIRES_IN", 24*time.Hour),
			RefreshIn: getEnvAsDuration("JWT_REFRESH_IN", 7*24*time.Hour),
		},
		Storage: StorageConfig{
			Provider:           getEnv("STORAGE_DRIVER", "local"),
			MaxUploadSize:      getEnvAsInt64("MAX_UPLOAD_SIZE", 100<<20), // 100MB
			AllowedFileTypes:   getEnvAsSlice("ALLOWED_FILE_TYPES", ",", []string{"image/*", "application/pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"}),
			TempDir:           getEnv("TEMP_DIR", "/tmp/freefileconverterz/temp"),
			UploadDir:         getEnv("UPLOAD_DIR", "/tmp/freefileconverterz/uploads"),
			S3Region:          getEnv("S3_REGION", ""),
			S3Bucket:          getEnv("S3_BUCKET", ""),
			S3AccessKeyID:     getEnv("S3_ACCESS_KEY", ""),
			S3SecretAccessKey: getEnv("S3_SECRET_KEY", ""),
			S3Endpoint:        getEnv("S3_ENDPOINT", ""),
			S3UseSSL:          getEnvAsBool("S3_USE_SSL", true),
		},
		Security: SecurityConfig{
			RateLimit:          getEnvAsInt("RATE_LIMIT", 100),
			RateLimitBurst:     getEnvAsInt("RATE_LIMIT_BURST", 50),
			CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", ",", []string{"*"}),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "text"),
			File:   getEnv("LOG_FILE", ""),
		},
	}

	// Create upload directory if it doesn't exist
	fmt.Printf("Attempting to create upload directory: %s\n", cfg.App.UploadDir)
	if err := os.MkdirAll(cfg.App.UploadDir, 0755); err != nil {
		fmt.Printf("Error creating upload directory: %v\n", err)
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	fmt.Printf("Successfully created/verified upload directory: %s\n", cfg.App.UploadDir)

	return cfg, nil
}

// Helper function to read an environment variable as a string
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper function to read an environment variable as an int64
func getEnvAsInt64(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultVal
}

// Helper function to read an environment variable as a boolean
func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper function to read an environment variable as a duration
func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(name, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper function to read an environment variable as a slice of strings
func getEnvAsSlice(name string, separator string, defaultVal []string) []string {
	valueStr := getEnv(name, "")
	if valueStr == "" {
		return defaultVal
	}

	values := []string{}
	for _, val := range strings.Split(valueStr, separator) {
		val = strings.TrimSpace(val)
		if val != "" {
			values = append(values, val)
		}
	}

	if len(values) == 0 {
		return defaultVal
	}

	return values
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// GetRedisAddr returns the Redis server address
func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

// GetUploadPath returns the full path to the upload directory
func (a *AppConfig) GetUploadPath() string {
	path, err := filepath.Abs(a.UploadDir)
	if err != nil {
		return a.UploadDir
	}
	return path
}
