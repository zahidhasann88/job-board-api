package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	Port            string
	JWTSecret       string
	FileStoragePath string
	RecruiterRole   string
	JobSeekerRole   string
}

func LoadConfig() (*Config, error) {
	godotenv.Load() // Load .env file if it exists

	config := &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/jobboard?sslmode=disable"),
		Port:            getEnv("PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		FileStoragePath: getEnv("FILE_STORAGE_PATH", "./uploads"),
		RecruiterRole:   getEnv("RECRUITER_ROLE", "recruiter"),
		JobSeekerRole:   getEnv("JOB_SEEKER_ROLE", "job_seeker"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
