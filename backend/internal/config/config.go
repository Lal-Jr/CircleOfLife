package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found; assuming environments are set")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/community?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
