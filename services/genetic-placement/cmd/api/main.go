package main

import (
	"log"

	"warehouse/services/genetic-placement/internal/config"
	"warehouse/services/genetic-placement/internal/handler"
	"warehouse/services/genetic-placement/internal/repository"
	"warehouse/services/genetic-placement/internal/service"
	"warehouse/services/genetic-placement/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Genetic Placement Service...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: DB=%s:%d, User=%s, Port=%s", 
		cfg.DBHost, cfg.DBPortInt, cfg.DBUser, cfg.ServerPort)

	// Initialize database connection
	dbConfig := &database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPortInt,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable", // Adjust as needed
	}

	log.Println("Attempting to connect to database...")
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	// Initialize dependencies
	log.Println("Initializing repository...")
	repo := repository.NewPostgresRepository(db)
	// Pass config to service for fitness calculation weights
	log.Println("Initializing service...")
	placementService := service.NewPlacementService(repo, cfg)
	
	log.Println("Initializing handler...")
	placementHandler := handler.NewPlacementHandler(placementService)

	// Setup router
	log.Println("Setting up router...")
	router := gin.Default()

	// Register routes
	log.Println("Registering routes...")
	placementHandler.RegisterRoutes(router)

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Genetic Placement Service starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 