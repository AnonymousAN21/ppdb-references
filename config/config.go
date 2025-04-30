package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database string
	JwtToken string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	return &Config{
		Port:     getEnv("PORT", ":8083"),
		Database: getEnv("DATABASE_URL", ""),
		JwtToken: getEnv("JWT_TOKEN", ""),
	}
}

func getEnv(key, defaultValues string) string {
	if values, exist := os.LookupEnv(key); exist {
		return values
	}
	return defaultValues
}
