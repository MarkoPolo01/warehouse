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
	// Add weights for genetic algorithm fitness function
	WeightDistance       float64
	WeightSize           float64
	WeightStorageConditions float64
}

func LoadConfig() *Config {
	port := getEnv("DB_PORT", "5432")
	portInt, _ := strconv.Atoi(port)
	
	// Default weights for fitness function
	weightDist, _ := strconv.ParseFloat(getEnv("WEIGHT_DISTANCE", "1.0"), 64)
	weightSize, _ := strconv.ParseFloat(getEnv("WEIGHT_SIZE", "1.0"), 64)
	weightStorage, _ := strconv.ParseFloat(getEnv("WEIGHT_STORAGE", "1.0"), 64)

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8085"), // Порт для Genetic service
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
		DBPortInt:  portInt,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "postgres"),
		WeightDistance: weightDist,
		WeightSize: weightSize,
		WeightStorageConditions: weightStorage,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 