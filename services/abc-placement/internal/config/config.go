package config

import (
	"log"
	os"
	"strconv"
)

type Config struct {
	ServerPort  string
	DBHost      string
	DBPort      string
	DBPortInt   int
	DBUser      string
	DBPassword  string
	DBName      string
}

func LoadConfig() *Config {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8082" // Default port for ABC service
	}

	dbPortInt, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	return &Config{
		ServerPort: serverPort,
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBPortInt:  dbPortInt,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
} 