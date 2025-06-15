package main

import (
	"log"

	"warehouse/services/abc-placement/internal/config"
	"warehouse/services/abc-placement/internal/handler"
	"warehouse/services/abc-placement/internal/repository"
	"warehouse/services/abc-placement/internal/service"
	"warehouse/services/abc-placement/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()


	dbConfig := &database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPortInt,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable", // Adjust as needed
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()


	repo := repository.NewPostgresRepository(db)
	placementService := service.NewPlacementService(repo)
	placementHandler := handler.NewPlacementHandler(placementService)


	router := gin.Default()


	placementHandler.RegisterRoutes(router)


	serverAddr := ":" + cfg.ServerPort
	log.Printf("ABC Placement Service starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 