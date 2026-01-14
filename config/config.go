package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment       string
	ServerPort        string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	JWTExpiration     string
	RefreshExpiration string
	UploadPath        string
	MaxUploadSize     int64
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Environment:       getEnv("ENVIRONMENT", "development"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "auth_api"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration:     getEnv("JWT_EXPIRATION", "15m"),
		RefreshExpiration: getEnv("REFRESH_EXPIRATION", "7d"),
		UploadPath:        getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize:     5 * 1024 * 1024,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
