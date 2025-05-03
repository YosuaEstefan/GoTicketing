package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ServerAddress string
	JWTSecret     string
	JWTExpiration time.Duration
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	jwtExpStr := getEnv("JWT_EXPIRATION_HOURS", "24")
	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration: %v", err)
	}

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "Yosua"),
		DBName:        getEnv("DB_NAME", "ticketing"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration: time.Duration(jwtExp) * time.Hour,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
