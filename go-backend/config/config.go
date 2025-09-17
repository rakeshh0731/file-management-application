package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds the application configuration, accessible globally.
var AppConfig Config

// Config defines all the configuration parameters for the application.
type Config struct {
	ServerPort     string
	AllowedOrigins string
	MongoURI       string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiresIn   time.Duration
	UploadDir      string
	MaxUploadSize  int64
}

// LoadConfig loads configuration from a .env file and the environment.
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = Config{
		ServerPort:     Getenv("SERVER_PORT", "8000"),
		AllowedOrigins: Getenv("ALLOWED_ORIGINS", "http://localhost:3000"),
		MongoURI:       Getenv("MONGO_URI", "mongodb://root:mongo@localhost:27017/filehub?authSource=admin"),
		DatabaseURL:    Getenv("DATABASE_URL", "postgres://admin:psql@localhost:5432/filehub_users?sslmode=disable"),
		JWTSecret:      Getenv("JWT_SECRET", "a-very-secret-key-that-should-be-long-and-random"),
		JWTExpiresIn:   getEnvAsDuration("JWT_EXPIRES_IN_HOURS", 24),
		UploadDir:      Getenv("UPLOAD_DIR", "uploads"),
		MaxUploadSize:  getEnvAsInt64("MAX_UPLOAD_SIZE_MB", 10) * 1024 * 1024, // Convert MB to bytes
	}
}

// Getenv retrieves an environment variable or returns a fallback.
func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt64 retrieves an environment variable as an int64 or returns a fallback.
func getEnvAsInt64(key string, fallback int64) int64 {
	if valueStr, ok := os.LookupEnv(key); ok {
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return value
		}
	}
	return fallback
}

// getEnvAsDuration retrieves an environment variable as a time.Duration (in hours) or returns a fallback.
func getEnvAsDuration(key string, fallbackHours int) time.Duration {
	hours := fallbackHours
	if valueStr, ok := os.LookupEnv(key); ok {
		if value, err := strconv.Atoi(valueStr); err == nil {
			hours = value
		}
	}
	return time.Duration(hours) * time.Hour
}
