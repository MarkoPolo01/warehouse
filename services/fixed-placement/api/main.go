package main

import (
	"log"

	"warehouse/services/fixed-placement/internal/config"
	"warehouse/services/fixed-placement/internal/handler"
	"warehouse/services/fixed-placement/internal/repository"
	"warehouse/services/fixed-placement/internal/service"
	"warehouse/services/fixed-placement/pkg/database"

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
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()


	repo := repository.NewPostgresRepository(db)
	placementService := service.NewPlacementService(repo)
	placementHandler := handler.NewPlacementHandler(placementService)


	router := gin.Default()
	placementHandler.RegisterRoutes(router)


	serverAddr := ":" + cfg.ServerPort
	log.Printf("Fixed Placement Service запущен на %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
} 