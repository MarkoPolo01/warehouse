package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBPortInt  int
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {
	port := getEnv("DB_PORT", "5432")
	portInt, _ := strconv.Atoi(port)
	
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8083"), // Порт для XYZ service
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
		DBPortInt:  portInt,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "postgres"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 