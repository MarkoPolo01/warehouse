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


	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: DB=%s:%d, User=%s, Port=%s", 
		cfg.DBHost, cfg.DBPortInt, cfg.DBUser, cfg.ServerPort)


	dbConfig := &database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPortInt,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	}

	log.Println("Attempting to connect to database...")
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database")


	log.Println("Initializing repository...")
	repo := repository.NewPostgresRepository(db)

	log.Println("Initializing service...")
	placementService := service.NewPlacementService(repo, cfg)
	
	log.Println("Initializing handler...")
	placementHandler := handler.NewPlacementHandler(placementService)


	log.Println("Setting up router...")
	router := gin.Default()


	log.Println("Registering routes...")
	placementHandler.RegisterRoutes(router)


	serverAddr := ":" + cfg.ServerPort
	log.Printf("Genetic Placement Service starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 