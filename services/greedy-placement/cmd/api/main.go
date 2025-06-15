package main

import (
	"log"

	"warehouse/services/greedy-placement/internal/config"
	"warehouse/services/greedy-placement/internal/handler"
	"warehouse/services/greedy-placement/internal/repository"
	"warehouse/services/greedy-placement/internal/service"
	"warehouse/services/greedy-placement/pkg/database"

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
		SSLMode:  "disable",
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
	log.Printf("Greedy Placement Service starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 