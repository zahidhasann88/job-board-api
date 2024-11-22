package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/zahidhasann88/job-board-api/pkg/validator"
)

type Config struct {
	DatabaseURL                string
	Port                       string
	JWTSecret                  string
	FileStoragePath            string
	RecruiterRole              string
	JobSeekerRole              string
	AllowedPorts               []string
	RateLimitRequestsPerMinute int
	RateLimitBurstRequestCount int
}

func LoadConfig() (*Config, error) {
	godotenv.Load() // Load .env file if it exists

	config := &Config{
		DatabaseURL:                getEnv("DATABASE_URL", "postgres://postgres:admin123@localhost:5432/job_board?sslmode=disable"),
		Port:                       getEnv("PORT", "8080"),
		JWTSecret:                  getEnv("JWT_SECRET", "your-secret-key"),
		FileStoragePath:            getEnv("FILE_STORAGE_PATH", "./uploads"),
		RecruiterRole:              getEnv("RECRUITER_ROLE", "recruiter"),
		JobSeekerRole:              getEnv("JOB_SEEKER_ROLE", "job_seeker"),
		AllowedPorts:               strings.Split(getEnv("ALLOWED_PORTS", "8080"), ","),
		RateLimitRequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
		RateLimitBurstRequestCount: getEnvAsInt("RATE_LIMIT_BURST_COUNT", 50),
	}

	// Validate database URL
	if err := validator.ValidateDatabaseURL(config.DatabaseURL); err != nil {
		log.Fatalf("Invalid database URL: %v", err)
	}

	// Validate server port
	if err := validator.ValidatePort(config.Port); err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get env as int with default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
