package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Storage  StorageConfig
	Security SecurityConfig
	Limits   RateLimits
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	Environment     string
}

type StorageConfig struct {
	UploadPath   string
	MaxUploadSize int64
}

type SecurityConfig struct {
	APIKey    string
	JWTSecret string
}

type RateLimits struct {
	Requests int
	Period   time.Duration
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 15*time.Second),
			Environment:     getEnv("ENV", "development"),
		},
		Storage: StorageConfig{
			UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
			MaxUploadSize: getInt64Env("MAX_UPLOAD_SIZE", 100<<20), // 100MB
		},
		Security: SecurityConfig{
			APIKey:    os.Getenv("API_KEY"),
			JWTSecret: getEnv("JWT_SECRET", "default-secret-key"),
		},
		Limits: RateLimits{
			Requests: getIntEnv("RATE_LIMIT_REQUESTS", 100),
			Period:   getDurationEnv("RATE_LIMIT_PERIOD", 1*time.Minute),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
